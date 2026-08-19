package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

func NewConfig(host, user, password, dbname, port string) *Config {
	return &Config{
		Host:     host,
		User:     user,
		Password: password,
		DBName:   dbname,
		Port:     port,
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.DBName, c.Port)
}

func Connect(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
