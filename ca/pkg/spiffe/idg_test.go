package spiffe

import (
	"fmt"
	"testing"
)

func TestParseIDGIdentity(t *testing.T) {
	cases := []string{
		"spiffe://siteid/clusterid/appid",
		"spiffe://test/test/test",
	}
	for _, a := range cases {
		id, err := ParseIDGIdentity(a)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(id.String())
	}
}
