package tokens_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccesscreate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefresh"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthrefresh"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokenstore"
	"github.com/stretchr/testify/require"
)

func TestUnitUserTokenCRUD(t *testing.T) {
	ctx := t.Context()
	store := tokenstore.NewMemory()
	mgr := tokens.New(tokens.ManagerOptions{Store: store})

	t.Run("Get missing user token returns nil without error", func(t *testing.T) {
		tok, err := mgr.GetUserToken(ctx)
		require.NoError(t, err)
		require.Nil(t, tok)
	})

	t.Run("Set and Get stores and retrieves user token", func(t *testing.T) {
		rfExp := time.Now().Add(30 * 24 * time.Hour).UnixMilli()
		input := &tokens.UserToken{
			UserAccessToken:       "access-token-123",
			UserRefreshToken:      "refresh-token-456",
			UserId:                "usr_abc",
			ExpiresAt:             time.Now().Add(1 * time.Hour).UnixMilli(),
			RefreshTokenExpiresAt: &rfExp,
		}

		err := mgr.SetUserToken(ctx, input)
		require.NoError(t, err)

		got, err := mgr.GetUserToken(ctx)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, input.UserAccessToken, got.UserAccessToken)
		require.Equal(t, input.UserRefreshToken, got.UserRefreshToken)
		require.Equal(t, input.UserId, got.UserId)
		require.Equal(t, input.ExpiresAt, got.ExpiresAt)
		require.Equal(t, *input.RefreshTokenExpiresAt, *got.RefreshTokenExpiresAt)
	})

	t.Run("RemoveUserToken deletes user token", func(t *testing.T) {
		err := mgr.RemoveUserToken(ctx)
		require.NoError(t, err)

		got, err := mgr.GetUserToken(ctx)
		require.NoError(t, err)
		require.Nil(t, got)
	})
}

func TestUnitDatasetTokenCRUD(t *testing.T) {
	ctx := t.Context()
	store := tokenstore.NewMemory()
	mgr := tokens.New(tokens.ManagerOptions{Store: store})

	t.Run("Get missing dataset token returns nil", func(t *testing.T) {
		tok, err := mgr.GetDatasetToken(ctx, "ds1")
		require.NoError(t, err)
		require.Nil(t, tok)
	})

	t.Run("Set and Get stores and retrieves dataset tokens", func(t *testing.T) {
		dt := &tokens.DatasetToken{
			Dataset:     "ds1",
			AccessToken: "ds-access-token-xyz",
			ExpiresAt:   time.Now().Add(5 * time.Minute).UnixMilli(),
			Grants:      grants.Read | grants.Write,
			Solution:    "relevance",
		}

		err := mgr.SetDatasetToken(ctx, "ds1", dt)
		require.NoError(t, err)

		got, err := mgr.GetDatasetToken(ctx, "ds1")
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, dt.Dataset, got.Dataset)
		require.Equal(t, dt.AccessToken, got.AccessToken)
		require.Equal(t, dt.Grants, got.Grants)

		// Check Map
		m, err := mgr.GetDatasetTokensMap(ctx)
		require.NoError(t, err)
		require.Len(t, m, 1)
		require.Equal(t, dt.AccessToken, m["ds1"].AccessToken)
	})

	t.Run("RemoveDatasetToken deletes specific token", func(t *testing.T) {
		err := mgr.RemoveDatasetToken(ctx, "ds1")
		require.NoError(t, err)

		got, err := mgr.GetDatasetToken(ctx, "ds1")
		require.NoError(t, err)
		require.Nil(t, got)
	})
}

func TestUnitDatasetRefreshTokenCRUD(t *testing.T) {
	ctx := t.Context()
	store := tokenstore.NewMemory()
	mgr := tokens.New(tokens.ManagerOptions{Store: store})

	t.Run("Set and Get robot Ed25519 keys", func(t *testing.T) {
		err := mgr.SetDatasetRefreshToken(ctx, "ds_robot", "ed25519_priv_key_base64")
		require.NoError(t, err)

		key, err := mgr.GetDatasetRefreshToken(ctx, "ds_robot")
		require.NoError(t, err)
		require.Equal(t, "ed25519_priv_key_base64", key)

		// Remove
		err = mgr.RemoveDatasetRefreshToken(ctx, "ds_robot")
		require.NoError(t, err)

		key, err = mgr.GetDatasetRefreshToken(ctx, "ds_robot")
		require.NoError(t, err)
		require.Empty(t, key)
	})

	t.Run("Mutual exclusivity: Setting user token clears robot key and derived tokens", func(t *testing.T) {
		err := mgr.SetDatasetRefreshToken(ctx, "ds1", "robot-key")
		require.NoError(t, err)
		err = mgr.SetDatasetToken(ctx, "ds1", &tokens.DatasetToken{Dataset: "ds1", AccessToken: "tok1"})
		require.NoError(t, err)

		// Set user token
		err = mgr.SetUserToken(ctx, &tokens.UserToken{UserAccessToken: "user-jwt"})
		require.NoError(t, err)

		// Robot key must be evicted
		rfKey, err := mgr.GetDatasetRefreshToken(ctx, "ds1")
		require.NoError(t, err)
		require.Empty(t, rfKey)

		// Derived token must be evicted
		dt, err := mgr.GetDatasetToken(ctx, "ds1")
		require.NoError(t, err)
		require.Nil(t, dt)

		// User token must exist
		ut, err := mgr.GetUserToken(ctx)
		require.NoError(t, err)
		require.NotNil(t, ut)
		require.Equal(t, "user-jwt", ut.UserAccessToken)
	})

	t.Run("Mutual exclusivity: Setting robot key clears user token and derived tokens", func(t *testing.T) {
		err := mgr.SetUserToken(ctx, &tokens.UserToken{UserAccessToken: "user-jwt"})
		require.NoError(t, err)
		err = mgr.SetDatasetToken(ctx, "ds1", &tokens.DatasetToken{Dataset: "ds1", AccessToken: "tok1"})
		require.NoError(t, err)

		// Set robot key
		err = mgr.SetDatasetRefreshToken(ctx, "ds2", "robot-key-2")
		require.NoError(t, err)

		// User token must be evicted
		ut, err := mgr.GetUserToken(ctx)
		require.NoError(t, err)
		require.Nil(t, ut)

		// Derived token must be evicted
		dt, err := mgr.GetDatasetToken(ctx, "ds1")
		require.NoError(t, err)
		require.Nil(t, dt)

		// Robot key must exist
		rfKey, err := mgr.GetDatasetRefreshToken(ctx, "ds2")
		require.NoError(t, err)
		require.Equal(t, "robot-key-2", rfKey)
	})
}

func TestUnitTokenExpiration(t *testing.T) {
	mgr := tokens.New(tokens.ManagerOptions{})

	t.Run("Nil token is expired", func(t *testing.T) {
		require.True(t, mgr.IsAccessTokenExpired(nil))
		require.True(t, mgr.IsRefreshTokenExpired(nil))
		require.True(t, mgr.IsDatasetTokenExpired(nil))
	})

	t.Run("Token with 0 expires_at does not expire", func(t *testing.T) {
		ut := &tokens.UserToken{UserAccessToken: "token", ExpiresAt: 0}
		require.False(t, mgr.IsAccessTokenExpired(ut))

		dt := &tokens.DatasetToken{AccessToken: "token", ExpiresAt: 0}
		require.False(t, mgr.IsDatasetTokenExpired(dt))
	})

	t.Run("Token expiring within skew is expired", func(t *testing.T) {
		// Expires in 10 seconds, but default skew is 30 seconds -> considered expired
		ut := &tokens.UserToken{
			UserAccessToken: "token",
			ExpiresAt:       time.Now().Add(10 * time.Second).UnixMilli(),
		}
		require.True(t, mgr.IsAccessTokenExpired(ut, 30.0))

		// If skew is 5 seconds, it's NOT expired yet
		require.False(t, mgr.IsAccessTokenExpired(ut, 5.0))
	})

	t.Run("Refresh token expired timestamp", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Minute).UnixMilli()
		future := time.Now().Add(10 * time.Minute).UnixMilli()

		expired := &tokens.UserToken{
			UserRefreshToken:      "ref",
			RefreshTokenExpiresAt: &past,
		}
		require.True(t, mgr.IsRefreshTokenExpired(expired))

		valid := &tokens.UserToken{
			UserRefreshToken:      "ref",
			RefreshTokenExpiresAt: &future,
		}
		require.False(t, mgr.IsRefreshTokenExpired(valid))
	})
}

func TestUnitSilentUserTokenRenewal(t *testing.T) {
	ctx := t.Context()

	t.Run("Returns valid cached access token without API call", func(t *testing.T) {
		store := tokenstore.NewMemory()
		var apiCalls int32
		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, _ apiz.Api) (int, error) {
				atomic.AddInt32(&apiCalls, 1)
				return 200, nil
			},
		}

		mgr := tokens.New(tokens.ManagerOptions{Store: store, Apis: apis})
		_ = mgr.SetUserToken(ctx, &tokens.UserToken{
			UserAccessToken:  "valid-jwt",
			UserRefreshToken: "ref-xyz",
			ExpiresAt:        time.Now().Add(1 * time.Hour).UnixMilli(),
		})

		token, err := mgr.GetValidAccessToken(ctx)
		require.NoError(t, err)
		require.Equal(t, "valid-jwt", token)
		require.Equal(t, int32(0), atomic.LoadInt32(&apiCalls))
	})

	t.Run("Silently refreshes expired access token", func(t *testing.T) {
		store := tokenstore.NewMemory()
		var refreshCalls int32
		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if rApi, ok := api.(*userapiauthrefresh.Api); ok {
					atomic.AddInt32(&refreshCalls, 1)
					require.Equal(t, "ref-old", rApi.Req.Payload.RefreshToken)
					rApi.Resp = &userapiauthrefresh.Resp{
						Payload: &userapiauthrefresh.RespPayload{
							UserAccessToken:  "fresh-jwt-token",
							UserRefreshToken: "ref-new",
							UserId:           "usr_123",
							ExpiresAt:        time.Now().Add(1 * time.Hour).UnixMilli(),
						},
					}
					return 200, nil
				}
				return 404, fmt.Errorf("unexpected api: %s", api.ApiSpec().Id)
			},
		}

		mgr := tokens.New(tokens.ManagerOptions{Store: store, Apis: apis})
		_ = mgr.SetUserToken(ctx, &tokens.UserToken{
			UserAccessToken:  "expired-jwt",
			UserRefreshToken: "ref-old",
			ExpiresAt:        time.Now().Add(-10 * time.Minute).UnixMilli(),
		})

		freshToken, err := mgr.GetValidAccessToken(ctx)
		require.NoError(t, err)
		require.Equal(t, "fresh-jwt-token", freshToken)
		require.Equal(t, int32(1), atomic.LoadInt32(&refreshCalls))

		// Check updated token in store
		saved, err := mgr.GetUserToken(ctx)
		require.NoError(t, err)
		require.Equal(t, "fresh-jwt-token", saved.UserAccessToken)
		require.Equal(t, "ref-new", saved.UserRefreshToken)
	})

	t.Run("Double-checked locking deduplicates concurrent refreshes", func(t *testing.T) {
		store := tokenstore.NewMemory()
		var refreshCalls int32
		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if rApi, ok := api.(*userapiauthrefresh.Api); ok {
					atomic.AddInt32(&refreshCalls, 1)
					time.Sleep(20 * time.Millisecond) // simulate network latency
					rApi.Resp = &userapiauthrefresh.Resp{
						Payload: &userapiauthrefresh.RespPayload{
							UserAccessToken:  "fresh-concurrent-jwt",
							UserRefreshToken: "ref-new",
							ExpiresAt:        time.Now().Add(1 * time.Hour).UnixMilli(),
						},
					}
					return 200, nil
				}
				return 404, nil
			},
		}

		mgr := tokens.New(tokens.ManagerOptions{Store: store, Apis: apis})
		_ = mgr.SetUserToken(ctx, &tokens.UserToken{
			UserAccessToken:  "expired-jwt",
			UserRefreshToken: "ref-old",
			ExpiresAt:        time.Now().Add(-10 * time.Minute).UnixMilli(),
		})

		var wg sync.WaitGroup
		concurrency := 8
		tokensOut := make([]string, concurrency)
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				tok, err := mgr.GetValidAccessToken(context.Background())
				if err == nil {
					tokensOut[idx] = tok
				}
			}(i)
		}
		wg.Wait()

		// Refresh should have happened exactly ONCE
		require.Equal(t, int32(1), atomic.LoadInt32(&refreshCalls))
		for _, tok := range tokensOut {
			require.Equal(t, "fresh-concurrent-jwt", tok)
		}
	})

	t.Run("401 error clears user token", func(t *testing.T) {
		store := tokenstore.NewMemory()
		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, _ apiz.Api) (int, error) {
				return 401, fmt.Errorf("invalid refresh token")
			},
		}

		mgr := tokens.New(tokens.ManagerOptions{Store: store, Apis: apis})
		_ = mgr.SetUserToken(ctx, &tokens.UserToken{
			UserAccessToken:  "expired-jwt",
			UserRefreshToken: "revoked-ref",
			ExpiresAt:        time.Now().Add(-10 * time.Minute).UnixMilli(),
		})

		_, err := mgr.GetValidAccessToken(ctx)
		require.Error(t, err)

		// Store should be cleared
		saved, err := mgr.GetUserToken(ctx)
		require.NoError(t, err)
		require.Nil(t, saved)
	})
}

func TestUnitDatasetTokenDerivation(t *testing.T) {
	ctx := t.Context()

	t.Run("Derives dataset token using user access token", func(t *testing.T) {
		store := tokenstore.NewMemory()
		var derivationCalls int32
		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if dApi, ok := api.(*userapidatasettokenaccesscreate.Api); ok {
					atomic.AddInt32(&derivationCalls, 1)
					require.Equal(t, lakedsz.MustParse("my-dataset"), dApi.Req.Params.Dataset)
					require.Equal(t, "Bearer valid-user-token", dApi.Req.Headers.Authorization)
					dApi.Resp = &userapidatasettokenaccesscreate.Resp{
						Payload: &userapidatasettokenaccesscreate.RespPayload{
							AccessToken: "derived-ds-token-789",
							Dataset:     lakedsz.MustParse("my-dataset"),
							Grants:      grants.All(),
							ExpiresAt:   time.Now().Add(10 * time.Minute).UnixMilli(),
						},
					}
					return 200, nil
				}
				return 404, nil
			},
		}

		mgr := tokens.New(tokens.ManagerOptions{Store: store, Apis: apis})
		_ = mgr.SetUserToken(ctx, &tokens.UserToken{
			UserAccessToken: "valid-user-token",
			ExpiresAt:       time.Now().Add(1 * time.Hour).UnixMilli(),
		})

		token, err := mgr.GetValidDatasetAccessToken(ctx, "my-dataset")
		require.NoError(t, err)
		require.Equal(t, "derived-ds-token-789", token)
		require.Equal(t, int32(1), atomic.LoadInt32(&derivationCalls))

		// Subsequent call is served from cache
		cachedToken, err := mgr.GetValidDatasetAccessToken(ctx, "my-dataset")
		require.NoError(t, err)
		require.Equal(t, "derived-ds-token-789", cachedToken)
		require.Equal(t, int32(1), atomic.LoadInt32(&derivationCalls))
	})

	t.Run("Derives dataset token using machine robot Ed25519 refresh key", func(t *testing.T) {
		store := tokenstore.NewMemory()
		_, privKey, err := ed25519.GenerateKey(rand.Reader)
		require.NoError(t, err)

		privKeyB64 := base64.RawURLEncoding.EncodeToString(privKey)

		var robotCalls int32
		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if rApi, ok := api.(*userapidatasettokenrefresh.Api); ok {
					atomic.AddInt32(&robotCalls, 1)
					require.Equal(t, lakedsz.MustParse("robot-dataset"), rApi.Req.Payload.Dataset)
					require.NotEmpty(t, rApi.Req.Headers.SignatureHeader)
					require.Len(t, rApi.Req.Headers.SignatureHeader, 162) // 162 chars for 121B assertion

					rApi.Resp = &userapidatasettokenrefresh.Resp{
						Payload: &userapidatasettokenrefresh.RespPayload{
							AccessToken: "robot-ds-token-abc",
							Dataset:     lakedsz.MustParse("robot-dataset"),
							Grants:      grants.Write | grants.Read,
							ExpiresAt:   time.Now().Add(5 * time.Minute).UnixMilli(),
						},
					}
					return 200, nil
				}
				return 404, nil
			},
		}

		mgr := tokens.New(tokens.ManagerOptions{Store: store, Apis: apis})
		err = mgr.SetDatasetRefreshToken(ctx, "robot-dataset", privKeyB64)
		require.NoError(t, err)

		tok, err := mgr.GetValidDatasetAccessToken(ctx, "robot-dataset")
		require.NoError(t, err)
		require.Equal(t, "robot-ds-token-abc", tok)
		require.Equal(t, int32(1), atomic.LoadInt32(&robotCalls))
	})
}
