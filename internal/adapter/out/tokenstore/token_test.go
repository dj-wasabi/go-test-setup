package tokenstore

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

var (
	logs              *slog.Logger
	rdb               *redis.Client
	mr                *miniredis.Miniredis
	ctx               context.Context
	serviceTokenstore out.PortStoreInterface
)

func setupRedis(t *testing.T) {
	logs = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	mr = miniredis.RunT(t)
	rdb = redis.NewClient(&redis.Options{
		Addr:     mr.Addr(),
		Password: "",
		DB:       0,
	})
}

func prepareTests() {
	logId := "myLogId"
	ctx = utils.NewContextWrapper(context.TODO(), logId).Build()
}

func TestAdd(t *testing.T) {
	setupRedis(t)
	prepareTests()

	serviceTokenstore = NewTokenstoreService(rdb, logs)
	tokenErr := serviceTokenstore.Add(ctx, "myusername", "randomtoken")

	assert.Nil(t, tokenErr)
	mr.CheckGet(t, "myusername", "randomtoken")

}

func TestGet(t *testing.T) {
	setupRedis(t)
	prepareTests()

	_ = mr.Set("myusername", "randomtoken")

	serviceTokenstore = NewTokenstoreService(rdb, logs)
	token, tokenErr := serviceTokenstore.Get(ctx, "myusername")

	assert.Nil(t, tokenErr)
	assert.Equal(t, token, "randomtoken")

}
