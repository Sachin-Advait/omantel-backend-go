package config

import (
	"fmt"
)

type Configuration struct {
	Server   Server   `mapstructure:"server"`
	Redis    Redis    `mapstructure:"redis"`
	Database Database `mapstructure:"database"`
}

type Server struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	ApiTimeout        int64  `mapstructure:"api_timeout"`
	MaxServerRequests int    `mapstructure:"max_server_requests"`
}

type Redis struct {
	IsCacheEnabled       bool    `mapstructure:"is_cache_enabled"`
	Address              string  `mapstructure:"address"`
	Password             string  `mapstructure:"password"`
	AutoCompleteApiTTL   float64 `mapstructure:"autocomplete_api_ttl"`
	PlaceDetailsApiTTL   float64 `mapstructure:"place_details_api_ttl"`
	ReverseGeoCodeApiTTL float64 `mapstructure:"reverse_geo_code_api_ttl"`
	ReadTimeout          float64 `mapstructure:"read_timeout"`
	WriteTimeout         float64 `mapstructure:"write_timeout"`
	Mode                 string  `mapstructure:"mode"`
}

type Database struct {
	Host        string `mapstructure:"host"`
	Name        string `mapstructure:"name"`
	Port        int    `mapstructure:"port"`
	UserName    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	MaxIdelConn int    `mapstructure:"maxIdelConn"`
	MaxOpenConn int    `mapstructure:"maxOpenConn"`
}

func (server Server) GetAddressString() string {
	return fmt.Sprintf("%s:%d", server.Host, server.Port)
}
