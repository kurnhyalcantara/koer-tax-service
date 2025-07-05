@echo off
setlocal enabledelayedexpansion

:: Check if any arguments are provided
if "%1"=="" (
    echo Usage: %0 ^<folder1^> [folder2] [folder3] ...
    exit /b 1
)

:: Get the full path of the current directory
for %%i in (".") do set "CURRENT_DIR=%%~fi"

set "OUTPUT_FOLDER=%CURRENT_DIR%\protogen"
set "BASE_PROTO=%CURRENT_DIR%\proto"
set "PLUGINS_PATH=%CURRENT_DIR%\proto\plugins"

:: Create output directory if it doesn't exist
if not exist "%OUTPUT_FOLDER%" mkdir "%OUTPUT_FOLDER%"

:: Process each provided folder argument
:process_args
if "%1"=="" goto :end

set "TARGET_FOLDER=%CURRENT_DIR%\proto\%1"

if not exist "%TARGET_FOLDER%" (
    echo Warning: Folder %TARGET_FOLDER% does not exist, skipping...
    shift
    goto :process_args
)

echo Processing module: %1
echo Target folder: %TARGET_FOLDER%

:: Process all .proto files
for /r "%TARGET_FOLDER%" %%F in (*.proto) do (
    echo Processing %%F
    protoc --proto_path="%BASE_PROTO%" ^
        --proto_path="%PLUGINS_PATH%" ^
        --plugin=protoc-gen-go="%GOPATH%\bin\protoc-gen-go.exe" ^
        --plugin=protoc-gen-govalidators="%GOPATH%\bin\protoc-gen-govalidators.exe" ^
        --go_out="%OUTPUT_FOLDER%" --go_opt=paths=source_relative ^
        --go-grpc_out="%OUTPUT_FOLDER%" --go-grpc_opt=paths=source_relative ^
        --govalidators_out="lang=go,paths=source_relative:%OUTPUT_FOLDER%" ^
        "%%F"
)

:: Process *_api.proto files
for /r "%TARGET_FOLDER%" %%F in (*_api.proto) do (
    echo Processing API: %%F
    protoc --proto_path="%BASE_PROTO%" ^
        --proto_path="%PLUGINS_PATH%" ^
        --plugin=protoc-gen-grpc-gateway="%GOPATH%\bin\protoc-gen-grpc-gateway.exe" ^
        --plugin=protoc-gen-openapiv2="%GOPATH%\bin\protoc-gen-openapiv2.exe" ^
        --plugin=protoc-gen-go-grpc="%GOPATH%\bin\protoc-gen-go-grpc.exe" ^
        --go-grpc_out="%OUTPUT_FOLDER%" --go-grpc_opt=paths=source_relative ^
        --grpc-gateway_out="%OUTPUT_FOLDER%" ^
        --grpc-gateway_opt=paths=source_relative,allow_delete_body=true,repeated_path_param_separator=ssv ^
        --openapiv2_out="%OUTPUT_FOLDER%" ^
        --openapiv2_opt=logtostderr=true,repeated_path_param_separator=ssv ^
        "%%F"
)

:: Process *_db.proto files
for /r "%TARGET_FOLDER%" %%F in (*_db.proto) do (
    echo Processing DB: %%F
    protoc --proto_path="%BASE_PROTO%" ^
        --proto_path="%PLUGINS_PATH%" ^
        --plugin=protoc-gen-gorm="%GOPATH%\bin\protoc-gen-gorm.exe" ^
        --go-grpc_out="%OUTPUT_FOLDER%" --go-grpc_opt=paths=source_relative ^
        --gorm_out=paths=source_relative:"%OUTPUT_FOLDER%" ^
        "%%F"
)

:: Move swagger file if it exists
if exist "%OUTPUT_FOLDER%\rnd-service\rnd_api.swagger.json" (
    if not exist "www" mkdir "www"
    move /Y "%OUTPUT_FOLDER%\rnd-service\rnd_api.swagger.json" "www\swagger.json"
)

shift
goto :process_args

:end
echo All modules processed successfully
endlocal