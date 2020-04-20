package services

import (
	"strings"
	"fmt"
	"os"
)

type ConfigurationService struct{}

func NewConfigurationService() ConfigurationService {
	return ConfigurationService{}
}

func (c *ConfigurationService) GetDasource() string {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DATABASE")

	clearDbUrl := os.Getenv("CLEARDB_DATABASE_URL")
	if len(clearDbUrl) > 0 {
		clearDbUrl = strings.Replace(clearDbUrl, "mysql://", "", 1)
		clearDbUrl = strings.Replace(clearDbUrl, "?reconnect=true", "", 1)
		parts := strings.Split(clearDbUrl, "@")
		userPass := strings.Split(parts[0], ":")
		user = userPass[0]
		pass = userPass[1]
		hostDbName := strings.Split(parts[1], "/")
		host = hostDbName[0]
		dbname = hostDbName[1]
		port = "3306"
	}

	return fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, dbname)
}

func (c *ConfigurationService) GetAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}

func (c *ConfigurationService) GetPort() string {
	return os.Getenv("PORT")
}

func (c *ConfigurationService) GetJwtSecret() string {
	return os.Getenv("JWT_SECRET")
}
