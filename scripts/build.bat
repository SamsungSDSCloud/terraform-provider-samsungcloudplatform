@echo off
SET SCRIPT_PATH=%~dp0
cd %SCRIPT_PATH%\..
:: OUTPUT DIRECTORY
mkdir out\windows_amd64
mkdir out\linux_amd64
:: BUILD (WINDOWS)
SET GOOS=windows
SET GOARCH=amd64
go build -o out\windows_amd64\terraform-provider-samsungcloudplatform.exe
:: BUILD (LINUX)
SET GOOS=linux
SET GOARCH=amd64
go build -o out\linux_amd64\terraform-provider-samsungcloudplatform
:: CREATE PATH
::mkdir %APPDATA%\terraform.d\plugins\registry.terraform.io\samsungsds\scp\0.0.1\windows_amd64
:: MOVE
::move out\windows_amd64\terraform-provider-scp.exe %APPDATA%\terraform.d\plugins\registry.terraform.io\samsungsds\scp\0.0.1\windows_amd64
