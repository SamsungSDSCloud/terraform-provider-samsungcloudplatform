# Terraform Provider for Samsung Cloud Platform

## Directory structure

- `docs` : Auto-generated documentation
- `examples` : Example samsungcloudplatform terraform files for testing & documentation
- `internal` : (OLD) User profile management
- `samsungcloudplatform` : Samsungcloudplatform terraform provider source code
- `tools` : Tool modules

## Build Requirements

Reference development environment

- [Terraform](https://www.terraform.io) 1.1.7
- [Go](https://go.dev) v1.18

Prepare third-party libraries : `go mod tidy`

## Setup credentials

### Create local setting file

Create `.scp` directory in your OS account home

```
cd %USERPROFILE%
mkdir ".scp"
```

Create `config.json` and `credentials.json` in `.scp` directory

### Add Samsungcloudplatform configuration

Insert following parameters in `.scp/config.json` file

```
{
    "host": "https://openapi.samsungsdscloud.com",
    "user-id": "1234",
    "email" : "your.email@samsung.com",
    "project-id": "PROJECT-XXXXXXXXXXXXXXXX"
}
```

### Add your credentials

Insert following parameters in `.scp/credentials.json` file

```
{
    "auth-method": "access-key",
    "access-key": "XXXXXXXXXXXXXXXX,
    "secret-key": "XXXXXXXXXXXXXXXX"
}
```

## Developing Provider

### Build provider executable

Build a dummy plugin for placeholder
1. Build terraform debug executable `go build -o terraform-provider-samsungcloudplatform.exe`
2. Copy to predefined location. On Windows : `%APPDATA%\terraform.d\plugins\registry.terraform.io\SamsungSDSCloud\samsungcloudplatform\3.7.1\windows_amd64`
3. Go to `*.tf` example directory
4. Use `terraform init` to initialize plugin
    * When succeeded, following message will appear
      ```
      Initializing provider plugins...

      Terraform has been successfully initialized!

      You may now begin working with Terraform. Try running "terraform plan" to see
      any changes that are required for your infrastructure. All Terraform commands
      should now work.

      If you ever set or change modules or backend configuration for Terraform,
      rerun this command to reinitialize your working directory. If you forget, other
      commands will detect it and remind you to do so if necessary.
      ```

Run plugin with debug mode
* `go run main.go -- --debug`


### Development guideline

* Use `error` interface to handle errors
* Create test cases when possible
* Let the linter format your code
    * See `go lint`
* Use special go comments to auto generate documentation (See `godoc`)
    * Comment right before `package`, `func`, `struct`, ... will auto-detect as description
    * `// BUGS(author) ` will detect as bug comment
    * Function names are auto-detected
    * Comment with more than 2 spaces will be detected as code
      ```
      // Comment
      //  fmt.Println("Hello, World!")
      // Comment
      ```


## License

Copyright 2024. Samsung SDS Co., Ltd. All rights reserved.

See [LICENSE](LICENSE_MPL2.0) for details.

