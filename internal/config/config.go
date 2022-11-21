package config

import (
	"fmt"
	"strings"
)

// AppConfig application configuration based on input flags
type AppConfig struct {
	FilePath     string
	Verbose      bool
	DbConfig     *DbConfig
	WorkersCount int
}

// DbConfig database configuration
type DbConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

// NewDbConfig creates a DbConfig given a dbUrl
// expected format "serverUrl:5432,database,user,password"
func NewDbConfig(dbUrl string) *DbConfig {
	dbSegments := strings.Split(dbUrl, ",")
	if len(dbSegments) != 4 {
		return nil
	}
	server := strings.Split(dbSegments[0], ":")
	host := server[0]
	port := "5432"
	if len(server) == 2 {
		port = server[1]
	}
	return &DbConfig{Host: host, Port: port, Database: dbSegments[1], User: dbSegments[2], Password: dbSegments[3]}
}

// DBConfig String is a format required by the sql.Open when we try to establish a connection
// e.g. db, err := sql.Open("postgres", dbConfig.String())
func (d *DbConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Database,
	)
}
