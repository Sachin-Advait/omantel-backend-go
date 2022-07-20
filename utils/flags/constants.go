package flags

// Flag constants
const (
	EnvKey                     = "env"
	EnvDefaultValue            = "dev"
	EnvUsage                   = "runtime environment such as dev, test, integration, pre-prod, prod"
	PortKey                    = "port"
	PortDefaultValue           = 8090
	PortUsage                  = "port for this service"
	BaseConfigPathKey          = "base-config-path"
	BaseConfigPathDefaultValue = "./resources"
	BaseConfigPathUsage        = "path to base folder that stores configurations"
)
