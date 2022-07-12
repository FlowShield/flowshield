package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

type Person struct {
	FirstName string `json:"first"`
	LastName  string `json:"last"`
}

func TestJson(t *testing.T) {
	const jsonStr = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	value := gjson.Get(jsonStr, "name")
	var person Person
	json.Unmarshal([]byte(value.String()), &person)
	fmt.Println(person)
}

func TestString(t *testing.T) {
	var str = "abcdefg"
	fmt.Println(str[len(str)-2:])
}
