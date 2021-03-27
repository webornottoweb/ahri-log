package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/joho/godotenv"
	lj "github.com/webornottoweb/jenv/pkg/json"
)

// AuthConfig represents structure with authentication params to get access to ssh servers
type AuthConfig struct {
	User lj.JEnvString
	Key  struct {
		Path     lj.JEnvString
		Password lj.JEnvString
	}
}

// EndpointsConfig represents configuration for all servers from which command will be executed
type EndpointsConfig struct {
	Servers []EndpointServer
}

// Endpoint server represents host:port pair
type EndpointServer struct {
	Host lj.JEnvString
	Port lj.JEnvInt
}

// ColorsConfig represents configuration for text output coloring
type ColorsConfig struct {
	Out   lj.JEnvString
	Error lj.JEnvString
}

// AuthConfig represents ssh auth config instance
var Auth *AuthConfig

// Endpoints represents EndpointsConfig instance
var Endpoints *EndpointsConfig

// Colors represents ColorsConfig instance
var Colors *ColorsConfig

func init() {
	if Endpoints != nil && Colors != nil {
		return
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	file, err := ioutil.ReadFile("configs/auth.json")
	if err != nil {
		log.Fatal("Error loading auth.json file")
		return
	}

	var auth AuthConfig
	err = json.Unmarshal([]byte(file), &auth)
	if err != nil {
		log.Fatal("Error filling endpoints")
		return
	}

	Auth = &auth

	file, err = ioutil.ReadFile("configs/endpoints.json")
	if err != nil {
		log.Fatal("Error loading endpoints.json file")
		return
	}

	var ep EndpointsConfig

	err = json.Unmarshal([]byte(file), &ep)
	if err != nil {
		log.Fatal("Error filling endpoints")
		return
	}

	Endpoints = &ep

	file, err = ioutil.ReadFile("configs/colors.json")
	if err != nil {
		log.Fatal("Error loading colors.json file")
		return
	}

	var cl ColorsConfig

	err = json.Unmarshal([]byte(file), &cl)
	if err != nil {
		log.Fatal("Error filling colors")
		return
	}

	Colors = &cl
}
