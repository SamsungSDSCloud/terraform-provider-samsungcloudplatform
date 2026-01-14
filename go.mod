module github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3

go 1.24.0

toolchain go1.24.7

require github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3 v3.16.1



require (
	github.com/antihax/optional v1.0.0
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-docs v0.9.0
	github.com/hashicorp/terraform-plugin-log v0.4.0
	github.com/inconshreveable/mousetrap v1.1.0
	github.com/pmezard/go-difflib v1.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.44.0
	golang.org/x/mod v0.29.0
	golang.org/x/sync v0.18.0
	golang.org/x/tools v0.38.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463
	gopkg.in/yaml.v3 v3.0.1
)

// For avoid blackduck vulnerability detection
replace golang.org/x/crypto => golang.org/x/crypto v0.44.0

replace github.com/containerd/containerd => github.com/containerd/containerd v1.7.28

replace github.com/go-git/go-git/v5 => github.com/go-git/go-git/v5 v5.16.2

require github.com/go-git/go-git/v5 v5.14.0 // indirect

require (
	github.com/Masterminds/goutils v1.1.1
	github.com/Masterminds/semver/v3 v3.1.1 // indirectsp
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/agext/levenshtein v1.2.2
	github.com/apparentlymart/go-textseg/v13 v13.0.0
	github.com/armon/go-radix v1.0.0
	github.com/bgentry/speakeasy v0.1.0
	github.com/davecgh/go-spew v1.1.1
	github.com/fatih/color v1.13.0
	github.com/google/uuid v1.6.0
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-checkpoint v0.5.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog v1.2.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-plugin v1.4.3
	github.com/hashicorp/go-uuid v1.0.3
	github.com/hashicorp/go-version v1.5.0
	github.com/hashicorp/hc-install v0.3.2
	github.com/hashicorp/hcl/v2 v2.12.0
	github.com/hashicorp/logutils v1.0.0
	github.com/hashicorp/terraform-exec v0.16.1
	github.com/hashicorp/terraform-json v0.14.0
	github.com/hashicorp/terraform-plugin-go v0.9.0
	github.com/hashicorp/terraform-registry-address v0.0.0-20210412075316-9b2996cce896
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d
	github.com/huandu/xstrings v1.3.2
	github.com/imdario/mergo v0.3.13
	github.com/mattn/go-colorable v0.1.12
	github.com/mattn/go-isatty v0.0.14
	github.com/mitchellh/cli v1.1.4
	github.com/mitchellh/copystructure v1.2.0
	github.com/mitchellh/go-testing-interface v1.14.1
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mitchellh/reflectwalk v1.0.2
	github.com/oklog/run v1.0.0
	github.com/posener/complete v1.2.3
	github.com/russross/blackfriday v1.6.0
	github.com/shopspring/decimal v1.3.1
	github.com/spf13/cast v1.5.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	github.com/vmihailenco/msgpack/v4 v4.3.12
	github.com/vmihailenco/tagparser v0.1.1
	github.com/zclconf/go-cty v1.10.0
	golang.org/x/oauth2 v0.30.0
	golang.org/x/sys v0.38.0
	golang.org/x/text v0.31.0
	google.golang.org/grpc v1.74.0-dev
)

require (
	github.com/golang/protobuf v1.5.4
	github.com/google/go-cmp v0.7.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.16.0
	golang.org/x/net v0.47.0
	google.golang.org/appengine v1.6.7
	google.golang.org/protobuf v1.36.6
)
