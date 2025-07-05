package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/kurnhyalcantara/koer-tax-service/protogen/koer-tax-service"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	custErr := status.Convert(err)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(runtime.HTTPStatusFromCode(custErr.Code()))

	body := &pb.GeneralBodyResponse{
		Error:   true,
		Code:    uint32(runtime.HTTPStatusFromCode(custErr.Code())),
		Message: custErr.Message(),
	}

	jErr := json.NewEncoder(w).Encode(body)

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}

func HttpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// download file, http response modifier
	if vals := md.HeaderMD.Get("file-download"); len(vals) > 0 {

		delete(md.HeaderMD, "file-download")
		delete(w.Header(), "Grpc-Metadata-File-Download")

		w.Header().Set("Content-Disposition", md.HeaderMD.Get("Content-Disposition")[0])
		w.Header().Set("Content-Length", md.HeaderMD.Get("Content-Length")[0])
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}

	return nil
}
