package conf

type Env struct {
	Host       string `mapstructure:"HOST"`
	Port       string `mapstructure:"PORT"`
	Production bool   `mapstructure:"PRODUCTION"`

	DBRead struct {
		Host   string `mapstructure:"DB_READ_HOST"`
		Port   string `mapstructure:"DB_READ_PORT"`
		Name   string `mapstructure:"DB_READ_NAME"`
		User   string `mapstructure:"DB_READ_USER"`
		Pass   string `mapstructure:"DB_READ_PASS"`
		Schema string `mapstructure:"DB_READ_SCHEMA"`
	}

	DBWrite struct {
		Host   string `mapstructure:"DB_WRITE_HOST"`
		Port   string `mapstructure:"DB_WRITE_PORT"`
		Name   string `mapstructure:"DB_WRITE_NAME"`
		User   string `mapstructure:"DB_WRITE_USER"`
		Pass   string `mapstructure:"DB_WRITE_PASS"`
		Schema string `mapstructure:"DB_WRITE_SCHEMA"`
	}

	Kafka struct {
		HostPort      string `mapstructure:"KAFKA_HOST_PORT"`
		ConsumerGroup string `mapstructure:"KAFKA_CONSUMER_GROUP"`
	}

	Redis struct {
		HostPort string `mapstructure:"REDIS_HOST_PORT"`
		Pass     string `mapstructure:"REDIS_PASS"`
		DB       string `mapstructure:"REDIS_DB"`
		TimeOut  string `mapstructure:"REDIS_TIMEOUT"`
	}

	Elastic struct {
		Host string `mapstructure:"ELASTIC_HOST"`
		Port string `mapstructure:"ELASTIC_PORT"`
		User string `mapstructure:"ELASTIC_USER"`
		Pass string `mapstructure:"ELASTIC_PASS"`
	}

	KeyCloak struct {
		Host         string `mapstructure:"KEYCLOAK_HOST"`
		Port         string `mapstructure:"KEYCLOAK_PORT"`
		Realm        string `mapstructure:"KEYCLOAK_REALM"`
		User         string `mapstructure:"KEYCLOAK_USER"`
		Pass         string `mapstructure:"KEYCLOAK_PASS"`
		ClientId     string `mapstructure:"KEYCLOAK_CLIENT_ID"`
		ClientSecret string `mapstructure:"KEYCLOAK_CLIENT_SECRET"`
	}
}
