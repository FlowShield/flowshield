package internal

import (
	"bufio"
	"context"
	"fmt"
	"github.com/cloudslit/cloudslit/client/internal/bll"
	"github.com/cloudslit/cloudslit/client/internal/config"
	"github.com/cloudslit/cloudslit/client/internal/schema"
	"github.com/cloudslit/cloudslit/client/pkg/errors"
	"github.com/cloudslit/cloudslit/client/pkg/logger"
	"github.com/cloudslit/cloudslit/client/pkg/util/json"
	"time"

	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// 登录状态
type State int

const (
	StateNotAuthenticated = State(iota)
	StateAuthenticating
	StateAuthenticated
)

type Up struct {
	UserDetail *schema.ControlUserDetail
	State      State
	UpCode     string
}

func NewUp() *Up {
	return &Up{
		State: StateNotAuthenticated,
	}
}

func InitClientServer(ctx context.Context) {
	up := NewUp()
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println("------------------------Interactive UI Start--------------------------")
	fmt.Println("----------------------------------------------------------------------")
	go func() {
		// Pre login
		err := up.preLogin()
		if err != nil {
			logger.Fatalf("login err:%v", err)
		}
		// Get client list
		client, err := up.printClients(ctx)
		if err != nil {
			logger.Fatalf("load client err:%v", err)
		}
		clientConfig, err := up.ParseW3sData(ctx, client)
		if err != nil {
			logger.Fatalf("parse client err:%v", err)
		}
		err = bll.NewClient().Listen(ctx, clientConfig)
		if err != nil {
			logger.Fatalf("start client err%v", err)
		}
		fmt.Println("########## start the client proxy #########")
	}()
}

func (a *Up) ParseW3sData(ctx context.Context, order *schema.ControlClient) (*schema.ClientConfig, error) {
	cfg := config.C.Web3
	// 解析配置
	tryCount := 0
retry:
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.W3S.Timeout)*time.Second)
	defer cancel()
	key := []byte(order.PeerId[len(order.PeerId)-8:])
	clientConfig, err := order.ToClientOrder(ctx, key)
	if err != nil {
		tryCount++
		logger.Warnf("get w3s data err:%v", err)
		if tryCount > cfg.W3S.RetryCount {
			return nil, err
		}
		goto retry
	}
	return clientConfig, nil
}

// printClients
func (a *Up) printClients(ctx context.Context) (*schema.ControlClient, error) {
	clients, err := GetControlClients()
	if err != nil {
		return nil, err
	}
	if len(clients) <= 0 {
		return nil, errors.NewWithStack("You have no available clients, Please go to " + config.C.App.ControlHost + " ,to place an order")
	}
	scanner := bufio.NewScanner(os.Stdin)
retry:
	fmt.Println("Please select one of the clients")
	for key, item := range clients {
		fmt.Println(key+1, " | ", item.Name)
	}
	fmt.Println("Please enter your client serial number:")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	if scanner.Text() == "" {
		fmt.Println("----------------------------------------------------------------------")
		fmt.Println("Error: Input errors, please re-enter")
		goto retry
	}
	index, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if index > len(clients) || index <= 0 {
		fmt.Println("Error: Input errors, please re-enter")
		goto retry
	}
	return clients[index-1], nil
}

func (a *Up) preLogin() error {
	// Get whether the device is logged in
	if config.C.Machine.Cookie != "" {
		// Validate cookies
		user, err := a.GetUserDetail()
		if err == nil {
			a.State = StateAuthenticated
			a.UserDetail = user
			return nil
		}
	}
	// Get login link
	upUrl, err := a.GetAuthUrl()
	if err != nil {
		return errors.WithStack(err)
	}
	// Output login connection
	fmt.Println("To authenticate, visit:")
	fmt.Println(upUrl.Data)
	a.State = StateAuthenticating
	a.UpCode = upUrl.GetCode()
	err = a.autoLogin()
	if err != nil {
		return errors.WithStack(err)
	}
	a.State = StateAuthenticated
	return nil
}

func (a *Up) autoLogin() error {
	result, err := a.GetLoginResult(110)
	if err != nil {
		return err
	} else {
		a.State = StateAuthenticated
		config.C.Machine.SetCookie(result.Data)
		_ = config.C.Machine.Write()
	}
	return nil
}

func (a *Up) GetUserDetail() (*schema.ControlUserDetail, error) {
	url := fmt.Sprintf("%s/api/v1/user/detail", config.C.App.ControlHost)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := httpControDo(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result schema.ControlUserDetail
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

// GetControlClients
func GetControlClients() (schema.ControlClients, error) {
	url := fmt.Sprintf("%s/api/v1/access/client?name=&page=1&limit_num=50&working=true", config.C.App.ControlHost)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := httpControDo(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result schema.ControlClientResult
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result.Code != 1001 {
		return nil, errors.NewWithStack(result.Message)
	}
	return result.Data.List, nil
}

func (a *Up) GetAuthUrl() (*schema.ControlMachineAuthResult, error) {
	url := fmt.Sprintf("%s/api/v1/controlplane/machine/%s", config.C.App.ControlHost, config.C.Machine.MachineId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := httpControDo(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var result schema.ControlMachineAuthResult
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result.Code != 1001 {
		return nil, errors.NewWithStack(result.Message)
	}
	return &result, nil
}

// GetLoginResult Get login result information
func (a *Up) GetLoginResult(timeout int) (*schema.ControlLoginResult, error) {
	url := fmt.Sprintf("%s/api/v1/controlplane/machine/auth/poll?timeout=%d&category=%s", config.C.App.ControlHost, timeout, a.UpCode)
	//fmt.Println(url)
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result schema.ControlLoginResult
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if result.Code != 1001 {
		return nil, errors.NewWithStack(result.Message)
	}
	return &result, nil
}

// httpControDo
func httpControDo(req *http.Request) ([]byte, error) {
	// Add request header
	req.Header.Add("Content-type", "application/json;charset=utf-8")
	// Add cookie
	cookie := &http.Cookie{
		Name:  "zta",
		Value: config.C.Machine.Cookie,
	}
	req.AddCookie(cookie)
	// Send request
	resp, err := config.Is.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.NewWithStack(resp.Status)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
