package testutil

import (
	"github.com/go-redis/redismock/v9"

	"github.com/jr-dragon/dynamic_link/internal/data"
)

func NewTestingClients() (*data.Clients, error) {
	c, err := data.NewClients(data.Config{})
	if err != nil {
		return nil, err
	}

	{
		c.RDB, c.RDBMock = redismock.NewClientMock()
	}

	return c, err
}
