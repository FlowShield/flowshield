package middle

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/cloudslit/cloudslit/fullnode/pkg/confer"
)

func Session(name string, cfg *confer.Redis) gin.HandlerFunc {
	//return sessions.Sessions(name, sessions.NewCookieStore([]byte("secret")))
	store, _ := sessions.NewRedisStore(10, "tcp", cfg.Addr, "", []byte("secret"))
	return sessions.Sessions(name, store)
}
