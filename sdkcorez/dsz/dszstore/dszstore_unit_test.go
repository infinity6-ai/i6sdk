package dszstore_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/domain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dszstore"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solutiondomain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
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
