package targetingfinancedataz

import (
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/commonz/validation"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz"
)

type ClientStatus struct {
	ID            string               `json:"id" bigquery:"id"`
	Dataset       bigquery.NullString  `json:"dataset" bigquery:"dataset"`
	IngestedAt    string               `json:"ingested_at" bigquery:"ingested_at"`
	CreatedAt     time.Time            `json:"created_at" bigquery:"created_at"`
	UserId        string               `json:"user_id" bigquery:"user_id"`
	UserCd        string               `json:"user_cd" bigquery:"user_cd"`
	UserVlSal     float64              `json:"user_vl_sal" bigquery:"user_vl_sal"`
	UserMnSal     float64              `json:"user_mn_sal" bigquery:"user_mn_sal"`
	UserMxSal     float64              `json:"user_mx_sal" bigquery:"user_mx_sal"`
	UserSmRmm     float64              `json:"user_sm_rmm" bigquery:"user_sm_rmm"`
	UserMdRmm     float64              `json:"user_md_rmm" bigquery:"user_md_rmm"`
	UserMnRmm     float64              `json:"user_mn_rmm" bigquery:"user_mn_rmm"`
	UserVlIrf     float64              `json:"user_vl_irf" bigquery:"user_vl_irf"`
	UserVlEmp     float64              `json:"user_vl_emp" bigquery:"user_vl_emp"`
	UserVlEpg     bigquery.NullFloat64 `json:"user_vl_epg" bigquery:"user_vl_epg"`
	UserNrDep     int64                `json:"user_nr_dep" bigquery:"user_nr_dep"`
	UserNrCad     bigquery.NullInt64   `json:"user_nr_cad" bigquery:"user_nr_cad"`
	UserVlCad     bigquery.NullFloat64 `json:"user_vl_cad" bigquery:"user_vl_cad"`
	UserSgGen     string               `json:"user_sg_gen" bigquery:"user_sg_gen"`
	UserTpRen     string               `json:"user_tp_ren" bigquery:"user_tp_ren"`
	UserFlCon     int64                `json:"user_fl_con" bigquery:"user_fl_con"`
	UserCdSuf     string               `json:"user_cd_suf" bigquery:"user_cd_suf"`
	UserCdMun     string               `json:"user_cd_mun" bigquery:"user_cd_mun"`
	UserPrCep     bigquery.NullString  `json:"user_pr_cep" bigquery:"user_pr_cep"`
	UserFlCco     int64                `json:"user_fl_cco" bigquery:"user_fl_cco"`
	UserCdSeg     string               `json:"user_cd_seg" bigquery:"user_cd_seg"`
	UserDsAcopen  int64                `json:"user_ds_acopen" bigquery:"user_ds_acopen"`
	UserVlLim     float64              `json:"user_vl_lim" bigquery:"user_vl_lim"`
	UserVlCon     float64              `json:"user_vl_con" bigquery:"user_vl_con"`
	UserFlCar     int64                `json:"user_fl_car" bigquery:"user_fl_car"`
	UserVlCar     float64              `json:"user_vl_car" bigquery:"user_vl_car"`
	UserVlCut     float64              `json:"user_vl_cut" bigquery:"user_vl_cut"`
	UserPcCut     bigquery.NullString  `json:"user_pc_cut" bigquery:"user_pc_cut"`
	UserMdTcr     float64              `json:"user_md_tcr" bigquery:"user_md_tcr"`
	UserPqTcr     float64              `json:"user_pq_tcr" bigquery:"user_pq_tcr"`
	UserSqTcr     float64              `json:"user_sq_tcr" bigquery:"user_sq_tcr"`
	UserTqTcr     float64              `json:"user_tq_tcr" bigquery:"user_tq_tcr"`
	UserMxTcr     float64              `json:"user_mx_tcr" bigquery:"user_mx_tcr"`
	UserRiskScore bigquery.NullFloat64 `json:"user_risk_score" bigquery:"user_risk_score"`
	UserFlAtr     bigquery.NullInt64   `json:"user_fl_atr" bigquery:"user_fl_atr"`
	UserVlAtr     float64              `json:"user_vl_atr" bigquery:"user_vl_atr"`
}

func (c *ClientStatus) Validate() error {
	if err := validation.StrNotEmpty(c.ID, "id is required"); err != nil {
		return fmt.Errorf("failed to validate id: %w", err)
	}
	if err := validation.False(c.CreatedAt.IsZero(), "created_at is required"); err != nil {
		return fmt.Errorf("failed to validate created_at: %w", err)
	}
	if err := validation.StrNotEmpty(c.UserId, "user_id is required"); err != nil {
		return fmt.Errorf("failed to validate user_id: %w", err)
	}
	if err := validation.LessOrEqual(len(c.UserId), 255, "user_id must have at most 255 characters"); err != nil {
		return fmt.Errorf("failed to validate user_id: %w", err)
	}
	if err := validation.StrNotEmpty(c.UserCd, "user_cd is required"); err != nil {
		return fmt.Errorf("failed to validate user_cd: %w", err)
	}
	if err := validation.LessOrEqual(len(c.UserCd), 255, "user_cd must have at most 255 characters"); err != nil {
		return fmt.Errorf("failed to validate user_cd: %w", err)
	}
	if err := validation.StrNotEmpty(c.UserSgGen, "user_sg_gen is required"); err != nil {
		return fmt.Errorf("failed to validate user_sg_gen: %w", err)
	}
	if err := validation.LessOrEqual(len(c.UserSgGen), 10, "user_sg_gen must have at most 10 characters"); err != nil {
		return fmt.Errorf("failed to validate user_sg_gen: %w", err)
	}
	if err := validation.StrNotEmpty(c.UserTpRen, "user_tp_ren is required"); err != nil {
		return fmt.Errorf("failed to validate user_tp_ren: %w", err)
	}
	if err := validation.LessOrEqual(len(c.UserTpRen), 2, "user_tp_ren must have at most 2 characters"); err != nil {
		return fmt.Errorf("failed to validate user_tp_ren: %w", err)
	}
	if err := validation.OneOf([]string{"pf", "pj"}, c.UserTpRen, "user_tp_ren must be one of: pf, pj"); err != nil {
		return fmt.Errorf("failed to validate user_tp_ren: %w", err)
	}
	if err := validation.GreaterOrEqual(c.UserFlCon, int64(0), "user_fl_con must be >= 0"); err != nil {
		return fmt.Errorf("failed to validate user_fl_con: %w", err)
	}
	if err := validation.LessOrEqual(c.UserFlCon, int64(1), "user_fl_con must be <= 1"); err != nil {
		return fmt.Errorf("failed to validate user_fl_con: %w", err)
	}
	if err := validation.StrNotEmpty(c.UserCdSuf, "user_cd_suf is required"); err != nil {
		return fmt.Errorf("failed to validate user_cd_suf: %w", err)
	}
	if err := validation.LessOrEqual(len(c.UserCdSuf), 3, "user_cd_suf must have at most 3 characters"); err != nil {
		return fmt.Errorf("failed to validate user_cd_suf: %w", err)
	}
	if err := validation.StrNotEmpty(c.UserCdMun, "user_cd_mun is required"); err != nil {
		return fmt.Errorf("failed to validate user_cd_mun: %w", err)
	}
	if err := validation.LessOrEqual(len(c.UserCdMun), 255, "user_cd_mun must have at most 255 characters"); err != nil {
		return fmt.Errorf("failed to validate user_cd_mun: %w", err)
	}
	if c.UserPrCep.Valid {
		if err := validation.LessOrEqual(len(c.UserPrCep.StringVal), 9, "user_pr_cep must have at most 9 characters"); err != nil {
			return fmt.Errorf("failed to validate user_pr_cep: %w", err)
		}
	}
	if err := validation.GreaterOrEqual(c.UserFlCco, int64(0), "user_fl_cco must be >= 0"); err != nil {
		return fmt.Errorf("failed to validate user_fl_cco: %w", err)
	}
	if err := validation.LessOrEqual(c.UserFlCco, int64(1), "user_fl_cco must be <= 1"); err != nil {
		return fmt.Errorf("failed to validate user_fl_cco: %w", err)
	}
	if err := validation.StrNotEmpty(c.UserCdSeg, "user_cd_seg is required"); err != nil {
		return fmt.Errorf("failed to validate user_cd_seg: %w", err)
	}
	if err := validation.LessOrEqual(len(c.UserCdSeg), 255, "user_cd_seg must have at most 255 characters"); err != nil {
		return fmt.Errorf("failed to validate user_cd_seg: %w", err)
	}
	if err := validation.GreaterOrEqual(c.UserFlCar, int64(0), "user_fl_car must be >= 0"); err != nil {
		return fmt.Errorf("failed to validate user_fl_car: %w", err)
	}
	if err := validation.LessOrEqual(c.UserFlCar, int64(1), "user_fl_car must be <= 1"); err != nil {
		return fmt.Errorf("failed to validate user_fl_car: %w", err)
	}
	if c.UserPcCut.Valid {
		if err := validation.LessOrEqual(len(c.UserPcCut.StringVal), 255, "user_pc_cut must have at most 255 characters"); err != nil {
			return fmt.Errorf("failed to validate user_pc_cut: %w", err)
		}
	}
	if c.UserFlAtr.Valid {
		if err := validation.GreaterOrEqual(c.UserFlAtr.Int64, int64(0), "user_fl_atr must be >= 0"); err != nil {
			return fmt.Errorf("failed to validate user_fl_atr: %w", err)
		}
		if err := validation.LessOrEqual(c.UserFlAtr.Int64, int64(1), "user_fl_atr must be <= 1"); err != nil {
			return fmt.Errorf("failed to validate user_fl_atr: %w", err)
		}
	}
	return nil
}

func (c *ClientStatus) GetCreatedAt() *time.Time {
	return &c.CreatedAt
}

func (c *ClientStatus) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("finance-predictivecampaign-client-status")
}

func (c *ClientStatus) SetDataset(ds lakedsz.Dataset) {
	c.Dataset = tabledataz.NewBqNullString(ds.Get())
}

func (c *ClientStatus) SetIngestedAt(p string) {
	c.IngestedAt = p
}

func (c *ClientStatus) SchemaQuery() *schemaz.Schema {
	if c == nil {
		c = &ClientStatus{}
	}
	return &schemaz.Schema{
		Desc: &schemaz.Desc{Summary: "ClientStatus is a row of the finance-predictivecampaign-client-status table. The id is provided by the sender."},
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
				"user_mn_rmm": {
					Desc:     &schemaz.Desc{Summary: "user_mn_rmm (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserMnRmm),
				},
				"user_vl_irf": {
					Desc:     &schemaz.Desc{Summary: "user_vl_irf (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlIrf),
				},
				"user_vl_emp": {
					Desc:     &schemaz.Desc{Summary: "user_vl_emp (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlEmp),
				},
				"user_vl_epg": {
					Desc:     &schemaz.Desc{Summary: "user_vl_epg (optional)."},
					UrlValue: tabledataz.ParseBqNullFloat64(&c.UserVlEpg),
				},
				"user_nr_dep": {
					Desc:     &schemaz.Desc{Summary: "user_nr_dep (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserNrDep),
				},
				"user_nr_cad": {
					Desc:     &schemaz.Desc{Summary: "user_nr_cad (optional)."},
					UrlValue: tabledataz.ParseBqNullInt64(&c.UserNrCad),
				},
				"user_vl_cad": {
					Desc:     &schemaz.Desc{Summary: "user_vl_cad (optional)."},
					UrlValue: tabledataz.ParseBqNullFloat64(&c.UserVlCad),
				},
				"user_sg_gen": {
					Desc:     &schemaz.Desc{Summary: "user_sg_gen (required)."},
					UrlValue: schemaz.ParseStr(&c.UserSgGen),
				},
				"user_tp_ren": {
					Desc:     &schemaz.Desc{Summary: "user_tp_ren (required)."},
					UrlValue: schemaz.ParseStr(&c.UserTpRen),
				},
				"user_fl_con": {
					Desc:     &schemaz.Desc{Summary: "user_fl_con (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserFlCon),
				},
				"user_cd_suf": {
					Desc:     &schemaz.Desc{Summary: "user_cd_suf (required)."},
					UrlValue: schemaz.ParseStr(&c.UserCdSuf),
				},
				"user_cd_mun": {
					Desc:     &schemaz.Desc{Summary: "user_cd_mun (required)."},
					UrlValue: schemaz.ParseStr(&c.UserCdMun),
				},
				"user_pr_cep": {
					Desc:     &schemaz.Desc{Summary: "user_pr_cep (optional)."},
					UrlValue: tabledataz.ParseBqNullString(&c.UserPrCep),
				},
				"user_fl_cco": {
					Desc:     &schemaz.Desc{Summary: "user_fl_cco (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserFlCco),
				},
				"user_cd_seg": {
					Desc:     &schemaz.Desc{Summary: "user_cd_seg (required)."},
					UrlValue: schemaz.ParseStr(&c.UserCdSeg),
				},
				"user_ds_acopen": {
					Desc:     &schemaz.Desc{Summary: "user_ds_acopen (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserDsAcopen),
				},
				"user_vl_lim": {
					Desc:     &schemaz.Desc{Summary: "user_vl_lim (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlLim),
				},
				"user_vl_con": {
					Desc:     &schemaz.Desc{Summary: "user_vl_con (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlCon),
				},
				"user_fl_car": {
					Desc:     &schemaz.Desc{Summary: "user_fl_car (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserFlCar),
				},
				"user_vl_car": {
					Desc:     &schemaz.Desc{Summary: "user_vl_car (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlCar),
				},
				"user_vl_cut": {
					Desc:     &schemaz.Desc{Summary: "user_vl_cut (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlCut),
				},
				"user_pc_cut": {
					Desc:     &schemaz.Desc{Summary: "user_pc_cut (optional)."},
					UrlValue: tabledataz.ParseBqNullString(&c.UserPcCut),
				},
				"user_md_tcr": {
					Desc:     &schemaz.Desc{Summary: "user_md_tcr (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserMdTcr),
				},
				"user_pq_tcr": {
					Desc:     &schemaz.Desc{Summary: "user_pq_tcr (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserPqTcr),
				},
				"user_sq_tcr": {
					Desc:     &schemaz.Desc{Summary: "user_sq_tcr (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserSqTcr),
				},
				"user_tq_tcr": {
					Desc:     &schemaz.Desc{Summary: "user_tq_tcr (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserTqTcr),
				},
				"user_mx_tcr": {
					Desc:     &schemaz.Desc{Summary: "user_mx_tcr (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserMxTcr),
				},
				"user_risk_score": {
					Desc:     &schemaz.Desc{Summary: "user_risk_score (optional)."},
					UrlValue: tabledataz.ParseBqNullFloat64(&c.UserRiskScore),
				},
				"user_fl_atr": {
					Desc:     &schemaz.Desc{Summary: "user_fl_atr (optional)."},
					UrlValue: tabledataz.ParseBqNullInt64(&c.UserFlAtr),
				},
				"user_vl_atr": {
					Desc:     &schemaz.Desc{Summary: "user_vl_atr (required)."},
					UrlValue: schemaz.ParseStrNumber(&c.UserVlAtr),
				},
			}
		},
	}
}

func (c *ClientStatus) SchemaBody() *schemaz.Schema {
	if c == nil {
		c = &ClientStatus{}
	}
	return &schemaz.Schema{
		Desc: &schemaz.Desc{Summary: "ClientStatus is a row of the finance-predictivecampaign-client-status table. The id is provided by the sender."},
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
				"user_mn_rmm": {
					Desc:  &schemaz.Desc{Summary: "user_mn_rmm (required)."},
					Field: func() any { return &c.UserMnRmm },
				},
				"user_vl_irf": {
					Desc:  &schemaz.Desc{Summary: "user_vl_irf (required)."},
					Field: func() any { return &c.UserVlIrf },
				},
				"user_vl_emp": {
					Desc:  &schemaz.Desc{Summary: "user_vl_emp (required)."},
					Field: func() any { return &c.UserVlEmp },
				},
				"user_vl_epg": {
					Desc:  &schemaz.Desc{Summary: "user_vl_epg (optional)."},
					Field: func() any { return &c.UserVlEpg },
				},
				"user_nr_dep": {
					Desc:  &schemaz.Desc{Summary: "user_nr_dep (required)."},
					Field: func() any { return &c.UserNrDep },
				},
				"user_nr_cad": {
					Desc:  &schemaz.Desc{Summary: "user_nr_cad (optional)."},
					Field: func() any { return &c.UserNrCad },
				},
				"user_vl_cad": {
					Desc:  &schemaz.Desc{Summary: "user_vl_cad (optional)."},
					Field: func() any { return &c.UserVlCad },
				},
				"user_sg_gen": {
					Desc:  &schemaz.Desc{Summary: "user_sg_gen (required)."},
					Field: func() any { return &c.UserSgGen },
				},
				"user_tp_ren": {
					Desc:  &schemaz.Desc{Summary: "user_tp_ren (required)."},
					Field: func() any { return &c.UserTpRen },
				},
				"user_fl_con": {
					Desc:  &schemaz.Desc{Summary: "user_fl_con (required)."},
					Field: func() any { return &c.UserFlCon },
				},
				"user_cd_suf": {
					Desc:  &schemaz.Desc{Summary: "user_cd_suf (required)."},
					Field: func() any { return &c.UserCdSuf },
				},
				"user_cd_mun": {
					Desc:  &schemaz.Desc{Summary: "user_cd_mun (required)."},
					Field: func() any { return &c.UserCdMun },
				},
				"user_pr_cep": {
					Desc:  &schemaz.Desc{Summary: "user_pr_cep (optional)."},
					Field: func() any { return &c.UserPrCep },
				},
				"user_fl_cco": {
					Desc:  &schemaz.Desc{Summary: "user_fl_cco (required)."},
					Field: func() any { return &c.UserFlCco },
				},
				"user_cd_seg": {
					Desc:  &schemaz.Desc{Summary: "user_cd_seg (required)."},
					Field: func() any { return &c.UserCdSeg },
				},
				"user_ds_acopen": {
					Desc:  &schemaz.Desc{Summary: "user_ds_acopen (required)."},
					Field: func() any { return &c.UserDsAcopen },
				},
				"user_vl_lim": {
					Desc:  &schemaz.Desc{Summary: "user_vl_lim (required)."},
					Field: func() any { return &c.UserVlLim },
				},
				"user_vl_con": {
					Desc:  &schemaz.Desc{Summary: "user_vl_con (required)."},
					Field: func() any { return &c.UserVlCon },
				},
				"user_fl_car": {
					Desc:  &schemaz.Desc{Summary: "user_fl_car (required)."},
					Field: func() any { return &c.UserFlCar },
				},
				"user_vl_car": {
					Desc:  &schemaz.Desc{Summary: "user_vl_car (required)."},
					Field: func() any { return &c.UserVlCar },
				},
				"user_vl_cut": {
					Desc:  &schemaz.Desc{Summary: "user_vl_cut (required)."},
					Field: func() any { return &c.UserVlCut },
				},
				"user_pc_cut": {
					Desc:  &schemaz.Desc{Summary: "user_pc_cut (optional)."},
					Field: func() any { return &c.UserPcCut },
				},
				"user_md_tcr": {
					Desc:  &schemaz.Desc{Summary: "user_md_tcr (required)."},
					Field: func() any { return &c.UserMdTcr },
				},
				"user_pq_tcr": {
					Desc:  &schemaz.Desc{Summary: "user_pq_tcr (required)."},
					Field: func() any { return &c.UserPqTcr },
				},
				"user_sq_tcr": {
					Desc:  &schemaz.Desc{Summary: "user_sq_tcr (required)."},
					Field: func() any { return &c.UserSqTcr },
				},
				"user_tq_tcr": {
					Desc:  &schemaz.Desc{Summary: "user_tq_tcr (required)."},
					Field: func() any { return &c.UserTqTcr },
				},
				"user_mx_tcr": {
					Desc:  &schemaz.Desc{Summary: "user_mx_tcr (required)."},
					Field: func() any { return &c.UserMxTcr },
				},
				"user_risk_score": {
					Desc:  &schemaz.Desc{Summary: "user_risk_score (optional)."},
					Field: func() any { return &c.UserRiskScore },
				},
				"user_fl_atr": {
					Desc:  &schemaz.Desc{Summary: "user_fl_atr (optional)."},
					Field: func() any { return &c.UserFlAtr },
				},
				"user_vl_atr": {
					Desc:  &schemaz.Desc{Summary: "user_vl_atr (required)."},
					Field: func() any { return &c.UserVlAtr },
				},
			}
		},
	}
}
