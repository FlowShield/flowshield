package caclient

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"golang.org/x/crypto/ocsp"
	"golang.org/x/sync/singleflight"

	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"github.com/pkg/errors"
)

var ocspBlockSign int64 = 0

var sg = new(singleflight.Group)

var ocspOpts = ocsp.RequestOptions{
	Hash: crypto.SHA1,
}

func SendOcspRequest(server string, req []byte, leaf, issuer *x509.Certificate) (*ocsp.Response, error) {
	if server == "" {
		server = leaf.OCSPServer[0]
	}
	var resp *http.Response
	var err error
	buf := bytes.NewBuffer(req)
	resp, err = httpClient.Post(server, "application/ocsp-request", buf)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.With("url", server, resp.Status, "body", string(body)).
			Warnf("Request error")
		return nil, fmt.Errorf("ocsp response err: %v", resp.Status)
	}

	switch {
	case bytes.Equal(body, ocsp.UnauthorizedErrorResponse):
		return nil, errors.New("OSCP unauthorized")
	case bytes.Equal(body, ocsp.MalformedRequestErrorResponse):
		return nil, errors.New("OSCP malformed")
	case bytes.Equal(body, ocsp.InternalErrorErrorResponse):
		return nil, errors.New("OSCP internal error")
	case bytes.Equal(body, ocsp.TryLaterErrorResponse):
		return nil, errors.New("OSCP try later")
	case bytes.Equal(body, ocsp.SigRequredErrorResponse):
		return nil, errors.New("OSCP signature required")
	}

	parsedOcspResp, err := ocsp.ParseResponseForCert(body, leaf, issuer)
	if err != nil {
		logger.With("body", string(body)).Errorf("ocsp Parsing error: %v", err)
		return nil, errors.Wrap(err, "ocsp Parsing error")
	}

	return parsedOcspResp, nil
}

// BlockOcspRequests Blocking OCSP requests will cause the MTLs handshake to fail
func BlockOcspRequests() {
	atomic.StoreInt64(&ocspBlockSign, 1)
}

// AllowOcspRequests
func AllowOcspRequests() {
	atomic.StoreInt64(&ocspBlockSign, 0)
}
