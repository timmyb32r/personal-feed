package pg

import (
	"fmt"
	"net/url"
)

type RepoConfig struct {
	userName string
	password string
	host     string
	portNum  int
	database string
	useTLS   bool
}

func (c *RepoConfig) ToConnString() string {
	result := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?prefer_simple_protocol=true", c.userName, url.QueryEscape(c.password), c.host, c.portNum, c.database)
	if c.useTLS {
		result += "&ssl=true"
	}
	return result
}

func NewConfig(userName string, password string, host string, portNum int, database string, useTLS bool) *RepoConfig {
	return &RepoConfig{
		userName: userName,
		password: password,
		host:     host,
		portNum:  portNum,
		database: database,
		useTLS:   useTLS,
	}
}
