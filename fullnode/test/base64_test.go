package test

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	data := base64.StdEncoding.EncodeToString([]byte("asdasd"))
	fmt.Println(data)
}
