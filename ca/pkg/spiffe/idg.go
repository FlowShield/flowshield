package spiffe

import (
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"strings"
)

// IDG Identity
// be like "spiffe://siteid/clusterid/unique_id"
type IDGIdentity struct {
	SiteID    string `json:"site_id"`
	ClusterID string `json:"cluster_id"`
	UniqueID  string `json:"unique_id"`
}

func ParseIDGIdentity(s string) (*IDGIdentity, error) {
	id, err := spiffeid.FromString(s)
	if err != nil {
		return nil, err
	}
	split := strings.Split(strings.Trim(id.Path(), "/"), "/")
	var idi IDGIdentity
	idi.SiteID = id.TrustDomain().String()
	if len(split) > 0 {
		idi.ClusterID = split[0]
	}
	if len(split) > 1 {
		idi.UniqueID = split[1]
	}
	return &idi, nil
}

func (i IDGIdentity) SpiffeID() spiffeid.ID {
	id, _ := spiffeid.New(i.SiteID, i.ClusterID, i.UniqueID)
	return id
}

func (i IDGIdentity) String() string {
	return i.SpiffeID().String()
}
