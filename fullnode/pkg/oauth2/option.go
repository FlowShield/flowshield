package oauth2

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
)

func GetOauth2RedirectURL(c *gin.Context, config *oauth2.Config) (redirectURL string, err error) {
	state := xid.New().String()
	redirectURL = config.AuthCodeURL(state, oauth2.ApprovalForce)
	session := sessions.Default(c)
	session.Set("state", state)
	err = session.Save()
	return
}
