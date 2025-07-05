package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	goruntime "runtime"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kurnhyalcantara/koer-tax-service/config"
	servicelogger "github.com/kurnhyalcantara/koer-tax-service/pkg/log"
	"github.com/kurnhyalcantara/koer-tax-service/pkg/validator"
	pb "github.com/kurnhyalcantara/koer-tax-service/protogen/koer-tax-service"
	grpcHandler "github.com/kurnhyalcantara/koer-tax-service/server/handler/grpc"
	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/db"
	inmemory "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/db/in_memory"
	grpcclient "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/grpc_client"
	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/interceptor"
	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/rescache"
	jwtmanager "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/security/jwt_manager"
	"github.com/kurnhyalcantara/koer-tax-service/server/repository"
	"github.com/kurnhyalcantara/koer-tax-service/server/usecase"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const defaultPort = 9090

var (
	appConfig *config.Config
	dbCore    db.DbCore
	s         *grpc.Server
)

func init() {
	appConfig = config.InitConfig()
}

func main() {
	app := cli.NewApp()
	app.Name = ""
	app.Commands = []cli.Command{
		grpcGatewayServerCmd(),
	}

	// Set max allocated cpu runtime for go routine parallelism
	setMaxAllocatedCPU()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		os.Exit(1)
	}
}

func grpcGatewayServerCmd() cli.Command {
	return cli.Command{
		Name:  "grpc-gw-server",
		Usage: "Starts gRPC and Gateway server",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port1",
				Value: defaultPort,
			},
			cli.IntFlag{
				Name:  "port2",
				Value: 3000,
			},
			cli.StringFlag{
				Name:  "grpc-endpoint",
				Value: ":" + fmt.Sprint(defaultPort),
				Usage: "the address of the running gRPC server to transcode to",
			},
		},
		Action: func(c *cli.Context) error {
			rpcPort, httpPort, grpcEndpoint := c.Int("port1"), c.Int("port2"), c.String("grpc-endpoint")

			go func() {
				if err := grpcServer(rpcPort); err != nil {
					log.Fatalf("failed RPC serve: %v", err)
				}
			}()

			go func() {
				if err := httpGatewayServer(httpPort, grpcEndpoint); err != nil {
					log.Fatalf("failed HTTP gateway serve: %v", err)
				}
			}()

			// Wait for Control C to exit
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)

			// Block until a signal is received
			<-ch

			CloseDBConnection(dbCore)

			log.Println("Stopping RPC and JSON Gateway server")
			s.Stop()
			log.Println("RPC server stopped")
			log.Println("JSON Gateway server stopped")

			return nil
		},
	}
}

func grpcServer(port int) error {

	logger := servicelogger.New(&servicelogger.LoggerConfig{
		Env:           appConfig.Env,
		ProductName:   fmt.Sprintf("%s", strings.ReplaceAll(strings.ToLower(appConfig.AppName), " ", "-")),
		ServiceName:   appConfig.LoggerTag,
		LogLevel:      appConfig.LoggerLevel,
		LogOutput:     appConfig.LoggerOutput,
		FluentbitHost: appConfig.FluentbitHost,
		FluentbitPort: appConfig.FluentbitPort,
	})
	// RPC
	logger.Info(servicelogger.LogPayload{
		Message: fmt.Sprintf("Starting RPC %s on port %d", appConfig.AppName, port),
	})

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed listen rpc serve: %v", err)
	}

	ctx := context.Background()

	dbCore = StartDBConnection()

	serviceClients := grpcclient.InitGrpcClients("", appConfig)

	jwtManager := jwtmanager.NewJwtManager(appConfig, serviceClients.AuthServiceClient)

	rdb := inmemory.New(appConfig.RedisAddress, appConfig.RedisPassword, appConfig.RedisDB)

	rdcache := rescache.NewResCache(
		rdb,
		jwtManager,
		fmt.Sprintf("%s", strings.ReplaceAll(strings.ToLower(appConfig.AppName), " ", "-")),
		appConfig.ResCacheTimeout,
		appConfig.ResCacheOnProgress,
		appConfig.ResCacheDoneTime,
		logger,
	)

	publisher := InitPublisherQueue(ctx)

	repo := repository.InitRepositories(dbCore, logger)

	uc := usecase.InitUseCases(usecase.Dependencies{
		Logger: logger,
		Repo:   repo,

		Manager:   jwtManager,
		Rcache:    rdcache,
		Publisher: publisher,

		AccountService: serviceClients.AccountServiceClient,
	})

	validator, err := validator.NewProtoValidator()
	if err != nil {
		log.Fatalf("failed init rpc validator: %v", err)
	}

	grpc_handler := grpcHandler.NewHandler(validator, uc, jwtManager)

	chainUnaryInterceptor := grpc.ChainUnaryInterceptor(
		interceptor.UnaryInterceptors(jwtManager),
	)

	streamInterceptor := grpc.StreamInterceptor(interceptor.StreamInterceptors(jwtManager))

	s = grpc.NewServer(chainUnaryInterceptor, streamInterceptor)

	pb.RegisterKoerTaxServiceServer(s, grpc_handler)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	return s.Serve(list)
}

func httpGatewayServer(port int, grpcEndpoint string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	conn, err := grpc.NewClient( // followed go 1.21.x version
		grpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(apmgrpc.NewStreamClientInterceptor()))
	if err != nil {
		return err
	}
	defer conn.Close()

	rmux := runtime.NewServeMux(
		runtime.WithErrorHandler(CustomHTTPError),
		runtime.WithForwardResponseOption(HttpResponseModifier),
	)

	client := pb.NewKoerTaxServiceClient(conn)
	err = pb.RegisterKoerTaxServiceHandlerClient(ctx, rmux, client)
	if err != nil {
		return err
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc("/api/tax/docs/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir(appConfig.SwaggerPath + "swagger-ui"))
	mux.Handle("/api/tax/docs/", http.StripPrefix("/api/tax/docs/", fs))

	// Start
	// logger.Info("", "httpGatewayServer", fmt.Sprintf("Starting JSON Gateway server on port %d...", port), nil, nil, nil, nil)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), apmhttp.Wrap(setHeaders(mux)))
}

func setMaxAllocatedCPU() {
	numCores := goruntime.NumCPU()

	// Use 90% of available cores
	desiredCores := int(float64(numCores) * 0.9)

	goruntime.GOMAXPROCS(desiredCores)
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, appConfig.SwaggerPath+"swagger.json")
}

func allowedOrigin(origin string) bool {
	if stringInSlice(viper.GetString("cors"), appConfig.CorsAllowedOrigins) {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

func setHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")

		if allowedOrigin(r.Header.Get("Origin")) {
			// w.Header().Set("Content-Security-Policy", "default-src 'self';img-src data: https:;object-src 'none'; upgrade-insecure-requests;")
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(appConfig.CorsAllowedOrigins, ", "))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(appConfig.CorsAllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(appConfig.CorsAllowedHeaders, ", "))
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
