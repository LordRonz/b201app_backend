package db

import (
	"github.com/lordronz/b201app_backend/pkg/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	pageSize = 10
)

// ClientInterface resembles a db interface to interact with an underlying db
type ClientInterface interface {
	Ping() error
	Connect(connectionString string) error
	GetUserByID(id int) *types.User
	SetUser(article *types.User) error
	GetUsers(pageID int) *types.UserList
}

// Client is a custom db client
type Client struct {
	Client *gorm.DB
}

// Ping allows the db to be pinged.
func (c *Client) Ping() error {
	sqlDB, err := c.Client.DB()

	if err != nil {
		return err
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
	c.autoMigrate()
	return nil
}

// autoMigrate creates the default database schema
func (c *Client) autoMigrate() {
	c.Client.AutoMigrate(&types.User{})
}

// GetUserByID queries an article from the database
func (c *Client) GetUserByID(id int) *types.User {
	user := &types.User{}

	c.Client.Where("id = ?", id).First(&user).Scan(user)

	return user
}

// SetUser writes a user to the database
func (c *Client) SetUser(user *types.User) error {
	// Upsert by trying to create and updating on conflict
	if err := c.Client.Create(&user).Error; err != nil {
		return c.Client.Model(&user).Where("id = ?", user.ID).Updates(&user).Error
	}
	return nil
}

// GetUsers returns all users from the database
func (c *Client) GetUsers(pageID int) *types.UserList {
	user := &types.UserList{}
	c.Client.Where("id >= ?", pageID).Order("id").Limit(pageSize + 1).Find(&user.Items)
	if len(user.Items) == pageSize+1 {
		user.NextPageID = user.Items[len(user.Items)-1].ID
		user.Items = user.Items[:pageSize]
	}
	return user
}
