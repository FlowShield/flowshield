package test

import (
	"github.com/ztalab/cfssl/hook"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	hook.ClientInsecureSkipVerify = true
	os.Chdir("./../../../")
	os.Setenv("IS_ENV", "test")
	//cli.Start(func(i *core.I) error {
	//	// CA Start
	//	go func() {
	//		err := singleca.Server()
	//		if err != nil {
	//			i.Logger.Fatal(err)
	//		}
	//	}()
	//	return nil
	//}, func(i *core.I) error {
	//	time.Sleep(2 * time.Second)
	//	if m.Run() > 0 {
	//		return errors.New("ERR")
	//	}
	//
	//	os.Exit(0)
	//	return nil
	//})
}
