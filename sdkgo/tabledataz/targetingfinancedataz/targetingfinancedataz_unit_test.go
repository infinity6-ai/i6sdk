package targetingfinancedataz_test

import (
	"encoding/json"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz/targetingfinancedataz"
)

var (
	_ tabledataz.TableData = (*targetingfinancedataz.ClientStatus)(nil)
	_ tabledataz.TableData = (*targetingfinancedataz.ClientEvent)(nil)
	_ tabledataz.TableData = (*targetingfinancedataz.Weight)(nil)
)

func validStatus() targetingfinancedataz.ClientStatus {
	return targetingfinancedataz.ClientStatus{
		ID:        "st_1",
		CreatedAt: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		StoreId:   "s1",
		UserId:    "u_1",
		UserCd:    "cd_1",
		UserSgGen: "f",
		UserTpRen: "pf",
		UserFlCon: 1,
		UserCdSuf: "SP",
		UserCdMun: "3550308",
		UserCdSeg: "seg_a",
	}
}

func validEvent() targetingfinancedataz.ClientEvent {
	return targetingfinancedataz.ClientEvent{
		ID:        "ev_1",
		CreatedAt: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		StoreId:   "s1",
		UserId:    "u_1",
		UserCd:    "cd_1",
		EventTp:   "purchase",
		EventCh:   "app",
		ItemId:    "item_1",
		ItemTpLin: "credit",
	}
}

func TestUnitClientStatusSerialization(t *testing.T) {
	t.Run("optional fields absent are null", func(t *testing.T) {
		var got targetingfinancedataz.ClientStatus
		require.NoError(t, json.Unmarshal([]byte(`{"id":"st_1","user_vl_epg":null,"user_tp_ren":"pj"}`), &got))
		require.Equal(t, "st_1", got.ID)
		require.False(t, got.UserVlEpg.Valid)
		require.False(t, got.UserNrCad.Valid)
		require.False(t, got.UserPcCut.Valid)
		require.Equal(t, "pj", got.UserTpRen)
	})

	t.Run("optional fields round trip", func(t *testing.T) {
		s := validStatus()
		s.UserVlEpg = tabledataz.NewBqNullFloat64(12.5)
		s.UserNrCad = tabledataz.NewBqNullInt64(3)
		s.UserPrCep = tabledataz.NewBqNullString("01310-100")
		s.UserPcCut = tabledataz.NewBqNullString("0.3")
		s.UserFlAtr = tabledataz.NewBqNullInt64(0)

		b, err := json.Marshal(s)
		require.NoError(t, err)
		require.Contains(t, string(b), `"user_vl_epg":12.5`)
		require.Contains(t, string(b), `"user_risk_score":null`)

		var got targetingfinancedataz.ClientStatus
		require.NoError(t, json.Unmarshal(b, &got))
		require.WithinDuration(t, s.CreatedAt, got.CreatedAt, time.Second)
		got.CreatedAt = s.CreatedAt
		require.Equal(t, s, got)
	})
}

func TestUnitClientStatusValidation(t *testing.T) {
	type scenario struct {
		mutate func(s *targetingfinancedataz.ClientStatus)
		errMsg string
	}
	check := func(t *testing.T, sc scenario) {
		t.Helper()
		s := validStatus()
		if sc.mutate != nil {
			sc.mutate(&s)
		}
		err := s.Validate()
		if sc.errMsg == "" {
			require.NoError(t, err)
			return
		}
		require.Error(t, err)
		require.Contains(t, err.Error(), sc.errMsg)
	}

	t.Run("valid", func(t *testing.T) { check(t, scenario{}) })
	t.Run("missing id", func(t *testing.T) {
		check(t, scenario{mutate: func(s *targetingfinancedataz.ClientStatus) { s.ID = "" }, errMsg: "failed to validate id"})
	})
	t.Run("missing created_at", func(t *testing.T) {
		check(t, scenario{mutate: func(s *targetingfinancedataz.ClientStatus) { s.CreatedAt = time.Time{} }, errMsg: "failed to validate created_at"})
	})
	t.Run("invalid user_tp_ren", func(t *testing.T) {
		check(t, scenario{mutate: func(s *targetingfinancedataz.ClientStatus) { s.UserTpRen = "xx" }, errMsg: "failed to validate user_tp_ren"})
	})
	t.Run("flag out of range", func(t *testing.T) {
		check(t, scenario{mutate: func(s *targetingfinancedataz.ClientStatus) { s.UserFlCar = 2 }, errMsg: "failed to validate user_fl_car"})
	})
	t.Run("optional flag out of range", func(t *testing.T) {
		check(t, scenario{mutate: func(s *targetingfinancedataz.ClientStatus) { s.UserFlAtr = tabledataz.NewBqNullInt64(5) }, errMsg: "failed to validate user_fl_atr"})
	})
	t.Run("optional flag null is valid", func(t *testing.T) {
		check(t, scenario{mutate: func(s *targetingfinancedataz.ClientStatus) { s.UserFlAtr = bigquery.NullInt64{} }})
	})
	t.Run("user_cd_suf too long", func(t *testing.T) {
		check(t, scenario{mutate: func(s *targetingfinancedataz.ClientStatus) { s.UserCdSuf = "ABCD" }, errMsg: "failed to validate user_cd_suf"})
	})
}

func TestUnitClientEventValidation(t *testing.T) {
	e := validEvent()
	require.NoError(t, e.Validate())

	e.UserPcCut = tabledataz.NewBqNullFloat64(0.25)
	require.NoError(t, e.Validate())

	e.ID = ""
	require.ErrorContains(t, e.Validate(), "failed to validate id")
}

func TestUnitWeightValidation(t *testing.T) {
	w := targetingfinancedataz.Weight{ID: "w_1", StoreId: "s1", Weight: 0.7}
	require.NoError(t, w.Validate())
	require.Nil(t, w.GetCreatedAt())

	w.ID = ""
	require.ErrorContains(t, w.Validate(), "failed to validate id")
}

func TestUnitSchemaQueryNullParsers(t *testing.T) {
	s := targetingfinancedataz.ClientStatus{}
	obj := s.SchemaQuery().Object(false)

	obj["user_vl_epg"].UrlValue().Parse("12.5")
	require.Equal(t, tabledataz.NewBqNullFloat64(12.5), s.UserVlEpg)
	v, _ := obj["user_vl_epg"].UrlValue().Format()
	require.Equal(t, "12.5", v)

	obj["user_nr_cad"].UrlValue().Parse("3")
	require.Equal(t, tabledataz.NewBqNullInt64(3), s.UserNrCad)

	obj["user_nr_cad"].UrlValue().Parse("")
	require.False(t, s.UserNrCad.Valid)
	v, _ = obj["user_nr_cad"].UrlValue().Format()
	require.Equal(t, "", v)
}

func TestUnitSchemas(t *testing.T) {
	ctx := t.Context()

	props := func(td *jstypez.Typedef) map[string]jstypez.Tuple {
		ret := map[string]jstypez.Tuple{}
		for _, p := range td.Properties.All() {
			ret[p.Name] = p
		}
		return ret
	}

	t.Run("ClientStatus body", func(t *testing.T) {
		td, _ := jsparser.ParseSchema(ctx, "ClientStatus", (&targetingfinancedataz.ClientStatus{}).SchemaBody())
		p := props(td)
		require.Equal(t, jstypez.TypString, p["id"].Typ)
		require.Equal(t, jstypez.TypTime, p["created_at"].Typ)
		require.Equal(t, jstypez.TypNumber, p["user_vl_sal"].Typ)
		require.Equal(t, jstypez.TypNumber, p["user_vl_epg"].Typ)
		require.Equal(t, jstypez.TypNumber, p["user_nr_cad"].Typ)
		require.Equal(t, jstypez.TypString, p["user_pr_cep"].Typ)
		require.Equal(t, jstypez.TypString, p["user_pc_cut"].Typ)
		require.NotContains(t, p, "partition")
		require.NotContains(t, p, "ingested_at")
	})

	t.Run("ClientEvent body", func(t *testing.T) {
		td, _ := jsparser.ParseSchema(ctx, "ClientEvent", (&targetingfinancedataz.ClientEvent{}).SchemaBody())
		p := props(td)
		require.Equal(t, jstypez.TypNumber, p["user_pc_cut"].Typ)
		require.Equal(t, jstypez.TypString, p["item_tp_sub"].Typ)
		require.NotContains(t, p, "updated_at")
	})

	t.Run("Query schemas", func(t *testing.T) {
		for name, s := range map[string]tabledataz.TableData{
			"ClientStatus": &targetingfinancedataz.ClientStatus{},
			"ClientEvent":  &targetingfinancedataz.ClientEvent{},
			"Weight":       &targetingfinancedataz.Weight{},
		} {
			td, _ := jsparser.ParseSchema(ctx, name, s.SchemaQuery())
			require.Equal(t, name, td.Head.Name)
			require.NotPanics(t, func() { s.Table() })
		}
	})
}
