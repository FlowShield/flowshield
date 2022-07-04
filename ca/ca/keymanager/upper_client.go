package keymanager

import (
	"crypto/tls"
	"net/url"

	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"github.com/pkg/errors"
	"github.com/ztalab/cfssl/api/client"
	"github.com/ztalab/cfssl/auth"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/cloudslit/cloudslit/ca/core"
)

type UpperClients interface {
	DoWithRetry(f func(*client.AuthRemote) error) error
	AllClients() map[string]*client.AuthRemote
}

type upperClients struct {
	// ip to client
	clients map[string]*client.AuthRemote
	logger  *zap.SugaredLogger
}

func (uc *upperClients) DoWithRetry(f func(*client.AuthRemote) error) error {
	if len(uc.clients) == 0 {
		return errors.New("No clients available")
	}
	var errGroup error
	for _, upperClient := range uc.clients {
		err := f(upperClient)
		if err == nil {
			// success
			return nil
		}
		uc.logger.With("upper", upperClient.Hosts()).Warnf("upper ca Execution error: %s", err)
		multierr.AppendInto(&errGroup, err)
	}
	return errGroup
}

func (uc *upperClients) AllClients() map[string]*client.AuthRemote {
	return uc.clients
}

func NewUpperClients(adds []string) (UpperClients, error) {
	if len(adds) == 0 {
		return nil, errors.New("Upper CA Address configuration error")
	}
	ap, err := auth.New(core.Is.Config.Singleca.CfsslConfig.AuthKeys["intermediate"].Key, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Auth key Configuration error")
	}
	clients := make(map[string]*client.AuthRemote)
	for _, addr := range adds {
		upperAddr, err := url.Parse(addr)
		if err != nil {
			return nil, errors.Wrap(err, "Upper CA Address resolution error")
		}
		upperClient := client.NewAuthServer(addr, &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		}, ap)
		clients[upperAddr.Host] = upperClient
	}
	logger.Infof("Upper CA Client Quantity: %v", len(clients))
	return &upperClients{
		clients: clients,
		logger:  logger.Named("upperca").SugaredLogger,
	}, nil
}
