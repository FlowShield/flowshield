package test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestKey(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("8828b5d74cdfa86ae17b11d2df83f627a888fab3b86a139c6d442ef7d0e9dd76")
	if err != nil {
		t.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println(address.String())
}
