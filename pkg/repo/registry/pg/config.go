package pg

type RepoConfigPG struct {
	DBHost     string `mapstructure:"db_host"`
	DBPort     int    `mapstructure:"db_port"`
	DBName     string `mapstructure:"db_name"`
	DBSchema   string `mapstructure:"db_schema"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
}

func (*RepoConfigPG) IsTypeTagged() {}
func (*RepoConfigPG) IsRepoConfig() {}
