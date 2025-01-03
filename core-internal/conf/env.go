package conf

type Env struct {
	Port           string `mapstructure:"PORT"`
	Debugging      bool `mapstructure:"DEBUGGING"`
	ContextTimeOut int `mapstructure:"CONTEXT_TIMEOUT"`

	GRPCHost string `mapstructure:"GRPC_HOST"`
	GRPCPort string `mapstructure:"GRPC_PORT"`

	GRPCProductHost string `mapstructure:"GRPC_PRODUCT_HOST"`

	OtelHttpExporter string `mapstructure:"OTEL_HTTP_EXPORTER"`
	OtelGrpcExporter string `mapstructure:"OTEL_GRPC_EXPORTER"`

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
		Host      string `mapstructure:"KAFKA_HOST"`
		ConsumerGroup string `mapstructure:"KAFKA_CONSUMER_GROUP"`
	}

	Redis struct {
		Host string `mapstructure:"REDIS_HOST_PORT"`
		Pass     string `mapstructure:"REDIS_PASS"`
		DB       string `mapstructure:"REDIS_DB"`
		TimeOut  int64 `mapstructure:"REDIS_TIMEOUT"`
	}

	Elastic struct {
		Host string `mapstructure:"ELASTIC_HOST"`
		User string `mapstructure:"ELASTIC_USER"`
		Pass string `mapstructure:"ELASTIC_PASS"`
	}

	KeyCloak struct {
		Host         string `mapstructure:"KEYCLOAK_HOST"`
		Realm        string `mapstructure:"KEYCLOAK_REALM"`
		User         string `mapstructure:"KEYCLOAK_USER"`
		Pass         string `mapstructure:"KEYCLOAK_PASS"`
		ClientId     string `mapstructure:"KEYCLOAK_CLIENT_ID"`
		ClientSecret string `mapstructure:"KEYCLOAK_CLIENT_SECRET"`
	}
}
