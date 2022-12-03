package util

import (
	"crypto/sha1"
	"encoding/hex"
	"math"
	"strconv"
	"time"

	"github.com/flowshield/flowshield/ca/pkg/memorycacher"
	"github.com/gin-gonic/gin"
)

var MapCache *memorycacher.Cache

func init() {
	MapCache = memorycacher.New(2*time.Minute, 10*time.Minute, math.MaxInt64)
}

func GinRequestHash(g *gin.Context) string {
	url := g.Request.URL.String()
	body := g.Request.ContentLength

	sha := sha1.Sum([]byte(url + strconv.FormatInt(body, 10)))

	return hex.EncodeToString(sha[:])
}
