package caclient

import (
	"log"

	"github.com/cloudSlit/cloudslit/casdk/pkg/logger"
	"go.uber.org/zap"
)

func init() {
	f := zap.RedirectStdLog(logger.S().Desugar())
	f()
	log.SetFlags(log.LstdFlags)
}
