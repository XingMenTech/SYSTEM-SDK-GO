package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/XingMenTech/SYSTEM-SDK-GO/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GrpcClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	fmt.Printf("Starting RPC %s \n", method)             // 在调用之前记录日志
	err := invoker(ctx, method, req, reply, cc, opts...) // 调用原始方法
	if err != nil {                                      // 处理错误和在调用之后记录日志
		fmt.Printf("Error in RPC %s: %v \n", method, err)
		return err
	}
	fmt.Printf("Finished RPC %s \n", method) // 在调用之后记录日志
	return nil
}

// newGrpcClient 创建一个新的gRPC客户端
func newGrpcClient(config *Config, log *logrus.Entry) *grpcClient {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(GrpcClientInterceptor),
	}

	conn, err := grpc.NewClient(config.ApiUrl, opts...)

	if err != nil {
		return nil
	}

	return &grpcClient{
		impl: impl{
			accessId: config.AccessId,
			aes:      NewAES([]byte(config.AccessId), []byte(config.AccessKey)),
			log:      log,
		},
		conn:   conn,
		client: pb.NewSystemServiceClient(conn),
	}
}

// Close 关闭gRPC连接
func (c *grpcClient) Close() error {
	return c.conn.Close()
}

// Login 登录
func (c *grpcClient) Login(name, password string) (*pb.AdminUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.LoginReq{
		Header: &pb.RpcHeader{
			Uuid:  "33333333",
			Token: "",
		},
		UserName: name,
		Password: password,
	}

	rpcParam := c.encrypt(req)
	resp, err := c.client.Login(ctx, rpcParam)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}

	respData := c.decrypt([]byte(resp.Data))
	var user pb.AdminUser
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *grpcClient) Logout(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.RpcHeader{
		Uuid:  "33333333",
		Token: token,
	}
	rpcParam := c.encrypt(req)
	resp, err := c.client.Logout(ctx, rpcParam)
	if err != nil {
		return err
	}
	if resp.Code != http.StatusOK {
		return errors.New(resp.Message)
	}
	return nil
}

func (c *grpcClient) CheckToken(token string) (*pb.AdminUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.RpcHeader{
		Uuid:  "33333333",
		Token: token,
	}
	rpcParam := c.encrypt(req)
	resp, err := c.client.CheckToken(ctx, rpcParam)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	respData := c.decrypt([]byte(resp.Data))
	var user pb.AdminUser
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *grpcClient) GetUserRoles(token string) ([]*pb.AdminRole, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.RpcHeader{
		Uuid:  "33333333",
		Token: token,
	}
	rpcParam := c.encrypt(req)
	resp, err := c.client.GetUserRoles(ctx, rpcParam)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	respData := c.decrypt([]byte(resp.Data))
	var user []*pb.AdminRole
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *grpcClient) GetUserPermissions(token string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.RpcHeader{
		Uuid:  "33333333",
		Token: token,
	}
	rpcParam := c.encrypt(req)
	resp, err := c.client.GetUserPermissions(ctx, rpcParam)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	respData := c.decrypt([]byte(resp.Data))
	var user []string
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *grpcClient) GetUserMenu(token string) ([]*pb.MenuTreeModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.RpcHeader{
		Uuid:  "33333333",
		Token: token,
	}
	rpcParam := c.encrypt(req)
	resp, err := c.client.GetUserMenu(ctx, rpcParam)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	respData := c.decrypt([]byte(resp.Data))
	var menu []*pb.MenuTreeModel
	if err := json.Unmarshal([]byte(respData), &menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (c *grpcClient) GetRoleMenu(token string, roleId int64) ([]*pb.MenuTreeModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pb.UserRoleMenuReq{
		Header: &pb.RpcHeader{Token: token},
		RoleId: roleId,
	}

	rpcParam := c.encrypt(req)
	resp, err := c.client.GetRoleMenu(ctx, rpcParam)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}
	respData := c.decrypt([]byte(resp.Data))
	var menu []*pb.MenuTreeModel
	if err := json.Unmarshal([]byte(respData), &menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (c *grpcClient) UserList(name, loginName string, roleId int64) (*pb.UserListResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

	rpcParam := c.encrypt(req)
	resp, err := c.client.UserList(ctx, rpcParam)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Message)
	}

	respData := c.decrypt([]byte(resp.Data))
	var user pb.UserListResp
	if err := json.Unmarshal([]byte(respData), &user); err != nil {
		return nil, err
	}
	return &user, nil
}
