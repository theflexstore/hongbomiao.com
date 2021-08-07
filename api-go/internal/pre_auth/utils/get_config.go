package utils

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	AppEnv                       string
	Port                         string
	OPAHost                      string
	OPAPort                      string
	DgraphHost                   string
	DgraphGRPCPort               string
	ElasticAPMServiceName        string
	ElasticAPMServerURL          string
	OpenCensusAgentHost          string
	OpenCensusAgentPort          string
	JaegerURL                    string
	JWTSecret                    string
	EnableOpenTelemetryStdoutLog string
}

func GetConfig() *Config {
	path := "config/pre_auth/"

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	_ = godotenv.Load(path + ".env." + appEnv + ".local")
	_ = godotenv.Load(path + ".env." + appEnv)

	return &Config{
		AppEnv:                       appEnv,
		Port:                         os.Getenv("PORT"),
		OPAHost:                      os.Getenv("OPA_HOST"),
		OPAPort:                      os.Getenv("OPA_PORT"),
		DgraphHost:                   os.Getenv("DGRAPH_HOST"),
		DgraphGRPCPort:               os.Getenv("DGRAPH_GRPC_PORT"),
		ElasticAPMServiceName:        os.Getenv("ELASTIC_APM_SERVICE_NAME"),
		ElasticAPMServerURL:          os.Getenv("ELASTIC_APM_SERVER_URL"),
		OpenCensusAgentHost:          os.Getenv("OPEN_CENSUS_AGENT_HOST"),
		OpenCensusAgentPort:          os.Getenv("OPEN_CENSUS_AGENT_PORT"),
		JaegerURL:                    os.Getenv("JAEGER_URL"),
		JWTSecret:                    os.Getenv("JWT_SECRET"),
		EnableOpenTelemetryStdoutLog: os.Getenv("ENABLE_OPEN_TELEMETRY_STDOUT_LOG"),
	}
}
