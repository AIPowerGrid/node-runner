package settings

import (
	"os"
)

type Settings struct {
	GOENV string
	Port  string
}

var env string
var settings = Settings{}

// Init reads env variables and assings settings object
func Init() {
	env = os.Getenv("GO_ENV")
	if env == "" {
		// fmt.Println("Warning: Setting development environment due to lack of GO_ENV value")
		env = "development"
	}
	LoadSettingsByEnv(env)

}

// LoadSettingsByEnv sets the global object
func LoadSettingsByEnv(env string) {
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "4444"
	}
	if env == "development" {

	} else if env == "production" {

	}
	settings = Settings{
		GOENV: env,
		Port:  port,
	}
}

// GetEnvironment returns env variable
func GetEnvironment() string {
	return env
}

// Get returns the settings
func Get() Settings {
	if settings == (Settings{}) {
		Init()
	}
	return settings
}
