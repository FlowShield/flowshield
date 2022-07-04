// Package ca Configuration class display
package ca

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/ztalab/cfssl/config"

	"github.com/cloudslit/cloudslit/ca/core"
)

type RoleProfile struct {
	Name           string        `json:"name"`
	Usages         []string      `json:"usages"`
	ExpiryString   string        `json:"expiry_string"`
	ExpiryDuration time.Duration `json:"expiry_duration" swaggertype:"string"`
	AuthKey        string        `json:"auth_key"`
	IsCa           bool          `json:"is_ca"`
}

// RoleProfiles Show environmental isolation status
//  No parameters are required
func (l *Logic) RoleProfiles() ([]RoleProfile, error) {
	cfg := core.Is.Config.Singleca.CfsslConfig
	if cfg == nil {
		l.logger.Error("cfssl config Empty")
		return nil, errors.New("cfssl config Empty")
	}
	roles := make([]RoleProfile, 0, len(cfg.Signing.Profiles)+1)

	parseRoleProfile := func(name string, profile *config.SigningProfile) RoleProfile {
		role := RoleProfile{
			Name:           strings.Title(name),
			Usages:         profile.Usage,
			ExpiryString:   profile.ExpiryString,
			ExpiryDuration: profile.Expiry,
			AuthKey:        profile.AuthKeyName,
			IsCa:           profile.CAConstraint.IsCA,
		}
		return role
	}

	for name, profile := range cfg.Signing.Profiles {
		roles = append(roles, parseRoleProfile(name, profile))
	}

	roles = append(roles, parseRoleProfile("default", cfg.Signing.Default))

	return roles, nil
}
