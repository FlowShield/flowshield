package caclient

import (
	"log"

	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"go.uber.org/zap"
)

func init() {
	f := zap.RedirectStdLog(logger.S().Desugar())
	f()
	log.SetFlags(log.LstdFlags)
}
