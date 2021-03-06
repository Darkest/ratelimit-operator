package ratelimiter

import (
	"github.com/stretchr/testify/assert"
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
	"strconv"
	"testing"
)

func Test_BuildNameForRedis(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build name for Redis", func(t *testing.T) {
		name := utils.BuildRandomString(3)
		a.Equal(name+"-redis", buildNameForRedis(name))
	})
}

func Test_BuildRedisUrl(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build Redis url", func(t *testing.T) {
		name := utils.BuildRandomString(3)
		a.Equal(buildNameForRedis(name)+":"+strconv.Itoa(int(constants.REDIS_PORT)), buildRedisUrl(name))
	})
}