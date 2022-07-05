package curl

import (
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

var Client = resty.New()

func init() {
	if os.Getenv("IDG_RUNTIME") != "production" {
		Client.SetDebug(true)
	}
	Client.SetTimeout(time.Second * 5)
}
