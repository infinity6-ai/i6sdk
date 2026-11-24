package targetingfinancedataz

import (
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz"
)

type ClientEvent struct {
	ID                       string               `json:"id" bigquery:"id"`
	Dataset                  bigquery.NullString  `json:"dataset" bigquery:"dataset"`
	StoreId                  string               `json:"store_id" bigquery:"store_id"`
	IngestedAt               string               `json:"ingested_at" bigquery:"ingested_at"`
	CreatedAt                time.Time            `json:"created_at" bigquery:"created_at"`
	UserId                   string               `json:"user_id" bigquery:"user_id"`
	UserCd                   string               `json:"user_cd" bigquery:"user_cd"`
	UserVlSal                float64              `json:"user_vl_sal" bigquery:"user_vl_sal"`
	UserMnSal                float64              `json:"user_mn_sal" bigquery:"user_mn_sal"`
	UserMxSal                float64              `json:"user_mx_sal" bigquery:"user_mx_sal"`
	UserSmRmm                float64              `json:"user_sm_rmm" bigquery:"user_sm_rmm"`
	UserMdRmm                float64              `json:"user_md_rmm" bigquery:"user_md_rmm"`
	UserRiskScore            bigquery.NullFloat64 `json:"user_risk_score" bigquery:"user_risk_score"`
	UserVlCon                float64              `json:"user_vl_con" bigquery:"user_vl_con"`
	UserVlLim                float64              `json:"user_vl_lim" bigquery:"user_vl_lim"`
	UserPcCut                bigquery.NullFloat64 `json:"user_pc_cut" bigquery:"user_pc_cut"`
	EventTp                  string               `json:"event_tp" bigquery:"event_tp"`
	EventCh                  string               `json:"event_ch" bigquery:"event_ch"`
	EventDt                  float64              `json:"event_dt" bigquery:"event_dt"`
	EventOr                  int64                `json:"event_or" bigquery:"event_or"`
	DaysSinceLastEvent       int64                `json:"days_since_last_event" bigquery:"days_since_last_event"`
	DaysSinceLastSameProduct int64                `json:"days_since_last_same_product" bigquery:"days_since_last_same_product"`
	ItemId                   string               `json:"item_id" bigquery:"item_id"`
	ItemPz                   float64              `json:"item_pz" bigquery:"item_pz"`
	ItemTx                   float64              `json:"item_tx" bigquery:"item_tx"`
	ItemTpLin                string               `json:"item_tp_lin" bigquery:"item_tp_lin"`
	ItemTpSub                bigquery.NullString  `json:"item_tp_sub" bigquery:"item_tp_sub"`
	ItemVlCon                float64              `json:"item_vl_con" bigquery:"item_vl_con"`
}

func (c *ClientEvent) Validate() error {
	_, err := storeid.Parse(c.StoreId)
	return tabledataz.FirstError(
		tabledataz.Required("id", c.ID),
		tabledataz.RequiredTime("created_at", c.CreatedAt),
		err,
		tabledataz.Str("user_id", c.UserId, 255),
		tabledataz.Str("user_cd", c.UserCd, 255),
		tabledataz.Str("event_tp", c.EventTp, 255),
		tabledataz.Str("event_ch", c.EventCh, 255),
		tabledataz.Str("item_id", c.ItemId, 255),
		tabledataz.Str("item_tp_lin", c.ItemTpLin, 255),
		tabledataz.NullMaxLen("item_tp_sub", c.ItemTpSub, 255),
	)
}

func (c *ClientEvent) GetCreatedAt() *time.Time {
	return &c.CreatedAt
}

func (c *ClientEvent) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("targeting-finance-client-events")
}

func (c *ClientEvent) SetDataset(ds lakedsz.Dataset) {
	c.Dataset = tabledataz.NewBqNullString(ds.Get())
}

func (c *ClientEvent) SetIngestedAt(p string) {
	c.IngestedAt = p
}

func (c *ClientEvent) SchemaQuery() *schemaz.Schema {
	if c == nil {
		c = &ClientEvent{}
	}
	return &schemaz.Schema{
		Desc: &schemaz.Desc{Summary: "ClientEvent is a row of the targeting-finance-client-events table. The id is provided by the sender."},
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
				"user_id": {
					Desc:     &schemaz.Desc{Summary: "user_id (required)."},
					UrlValue: schemaz.ParseStr(&c.UserId),
				},
				"user_cd": {
					Desc:     &schemaz.Desc{Summary: "user_cd (required)."},
					UrlValue: schemaz.ParseStr(&c.UserCd),
				},
				"user_vl_sal": {
					Desc:     &schemaz.Desc{Summary: "user_vl_sal (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlSal),
				},
				"user_mn_sal": {
					Desc:     &schemaz.Desc{Summary: "user_mn_sal (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserMnSal),
				},
				"user_mx_sal": {
					Desc:     &schemaz.Desc{Summary: "user_mx_sal (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserMxSal),
				},
				"user_sm_rmm": {
					Desc:     &schemaz.Desc{Summary: "user_sm_rmm (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserSmRmm),
				},
				"user_md_rmm": {
					Desc:     &schemaz.Desc{Summary: "user_md_rmm (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserMdRmm),
				},
				"user_risk_score": {
					Desc:     &schemaz.Desc{Summary: "user_risk_score (optional)."},
					UrlValue: tabledataz.ParseBqNullFloat64(&c.UserRiskScore),
				},
				"user_vl_con": {
					Desc:     &schemaz.Desc{Summary: "user_vl_con (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlCon),
				},
				"user_vl_lim": {
					Desc:     &schemaz.Desc{Summary: "user_vl_lim (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlLim),
				},
				"user_pc_cut": {
					Desc:     &schemaz.Desc{Summary: "user_pc_cut (optional)."},
					UrlValue: tabledataz.ParseBqNullFloat64(&c.UserPcCut),
				},
				"event_tp": {
					Desc:     &schemaz.Desc{Summary: "event_tp (required)."},
					UrlValue: schemaz.ParseStr(&c.EventTp),
				},
				"event_ch": {
					Desc:     &schemaz.Desc{Summary: "event_ch (required)."},
					UrlValue: schemaz.ParseStr(&c.EventCh),
				},
				"event_dt": {
					Desc:     &schemaz.Desc{Summary: "event_dt (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.EventDt),
				},
				"event_or": {
					Desc:     &schemaz.Desc{Summary: "event_or (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.EventOr),
				},
				"days_since_last_event": {
					Desc:     &schemaz.Desc{Summary: "days_since_last_event (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.DaysSinceLastEvent),
				},
				"days_since_last_same_product": {
					Desc:     &schemaz.Desc{Summary: "days_since_last_same_product (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.DaysSinceLastSameProduct),
				},
				"item_id": {
					Desc:     &schemaz.Desc{Summary: "item_id (required)."},
					UrlValue: schemaz.ParseStr(&c.ItemId),
				},
				"item_pz": {
					Desc:     &schemaz.Desc{Summary: "item_pz (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.ItemPz),
				},
				"item_tx": {
					Desc:     &schemaz.Desc{Summary: "item_tx (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.ItemTx),
				},
				"item_tp_lin": {
					Desc:     &schemaz.Desc{Summary: "item_tp_lin (required)."},
					UrlValue: schemaz.ParseStr(&c.ItemTpLin),
				},
				"item_tp_sub": {
					Desc:     &schemaz.Desc{Summary: "item_tp_sub (optional)."},
					UrlValue: tabledataz.ParseBqNullString(&c.ItemTpSub),
				},
				"item_vl_con": {
					Desc:     &schemaz.Desc{Summary: "item_vl_con (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.ItemVlCon),
				},
			}
		},
	}
}

func (c *ClientEvent) SchemaBody() *schemaz.Schema {
	if c == nil {
		c = &ClientEvent{}
	}
	return &schemaz.Schema{
		Desc: &schemaz.Desc{Summary: "ClientEvent is a row of the targeting-finance-client-events table. The id is provided by the sender."},
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
				"user_id": {
					Desc:  &schemaz.Desc{Summary: "user_id (required)."},
					Field: func() any { return &c.UserId },
				},
				"user_cd": {
					Desc:  &schemaz.Desc{Summary: "user_cd (required)."},
					Field: func() any { return &c.UserCd },
				},
				"user_vl_sal": {
					Desc:  &schemaz.Desc{Summary: "user_vl_sal (required)."},
					Field: func() any { return &c.UserVlSal },
				},
				"user_mn_sal": {
					Desc:  &schemaz.Desc{Summary: "user_mn_sal (required)."},
					Field: func() any { return &c.UserMnSal },
				},
				"user_mx_sal": {
					Desc:  &schemaz.Desc{Summary: "user_mx_sal (required)."},
					Field: func() any { return &c.UserMxSal },
				},
				"user_sm_rmm": {
					Desc:  &schemaz.Desc{Summary: "user_sm_rmm (required)."},
					Field: func() any { return &c.UserSmRmm },
				},
				"user_md_rmm": {
					Desc:  &schemaz.Desc{Summary: "user_md_rmm (required)."},
					Field: func() any { return &c.UserMdRmm },
				},
				"user_risk_score": {
					Desc:  &schemaz.Desc{Summary: "user_risk_score (optional)."},
					Field: func() any { return &c.UserRiskScore },
				},
				"user_vl_con": {
					Desc:  &schemaz.Desc{Summary: "user_vl_con (required)."},
					Field: func() any { return &c.UserVlCon },
				},
				"user_vl_lim": {
					Desc:  &schemaz.Desc{Summary: "user_vl_lim (required)."},
					Field: func() any { return &c.UserVlLim },
				},
				"user_pc_cut": {
					Desc:  &schemaz.Desc{Summary: "user_pc_cut (optional)."},
					Field: func() any { return &c.UserPcCut },
				},
				"event_tp": {
					Desc:  &schemaz.Desc{Summary: "event_tp (required)."},
					Field: func() any { return &c.EventTp },
				},
				"event_ch": {
					Desc:  &schemaz.Desc{Summary: "event_ch (required)."},
					Field: func() any { return &c.EventCh },
				},
				"event_dt": {
					Desc:  &schemaz.Desc{Summary: "event_dt (required)."},
					Field: func() any { return &c.EventDt },
				},
				"event_or": {
					Desc:  &schemaz.Desc{Summary: "event_or (required)."},
					Field: func() any { return &c.EventOr },
				},
				"days_since_last_event": {
					Desc:  &schemaz.Desc{Summary: "days_since_last_event (required)."},
					Field: func() any { return &c.DaysSinceLastEvent },
				},
				"days_since_last_same_product": {
					Desc:  &schemaz.Desc{Summary: "days_since_last_same_product (required)."},
					Field: func() any { return &c.DaysSinceLastSameProduct },
				},
				"item_id": {
					Desc:  &schemaz.Desc{Summary: "item_id (required)."},
					Field: func() any { return &c.ItemId },
				},
				"item_pz": {
					Desc:  &schemaz.Desc{Summary: "item_pz (required)."},
					Field: func() any { return &c.ItemPz },
				},
				"item_tx": {
					Desc:  &schemaz.Desc{Summary: "item_tx (required)."},
					Field: func() any { return &c.ItemTx },
				},
				"item_tp_lin": {
					Desc:  &schemaz.Desc{Summary: "item_tp_lin (required)."},
					Field: func() any { return &c.ItemTpLin },
				},
				"item_tp_sub": {
					Desc:  &schemaz.Desc{Summary: "item_tp_sub (optional)."},
					Field: func() any { return &c.ItemTpSub },
				},
				"item_vl_con": {
					Desc:  &schemaz.Desc{Summary: "item_vl_con (required)."},
					Field: func() any { return &c.ItemVlCon },
				},
			}
		},
	}
}
