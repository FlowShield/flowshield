package upperca

import (
	"strings"

	cfssl_client "github.com/ztalab/cfssl/api/client"

	"github.com/cloudslit/cloudslit/ca/ca/keymanager"
)

func ProxyRequest(f func(host string) error) error {
	return keymanager.GetKeeper().RootClient.DoWithRetry(func(remote *cfssl_client.AuthRemote) error {
		host := strings.TrimSuffix(remote.Hosts()[0], "/")
		return f(host)
	})
}
