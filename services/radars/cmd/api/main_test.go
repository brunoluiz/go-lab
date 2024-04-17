package main_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/core/storage/postgres"
	"github.com/brunoluiz/go-lab/core/tester/httpjson"
	"github.com/brunoluiz/go-lab/core/xgin"
	"github.com/brunoluiz/go-lab/services/radars/cmd/api/provider"
	"github.com/brunoluiz/go-lab/services/radars/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgres_test "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestGolden(t *testing.T) {
	t.Log("Starting test")
	ctx := context.Background()
	postgresContainer, err := postgres_test.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres_test.WithDatabase("postgres"),
		postgres_test.WithUsername("postgres"),
		postgres_test.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, postgresContainer.Terminate(ctx))
	}()

	pgAddress, err := postgresContainer.Endpoint(ctx, "")
	require.NoError(t, err)

	a := fx.Module("TestGolden",
		fx.Provide(func() *config.Config {
			return &config.Config{
				CommonConfig: app.CommonConfig{
					Env: "test",
				},
				HTTP: xgin.HTTPConfig{
					Address: "127.0.0.1",
					Port:    "6666",
				},
				DB: postgres.EnvConfig{
					DSN: fmt.Sprintf("postgres://postgres:postgres@%s/postgres?sslmode=disable", pgAddress),
				},
			}
		},
			func(c *config.Config) postgres.EnvConfig {
				return c.DB
			},
		),
		provider.InjectApp(),
	)

	steps := []httpjson.Step{
		{
			RequestTemplatePath: "./testdata/create_radar.json",
			PreHook: func(t *testing.T, kv *httpjson.KV) map[string]any {
				return map[string]any{
					"name": "Create first radar",
				}
			},
			PostHook: func(t *testing.T, kv *httpjson.KV, w *httptest.ResponseRecorder) {
				require.Equal(t, 201, w.Code)
				kv.Set("first.radar_id", gjson.Get(w.Body.String(), "data.radar.id").String())
			},
		},
		{
			RequestTemplatePath: "./testdata/create_radar.json",
			PreHook: func(t *testing.T, kv *httpjson.KV) map[string]any {
				return map[string]any{
					"name": "Create second radar",
				}
			},
			PostHook: func(t *testing.T, kv *httpjson.KV, w *httptest.ResponseRecorder) {
				require.Equal(t, 201, w.Code)
				kv.Set("second.radar_id", gjson.Get(w.Body.String(), "data.radar.id").String())
			},
		},
		{
			RequestTemplatePath: "./testdata/get_radar.json",
			PreHook: func(t *testing.T, kv *httpjson.KV) map[string]any {
				return map[string]any{
					"radar_id": kv.Get("first.radar_id").(string),
					"name":     "Get first radar",
				}
			},
			PostHook: func(t *testing.T, kv *httpjson.KV, w *httptest.ResponseRecorder) {
				require.Equal(t, 200, w.Code)
			},
		},
		{
			RequestTemplatePath: "./testdata/delete_radar.json",
			PreHook: func(t *testing.T, kv *httpjson.KV) map[string]any {
				return map[string]any{
					"radar_id": kv.Get("second.radar_id").(string),
					"name":     "Delete second radar",
				}
			},
			PostHook: func(t *testing.T, kv *httpjson.KV, w *httptest.ResponseRecorder) {
				require.Equal(t, 200, w.Code)
			},
		},
	}

	appTest := fxtest.New(t, a,
		fx.Replace(func(t *testing.T) []any {
			return []any{}
		}),
		fx.Invoke(func(r *gin.Engine) {
			httpjson.RunSequence(ctx, t, r, steps)
		}),
	)
	defer appTest.RequireStop()
	appTest.RequireStart()
}
