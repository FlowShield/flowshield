package keyprovider

import (
	"fmt"
	"github.com/flowshield/flowshield/ca/pkg/keygen"
	"github.com/flowshield/flowshield/ca/pkg/spiffe"
	"testing"
)

var (
	testID = &spiffe.IDGIdentity{
		SiteID:    "test_site",
		ClusterID: "test_cluster",
		UniqueID:  "unique_id",
	}
)

func TestXKeyProvider_Generate(t *testing.T) {
	kp, _ := NewXKeyProvider(&spiffe.IDGIdentity{})
	err := kp.Generate(string(keygen.EcdsaSigAlg), 0)
	if err != nil {
		t.Error(err)
	}
	err = kp.Generate(string(keygen.RsaSigAlg), 0)
	if err != nil {
		t.Error(err)
	}
}

func TestXKeyProvider_CertificateRequest(t *testing.T) {
	kp, _ := NewXKeyProvider(testID)
	err := kp.Generate(string(keygen.EcdsaSigAlg), 0)
	if err != nil {
		t.Error(err)
	}
	csr, err := kp.CertificateRequest(nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("----------------------------\n",
		string(csr))
}
