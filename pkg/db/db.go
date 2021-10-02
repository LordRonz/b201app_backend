package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/lordronz/b201app_backend/pkg/types"
)

// Client is a custom db client
type Client struct {
	Client *gorm.DB
}

// Ping allows the db to be pinged.
func (c *Client) Ping() error {
	sqlDB, err := c.Client.DB()

	if err != nil {
		return err;
	}

	return sqlDB.Ping()
}

// Connect establishes a connection to the database and auto migrates the database schema
func (c *Client) Connect(dsn string) error {
	var err error
	// Create the database connection
	c.Client, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	// End the program with an error if it could not connect to the database
	if err != nil {
		return err
	}
	// c.autoMigrate()
	return nil
}

// autoMigrate creates the default database schema
func (c *Client) autoMigrate() {
	c.Client.AutoMigrate(&types.User{})
}
