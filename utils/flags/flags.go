package flags

import (
	flag "github.com/spf13/pflag"
)

var (
	env            = flag.String(EnvKey, EnvDefaultValue, EnvUsage)
	port           = flag.Int(PortKey, PortDefaultValue, PortUsage)
	baseConfigPath = flag.String(BaseConfigPathKey, BaseConfigPathDefaultValue, BaseConfigPathUsage)
)

func init() {
	flag.Parse()
}

// Env is the runtime environment such as dev, test, integration, pre-prod, prod
func Env() string {
	return *env
}

// Port is the given port for this service where the process will be started
func Port() int {
	return *port
}

// BaseConfigPath is the path that holds the base configuration filee path
func BaseConfigPath() string {
	return *baseConfigPath
}
