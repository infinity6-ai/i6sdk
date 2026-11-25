module github.com/infinity6-ai/i6sdk/apitfz

go 1.26.1

require (
	github.com/infinity6-ai/gox/commonz v0.0.79
	github.com/infinity6-ai/gox/routez v0.0.79
	github.com/infinity6-ai/gox/schemaz v0.0.79
	github.com/infinity6-ai/i6sdk/sdkcorez v0.0.9
	github.com/infinity6-ai/i6sdk/versionz v0.0.9
	github.com/spf13/cobra v1.10.2
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.1.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/exp v0.0.0-20261121134650-e1094eeca1b5 // indirect
)

replace github.com/infinity6-ai/i6sdk/apicollectionz => ../apicollectionz

replace github.com/infinity6-ai/i6sdk/apidatasetz => ../apidatasetz

replace github.com/infinity6-ai/i6sdk/apiingestz => ../apiingestz

replace github.com/infinity6-ai/i6sdk/apipipez => ../apipipez

replace github.com/infinity6-ai/i6sdk/apiuserz => ../apiuserz

replace github.com/infinity6-ai/i6sdk/forgez => ../forgez

replace github.com/infinity6-ai/i6sdk/mold => ../mold

replace github.com/infinity6-ai/i6sdk/sdkcorez => ../sdkcorez

replace github.com/infinity6-ai/i6sdk/sdkgo => ../sdkgo

replace github.com/infinity6-ai/i6sdk/sdkpy => ../sdkpy

replace github.com/infinity6-ai/i6sdk/sdkweb => ../sdkweb

replace github.com/infinity6-ai/i6sdk/versionz => ../versionz
