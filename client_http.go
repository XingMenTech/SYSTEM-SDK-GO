package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/XingMenTech/SYSTEM-SDK-GO/pb"
	"github.com/sirupsen/logrus"
)

const (
	LoginPath              = "/admin/api/login"
	LogoutPath             = "/admin/api/logout"
	CheckTokenPath         = "/admin/api/checkToken"
	GetUserRolesPath       = "/admin/api/getUserRoles"
	GetUserPermissionsPath = "/admin/api/getUserPermissions"
	UserMenuPath           = "/admin/api/getUserMenu"
	RoleMenu               = "/admin/api/getRoleMenu"
	UserList               = "/admin/api/getUserList"
)

func newHttpClient(config *Config, log *logrus.Entry) Client {
	if log == nil {
		log = logrus.WithField("model", "HttpClient")
		log.Level = logrus.DebugLevel
	}

	return &httpClient{
		apiUrl: config.ApiUrl,
		impl: impl{
			accessId: config.AccessId,
			aes:      NewAES([]byte(config.AccessId), []byte(config.AccessKey)),
			log:      log,
		},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
func (c *httpClient) Login(userName, password string) (data *pb.AdminUser, err error) {
	req := &pb.LoginReq{
		Header: &pb.RpcHeader{
			Uuid:  "33333333333",
			Token: "",
		},
		UserName: userName,
		Password: password,
	}
	err = c.doRequest(LoginPath, http.MethodPost, req, &data)

	return
}
func (c *httpClient) Logout(token string) (err error) {
	req := &pb.RpcHeader{Token: token}
	err = c.doRequest(LogoutPath, http.MethodPost, req, "")
	return
}

func (c *httpClient) CheckToken(token string) (data *pb.AdminUser, err error) {
	req := &pb.RpcHeader{Token: token}
	err = c.doRequest(CheckTokenPath, http.MethodPost, req, &data)

	return
}

func (c *httpClient) GetUserRoles(token string) (data []*pb.AdminRole, err error) {
	req := &pb.RpcHeader{Token: token}
	err = c.doRequest(GetUserRolesPath, http.MethodPost, req, &data)
	return
}

func (c *httpClient) GetUserPermissions(token string) (data []string, err error) {
	req := &pb.RpcHeader{Token: token}
	err = c.doRequest(GetUserPermissionsPath, http.MethodPost, req, &data)
	return
}

func (c *httpClient) GetUserMenu(token string) (data []*pb.MenuTreeModel, err error) {

	req := &pb.RpcHeader{Token: token}
	err = c.doRequest(UserMenuPath, http.MethodPost, req, &data)

	return
}

func (c *httpClient) GetRoleMenu(token string, roleId int64) (data []*pb.MenuTreeModel, err error) {

	req := &pb.UserRoleMenuReq{
		Header: &pb.RpcHeader{Token: token},
		RoleId: roleId,
	}
	err = c.doRequest(RoleMenu, http.MethodPost, req, &data)

	return
}

func (c *httpClient) UserList(name, loginName string, roleId int64) (data *pb.UserListResp, err error) {

	req := &pb.UserListReq{
		Header: &pb.RpcHeader{
			Uuid:  "33333333",
			Token: "",
		},
		RoleId:    roleId,
		UserName:  name,
		LoginName: loginName,
		JobNum:    0,
	}

	err = c.doRequest(UserList, http.MethodPost, req, &data)
	return
}

func (c *httpClient) doRequest(path, method string, params interface{}, result interface{}) error {
	url := c.apiUrl + path

	requestParam := c.encrypt(params)
	reqParam, _ := json.Marshal(requestParam)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqParam))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println("pay center http request timeout")
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("network error: (%d)", resp.StatusCode)
		c.log.Error(msg)
		return errors.New(msg)
	}
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("response body read error")
		return errors.New("response body read error")
	}
	c.log.Debug("解码前响应数据：", string(bodyByte))

	var res *pb.RpcResp
	err = json.Unmarshal(bodyByte, &res)
	if err != nil {
		c.log.Error("response body unmarshal error")
		return errors.New("response body unmarshal error")
	}
	if res.Code != http.StatusOK {
		c.log.Error(res.Message)
		return errors.New(res.Message)
	}
	if res.Data == "" {
		return nil
	}

	decrypt := c.decrypt([]byte(res.Data))

	c.log.Debug("解码后响应数据：", decrypt)

	err = json.Unmarshal([]byte(decrypt), result)
	if err != nil {
		c.log.Error("response body unmarshal error")
		return err
	}
	return nil
}
