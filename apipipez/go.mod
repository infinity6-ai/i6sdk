module github.com/infinity6-ai/i6sdk/apipipez

go 1.26.1

replace github.com/infinity6-ai/i6sdk/sdkcorez => ../sdkcorez

require (
	github.com/infinity6-ai/gox/commonz v0.0.77
	github.com/infinity6-ai/gox/routez v0.0.77
	github.com/infinity6-ai/gox/schemaz v0.0.77
	github.com/infinity6-ai/i6sdk/sdkcorez v0.0.0-00010101000000-000000000000
	github.com/infinity6-ai/i6sdk/versionz v0.0.1
	github.com/spf13/cobra v1.10.2
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/exp v0.0.0-20260824195058-e88cd73687aa // indirect
)

replace github.com/infinity6-ai/i6sdk/versionz => ../versionz
