package pg

import "fmt"

type config struct {
	userName string
	password string
	host     string
	portNum  int
	database string
	useTLS   bool
}

func (c *config) ToConnString() string {
	result := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?prefer_simple_protocol=true", c.userName, c.password, c.host, c.portNum, c.database)
	if c.useTLS {
		result += "&ssl=true"
	}
	return result
}

func NewConfig(userName string, password string, host string, portNum int, database string, useTLS bool) *config {
	return &config{
		userName: userName,
		password: password,
		host:     host,
		portNum:  portNum,
		database: database,
		useTLS:   useTLS,
	}
}
