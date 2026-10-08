package sdkobjectz

import "github.com/infinity6-ai/gox/schemaz/schemaz"

type SdkObject interface {
	Schema() *schemaz.Schema
}
