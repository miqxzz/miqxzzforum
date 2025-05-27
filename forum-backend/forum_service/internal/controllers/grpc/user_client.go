package grpc

import (
	"context"
	"github.com/Engls/forum-project2/forum_service/internal/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client user.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &UserClient{
		conn:   conn,
		client: user.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) GetUsername(ctx context.Context, userID int) (string, error) {
	resp, err := c.client.GetUsername(ctx, &user.UserRequest{UserId: int32(userID)})
	if err != nil {
		log.Printf("Failed to get username: %v", err)
		return "", err
	}
	return resp.Username, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}
