package docker

import (
	"context"

	dockercli "github.com/docker/docker/client"
)

// Client represents an active docker object
type Client struct {
	ctx context.Context
	*dockercli.Client
}

// New initializes a new Docker API client from env
func New(ctx context.Context) (*Client, error) {
	cli, err := dockercli.NewClientWithOpts(dockercli.FromEnv, dockercli.WithVersion("1.12"))
	if err != nil {
		return nil, err
	}
	if _, err = cli.ServerVersion(ctx); err != nil {
		return nil, err
	}
	return &Client{
		ctx:    ctx,
		Client: cli,
	}, err
}

// Close closes the docker client
func (c *Client) Close() {
	if c.Client != nil {
		_ = c.Client.Close()
	}
}
