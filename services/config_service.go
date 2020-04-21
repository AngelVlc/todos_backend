package services

import (
	"fmt"
	"strings"

	"github.com/AngelVlc/todos/providers"
)

type ConfigurationService struct {
	eg providers.EnvGetter
}

func NewConfigurationService(eg providers.EnvGetter) ConfigurationService {
	return ConfigurationService{eg}
}

func (c *ConfigurationService) GetDatasource() string {
	host := c.eg.Getenv("MYSQL_HOST")
	port := c.eg.Getenv("MYSQL_PORT")
	user := c.eg.Getenv("MYSQL_USER")
	pass := c.eg.Getenv("MYSQL_PASSWORD")
	dbname := c.eg.Getenv("MYSQL_DATABASE")

	clearDbUrl := c.eg.Getenv("CLEARDB_DATABASE_URL")
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
	return c.eg.Getenv("ADMIN_PASSWORD")
}

func (c *ConfigurationService) GetPort() string {
	return c.eg.Getenv("PORT")
}

func (c *ConfigurationService) GetJwtSecret() string {
	return c.eg.Getenv("JWT_SECRET")
}
