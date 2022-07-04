package caclient

import (
	"time"

	"github.com/ztalab/cfssl/transport/roots"
	"go.uber.org/zap"
)

// RotateController ...
type RotateController struct {
	transport   *Transport
	rotateAfter time.Duration
	logger      *zap.SugaredLogger
}

// Run ...
func (rc *RotateController) Run() {
	log := rc.logger
	ticker := time.NewTicker(60 * time.Minute)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			// Automatically update certificates
			err := rc.transport.AutoUpdate()
			if err != nil {
				log.Errorf("Certificate rotation failed: %v", err)
			}
			rc.AddCert()
		}
	}
}

func (rc *RotateController) AddCert() {
	log := rc.logger
	store, err := roots.New(rc.transport.Identity.Roots)
	if err != nil {
		log.Errorf("Failed to get roots: %v", err)
		return
	}
	rc.transport.TrustStore.AddCerts(store.Certificates())

	if len(rc.transport.Identity.ClientRoots) > 0 {
		store, err = roots.New(rc.transport.Identity.ClientRoots)
		if err != nil {
			log.Errorf("Failed to get client roots: %v", err)
			return
		}
		rc.transport.ClientTrustStore.AddCerts(store.Certificates())
	}
	return
}
