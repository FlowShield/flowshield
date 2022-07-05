package certleaf

import (
	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"go.uber.org/zap"

	logic "github.com/cloudslit/cloudslit/ca/logic/certleaf"
)

type API struct {
	logger *zap.SugaredLogger
	logic  *logic.Logic
}

func NewAPI() *API {
	return &API{
		logger: logger.Named("api").SugaredLogger,
		logic:  logic.NewLogic(),
	}
}
