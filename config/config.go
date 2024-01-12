package config

import (
	"github/yogabagas/join-app/pkg/config"
	"log"
	"os"
)

var GlobalCfg *Config

type (
	Config struct {
		App                    App       `json:"app"`
		DB                     DB        `json:"db"`
		Cache                  Cache     `json:"cache"`
		Whitelist              Whitelist `json:"whitelist"`
		JWK                    JWK       `json:"jwk"`
		PasswordAlg            string    `json:"password_alg"`
		TokenExpiration        int       `json:"token_exp"`
		RefreshTokenExpiration int       `json:"refresh_token_exp"`
		Storage                Storage   `json:"storage"`
	}

	App struct {
		Name         string `json:"name"`
		Host         string `json:"host"`
		Port         string `json:"port"`
		ReadTimeout  int    `json:"read_timeout"`
		WriteTimeout int    `json:"write_timeout"`
		JWTSecret    string `json:"jwt_secret"`
	}

	DB struct {
		SQL struct {
			User     string `json:"user"`
			Password string `json:"password"`
			Host     string `json:"host"`
			Schema   string `json:"schema"`
		} `json:"sql"`
	}

	Cache struct {
		Redis struct {
			User     string `json:"user"`
			Password string `json:"password"`
			Host     string `json:"host"`
		} `json:"redis"`
	}
	Whitelist struct {
		APIs []API `json:"API"`
	}

	API struct {
		Endpoint string   `json:"endpoint"`
		Methods  []string `json:"methods"`
	}

	JWK struct {
		Size      int    `json:"size"`
		KeyID     string `json:"key_id"`
		Algorithm string `json:"alg"`
		Use       string `json:"use"`
		Expired   int    `json:"ttl_in_hours"`
	}

	Storage struct {
		Minio struct {
			Endpoint  string `json:"endpoint"`
			Bucket    string `json:"bucket"`
			AccessKey string `json:"access_key"`
			SecretKey string `json:"secret_key"`
			Region    string `json:"region"`
			Secure    bool   `json:"secure"`
			Token     string `json:"token"`
		} `json:"minio"`
	}
)

func LoadConfig(path string) interface{} {

	env := os.Getenv("APP_ENV")

	log.Println("environment", env)

	if GlobalCfg == nil {
		err := config.ReadModuleConfig(
			&config.Cfg{
				Target: &GlobalCfg,
				Path:   path,
				Module: "config",
				Env:    env,
			})
		if err != nil {
			log.Fatalln("can't load file config", err)
		}
	}
	return GlobalCfg
}
