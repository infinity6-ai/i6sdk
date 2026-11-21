package dszstore_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/domain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dszstore"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solutiondomain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/storeid"
)

func TestUnitDszStoreSchema(t *testing.T) {
	store := &dszstore.DszStore{
		Id:   storeid.MustParse("my-store"),
		Name: "My Awesome Store",
		SolutionDomains: []solutiondomain.SolutionDomain{
			solutiondomain.MustParse(solutiondomain.Value{
				Solution: solution.Relevance().Solution,
				Domain:   domain.MustParse("fashion"),
			}),
		},
	}

	schema := store.Schema()
	require.NotNil(t, schema)
	require.NotNil(t, schema.Object)

	props := schema.Object(true)
	require.Len(t, props, 3)

	require.Contains(t, props, "id")
	require.Contains(t, props, "name")
	require.Contains(t, props, "solution_domains")
	require.NotNil(t, props["solution_domains"].Array)
}
