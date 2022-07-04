package singleca

import "github.com/pkg/errors"

var errBadSigner = errors.New("signer not initialized")
var errNoCertDBConfigured = errors.New("cert database not configured (missing -database-config)")
