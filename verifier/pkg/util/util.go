package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func FloatToString(Num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(Num, 'f', 2, 64)
}

func MD5Bytes(s []byte) string {
	ret := md5.Sum(s)
	return hex.EncodeToString(ret[:])
}

func GetCookieFromGin(ctx *gin.Context, key string) (value string) {
	value, _ = ctx.Cookie(key)
	return
}

func MD5(s string) string {
	return MD5Bytes([]byte(s))
}

func MD5File(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return MD5Bytes(data), nil
}

func ToInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

func StructToMap(item interface{}) (data map[string]interface{}, err error) {
	jsonByte, err := json.Marshal(item)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonByte, &data)
	if err != nil {
		return
	}
	return
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func JsonToMap(jsonStr string) (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func UtilIsEmpty(data string) bool {
	return strings.Trim(data, " ") == ""
}

func GetUniqueID() string {
	if uniqueID := os.Getenv("IDG_UNIQUEID"); len(uniqueID) > 0 {
		return uniqueID
	}
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s-%s", hostname, GetRandomString(6))
}

func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
