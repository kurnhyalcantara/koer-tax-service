@echo off
setlocal enabledelayedexpansion

if "%1"=="" (
    echo Usage: %0 ^<folder^>
    exit /b 1
)

:: Get absolute path of current directory
for %%i in (".") do set "CURRENT_DIR=%%~fi"

set "TARGET_FOLDER=%CURRENT_DIR%\%1"
set "PLUGINS_PATH=%CURRENT_DIR%\plugins"

if not exist "%TARGET_FOLDER%" (
    echo Folder %TARGET_FOLDER% does not exist.
    exit /b 1
)

if exist verify rmdir /s /q verify
mkdir verify

echo Processing %TARGET_FOLDER%

:: Process all .proto files
for /r "%TARGET_FOLDER%" %%F in (*.proto) do (
    echo Processing %%F
    protoc --proto_path="%CURRENT_DIR%" ^
        --proto_path="%PLUGINS_PATH%" ^
        --plugin=protoc-gen-go="%GOPATH%\bin\protoc-gen-go.exe" ^
        --plugin=protoc-gen-govalidators="%GOPATH%\bin\protoc-gen-govalidators.exe" ^
        --go_out=verify --go_opt=paths=source_relative ^
        --go-grpc_out=verify --go-grpc_opt=paths=source_relative ^
        --govalidators_out="lang=go,paths=source_relative:verify" ^
        "%%F"
)

:: Process *_api.proto files
for /r "%TARGET_FOLDER%" %%F in (*_api.proto) do (
    echo Processing %%F
    protoc --proto_path="%CURRENT_DIR%" ^
        --proto_path="%PLUGINS_PATH%" ^
        --plugin=protoc-gen-grpc-gateway="%GOPATH%\bin\protoc-gen-grpc-gateway.exe" ^
        --plugin=protoc-gen-openapiv2="%GOPATH%\bin\protoc-gen-openapiv2.exe" ^
        --plugin=protoc-gen-go-grpc="%GOPATH%\bin\protoc-gen-go-grpc.exe" ^
        --go-grpc_out=verify --go-grpc_opt=paths=source_relative ^
        --grpc-gateway_out=verify ^
        --grpc-gateway_opt=paths=source_relative,allow_delete_body=true,repeated_path_param_separator=ssv ^
        --openapiv2_out=verify ^
        --openapiv2_opt=logtostderr=true,repeated_path_param_separator=ssv ^
        "%%F"
)

:: Process *_db.proto files
for /r "%TARGET_FOLDER%" %%F in (*_db.proto) do (
    echo Processing %%F
    protoc --proto_path="%CURRENT_DIR%" ^
        --proto_path="%PLUGINS_PATH%" ^
        --plugin=protoc-gen-gorm="%GOPATH%\bin\protoc-gen-gorm.exe" ^
        --go-grpc_out=verify --go-grpc_opt=paths=source_relative ^
        --gorm_out=paths=source_relative:verify ^
        "%%F"
)

rmdir /s /q verify

echo Verify successful