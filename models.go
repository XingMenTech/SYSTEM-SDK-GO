package client

import (
	"encoding/json"
	"net/http"

	"github.com/XingMenTech/SYSTEM-SDK-GO/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	ClientTypeHttp = "http"
	ClientTypeGrpc = "grpc"
)

type Config struct {
	ClientType string `yaml:"client_type" json:"clientType" comment:"客户端类型：http/grpc"`
	ApiUrl     string `yaml:"api_url" json:"apiUrl" comment:"API地址"`
	AccessId   string `yaml:"access_id" json:"accessId" comment:"accessId"`
	AccessKey  string `yaml:"access_key" json:"accessKey" comment:"accessKey"`
}

type Client interface {
	Login(userName, password string) (data *pb.AdminUser, err error)
	Logout(token string) (err error)
	CheckToken(token string) (data *pb.AdminUser, err error)
	GetUserRoles(token string) (data []*pb.AdminRole, err error)
	GetUserPermissions(token string) (data []string, err error)
	GetUserMenu(token string) (data []*pb.MenuTreeModel, err error)
	GetRoleMenu(token string, roleId int64) (data []*pb.MenuTreeModel, err error)
	UserList(name, loginName string, roleId int64) (data *pb.UserListResp, err error)
}

func NewClient(config *Config, log *logrus.Entry) Client {
	if config.ClientType == ClientTypeGrpc {
		return newGrpcClient(config, log)
	} else {
		return newHttpClient(config, log)
	}
}

type impl struct {
	aes      *AES
	accessId string
	log      *logrus.Entry
}

type grpcClient struct {
	impl
	conn   *grpc.ClientConn
	client pb.SystemServiceClient
}

type httpClient struct {
	impl
	apiUrl string
	client *http.Client
}

func (c *impl) decrypt(body []byte) string {
	decrypt, err := c.aes.Decrypt(body)
	if err != nil {
		return ""
	}

	return string(decrypt)
}

func (c *impl) encrypt(param interface{}) *pb.RpcParam {

	result := &pb.RpcParam{
		AppKey: c.accessId,
	}

	marshal, _ := json.Marshal(param)

	encrypt, err := c.aes.Encrypt(marshal)
	if err != nil {
		return nil
	}
	result.Data = encrypt
	return result
}
