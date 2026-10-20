package dsztoken

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type DatasetToken struct {
	Id             string          `json:"id"`
	Jti            string          `json:"jti"`
	Dataset        lakedsz.Dataset `json:"dataset"`
	Annotation     string          `json:"annotation"`
	ExpirationDate int64           `json:"expiration_date"`
	Audience       string          `json:"audience"`
	Grants         []string        `json:"grants"`
	Channel        string          `json:"channel,omitempty"`
	Token          string          `json:"token,omitempty"`
}

func (d *DatasetToken) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Raw: func() any { return &d.Id },
				},
				"jti": {
					Raw: func() any { return &d.Jti },
				},
				"dataset": {
					Raw: func() any { return &d.Dataset },
				},
				"annotation": {
					Raw: func() any { return &d.Annotation },
				},
				"expiration_date": {
					Raw: func() any { return &d.ExpirationDate },
				},
				"audience": {
					Raw: func() any { return &d.Audience },
				},
				"grants": {
					Array: schemaz.NewArray(schemaz.ArrayOptions[[]string, string]{
						Slice: func() *[]string { return &d.Grants },
						Get: func(idx int, element *string) *schemaz.Schema {
							return &schemaz.Schema{
								Raw: func() any { return element },
							}
						},
					}),
				},
				"channel": {
					Raw: func() any { return &d.Channel },
				},
				"token": {
					Raw: func() any { return &d.Token },
				},
			}
		},
	}
}
