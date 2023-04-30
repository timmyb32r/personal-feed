package pg

import (
	"fmt"
	"net/url"
)

type RepoConfigPG struct {
	Host     string `mapstructure:"db_host"`
	Port     int    `mapstructure:"db_port"`
	Name     string `mapstructure:"db_name"`
	Schema   string `mapstructure:"db_schema"`
	User     string `mapstructure:"db_user"`
	Password string `mapstructure:"db_password"`
	UseTLS   bool   `mapstructure:"use_tls"`
}

func (*RepoConfigPG) IsTypeTagged() {}
func (*RepoConfigPG) IsRepoConfig() {}

func (c *RepoConfigPG) ToConnString() string {
	result := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?prefer_simple_protocol=true", c.User, url.QueryEscape(c.Password), c.Host, c.Port, c.Name)
	if c.UseTLS {
		result += "&ssl=true"
	}
	return result
}
