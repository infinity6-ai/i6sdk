package targetingfinancedataz

import (
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/commonz/validation"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz"
)

type ClientFeatures struct {
	ID            string              `json:"id" bigquery:"id"`
	Dataset       bigquery.NullString `json:"dataset" bigquery:"dataset"`
	StoreId       string              `json:"store_id" bigquery:"store_id"`
	IngestedAt    string              `json:"ingested_at" bigquery:"ingested_at"`
	CreatedAt     time.Time           `json:"created_at" bigquery:"created_at"`
	FeaturesFloat []*FeatureFloat     `json:"features_float" bigquery:"features_float"`
}

func (c *ClientFeatures) Validate() error {
	_, err := storeid.Parse(c.StoreId)
	if err := tabledataz.FirstError(
		tabledataz.Required("id", c.ID),
		tabledataz.RequiredTime("created_at", c.CreatedAt),
		err,
	); err != nil {
		return err
	}
	if err := validation.NotEmpty(c.FeaturesFloat, "features_float must not be empty"); err != nil {
		return fmt.Errorf("failed to validate features_float: %w", err)
	}
	seen := make(map[string]struct{}, len(c.FeaturesFloat))
	for i, f := range c.FeaturesFloat {
		if f == nil {
			return fmt.Errorf("failed to validate features_float at index %d: feature cannot be nil", i)
		}
		if err := f.Validate(); err != nil {
			return fmt.Errorf("failed to validate features_float at index %d: %w", i, err)
		}
		if _, dup := seen[f.Key]; dup {
			return fmt.Errorf("failed to validate features_float at index %d: duplicated key %q", i, f.Key)
		}
		seen[f.Key] = struct{}{}
	}
	return nil
}

func (c *ClientFeatures) GetCreatedAt() *time.Time {
	return &c.CreatedAt
}

func (c *ClientFeatures) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("targeting-finance-client-features")
}

func (c *ClientFeatures) SetDataset(ds lakedsz.Dataset) {
	c.Dataset = tabledataz.NewBqNullString(ds.Get())
}

func (c *ClientFeatures) SetIngestedAt(p string) {
	c.IngestedAt = p
}

type FeatureFloat struct {
	Key   string  `json:"key" bigquery:"key"`
	Value float64 `json:"value" bigquery:"value"`
}

func (f *FeatureFloat) Validate() error {
	return tabledataz.Required("key", f.Key)
}

func (f *FeatureFloat) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Desc: &schemaz.Desc{
			Summary: "Generic numeric feature. Key is free-form (required, unique per row).",
		},
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"key": {
					Desc:  &schemaz.Desc{Summary: "Feature name, any value (required)."},
					Field: func() any { return &f.Key },
				},
				"value": {
					Desc:  &schemaz.Desc{Summary: "Feature numeric value."},
					Field: func() any { return &f.Value },
				},
			}
		},
	}
}

func (c *ClientFeatures) SchemaQuery() *schemaz.Schema {
	if c == nil {
		c = &ClientFeatures{}
	}
	return &schemaz.Schema{
		Desc: &schemaz.Desc{Summary: "ClientFeatures holds generic key/value columns for a client. The id is provided by the sender."},
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Desc:     &schemaz.Desc{Summary: "Unique row identifier, provided by the sender (required)."},
					UrlValue: schemaz.ParseStr(&c.ID),
				},
				"created_at": {
					Desc:     &schemaz.Desc{Summary: "Reference timestamp of the row (required)."},
					UrlValue: tabledataz.ParseTime(&c.CreatedAt),
				},
				"features_float": {
					Desc:     &schemaz.Desc{Summary: "List of generic numeric features (must not be empty)."},
					UrlValue: schemaz.ParseJson(&c.FeaturesFloat),
				},
			}
		},
	}
}

func (c *ClientFeatures) SchemaBody() *schemaz.Schema {
	if c == nil {
		c = &ClientFeatures{}
	}
	return &schemaz.Schema{
		Desc: &schemaz.Desc{Summary: "ClientFeatures holds generic key/value columns for a client. The id is provided by the sender."},
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Desc:  &schemaz.Desc{Summary: "Unique row identifier, provided by the sender (required)."},
					Field: func() any { return &c.ID },
				},
				"created_at": {
					Desc:  &schemaz.Desc{Summary: "Reference timestamp of the row (required)."},
					Field: func() any { return &c.CreatedAt },
				},
				"features_float": {
					Desc: &schemaz.Desc{Summary: "List of generic numeric features (must not be empty)."},
					Array: schemaz.NewArray(schemaz.ArrayOptions[[]*FeatureFloat, *FeatureFloat]{
						NewElement: func(idx int) *FeatureFloat { return &FeatureFloat{} },
						Slice:      func() *[]*FeatureFloat { return &c.FeaturesFloat },
						Get:        func(idx int, element **FeatureFloat) *schemaz.Schema { return (*element).Schema() },
					}),
				},
			}
		},
	}
}
