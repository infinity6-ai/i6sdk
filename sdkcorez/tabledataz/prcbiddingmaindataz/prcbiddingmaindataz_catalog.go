package prcbiddingmaindataz

import (
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz"
)

type Catalog struct {
	ID            string              `json:"id" bigquery:"id"`
	Dataset       bigquery.NullString `json:"dataset" bigquery:"dataset"`
	Desc          string              `json:"desc" bigquery:"desc"`
	Supplier      string              `json:"supplier" bigquery:"supplier"`
	Status        string              `json:"status" bigquery:"status"`
	ExclusiveSale string              `json:"exclusive_sale" bigquery:"exclusive_sale"`
	IngestedAt    string              `json:"ingested_at" bigquery:"ingested_at"`
}

func (c *Catalog) Validate() error {
	return tabledataz.FirstError(
		tabledataz.Str("id", c.ID, 64),
		tabledataz.Required("desc", c.Desc),
		tabledataz.Required("supplier", c.Supplier),
		tabledataz.Required("status", c.Status),
		tabledataz.Required("exclusive_sale", c.ExclusiveSale),
	)
}

func (c *Catalog) GetCreatedAt() *time.Time {
	return nil
}

func (c *Catalog) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("prcbidding-main-catalog")
}

func (c *Catalog) SetDataset(ds lakedsz.Dataset) {
	c.Dataset = tabledataz.NewBqNullString(ds.Get())
}

func (c *Catalog) SetIngestedAt(p string) {
	c.IngestedAt = p
}

func (c *Catalog) SchemaQuery() *schemaz.Schema {
	if c == nil {
		c = &Catalog{}
	}
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Desc:     &schemaz.Desc{Summary: "Catalog entry id."},
					UrlValue: schemaz.ParseStr(&c.ID),
				},
				"desc": {
					Desc:     &schemaz.Desc{Summary: "Catalog entry description."},
					UrlValue: schemaz.ParseStr(&c.Desc),
				},
				"supplier": {
					Desc:     &schemaz.Desc{Summary: "Catalog entry supplier."},
					UrlValue: schemaz.ParseStr(&c.Supplier),
				},
				"status": {
					Desc:     &schemaz.Desc{Summary: "Catalog entry status."},
					UrlValue: schemaz.ParseStr(&c.Status),
				},
				"exclusive_sale": {
					Desc:     &schemaz.Desc{Summary: "Catalog entry exclusive sale."},
					UrlValue: schemaz.ParseStr(&c.ExclusiveSale),
				},
			}
		},
	}
}

func (c *Catalog) SchemaBody() *schemaz.Schema {
	if c == nil {
		c = &Catalog{}
	}
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Desc:  &schemaz.Desc{Summary: "Catalog entry id."},
					Field: func() any { return &c.ID },
				},
				"desc": {
					Desc:  &schemaz.Desc{Summary: "Catalog entry description."},
					Field: func() any { return &c.Desc },
				},
				"supplier": {
					Desc:  &schemaz.Desc{Summary: "Catalog entry supplier."},
					Field: func() any { return &c.Supplier },
				},
				"status": {
					Desc:  &schemaz.Desc{Summary: "Catalog entry status."},
					Field: func() any { return &c.Status },
				},
				"exclusive_sale": {
					Desc:  &schemaz.Desc{Summary: "Catalog entry exclusive sale."},
					Field: func() any { return &c.ExclusiveSale },
				},
			}
		},
	}
}
