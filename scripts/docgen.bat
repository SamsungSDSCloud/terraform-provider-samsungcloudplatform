@echo off
SET SCRIPT_PATH=%~dp0
cd %SCRIPT_PATH%\..

:: Generate documentation
:: go generate
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --examples-dir .\examples --rendered-provider-name scp --website-source-dir ./template --website-source-dir ./templates

cd %SCRIPT_PATH%\..
del docs\data-sources\firewall.md
::del docs\data-sources\firewalls.md
del docs\data-sources\kubernetes_apps_image.md
del docs\data-sources\public_ip.md
del docs\data-sources\region.md
del docs\data-sources\standard_image.md
