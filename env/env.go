package env

import (
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	GrpcPort              string
	HTTPPort              string
	TokenType             string
	JWTAccessTokenSecret  string
	JWTRefreshTokenSecret string
	AccessTokenLife       int64 //seconds
	RefreshTokenLife      int64 //seconds
	MaxLoginAttemps       int8
	LoginAgainAllowLife   int64 //seconds
	AllowedDomains        []string
	UserService           *UserServiceConfig
	Database              *DatabaseConfig
	SignupLinkLife        int
	JWTSignupTokenSecret  string
}

type UserServiceConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	PSQLConfig
}

type PSQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	SslMode  string
}

var Settings *Config

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	viper.BindEnv("GRPC_PORT")
	viper.BindEnv("HTTP_PORT")
	viper.BindEnv("TOKEN_TYPE")
	viper.BindEnv("JWT_ACCESS_TOKEN_SECRET")
	viper.BindEnv("JWT_REFRESH_TOKEN_SECRET")
	viper.BindEnv("ALLOWED_DOMAINS")
	viper.BindEnv("ACCESS_TOKEN_LIFE")  //1min
	viper.BindEnv("REFRESH_TOKEN_LIFE") //15days
	viper.BindEnv("MAX_LOGIN_ATTEMPT")
	viper.BindEnv("LOGIN_AGAIN_ALLOW_LIFE") //1hour

	viper.BindEnv("POSTGRES_HOST")
	viper.BindEnv("POSTGRES_PORT")
	viper.BindEnv("POSTGRES_USERNAME")
	viper.BindEnv("POSTGRES_PASSWORD")
	viper.BindEnv("POSTGRES_DBNAME")
	viper.BindEnv("POSTGRES_SSL_MODE")

	viper.BindEnv("USERS_SERVICE_HOST")
	viper.BindEnv("USERS_SERVICE_PORT")

	Settings = &Config{
		GrpcPort:              viper.GetString("GRPC_PORT"),
		HTTPPort:              viper.GetString("HTTP_PORT"),
		TokenType:             viper.GetString("TOKEN_TYPE"),
		JWTAccessTokenSecret:  viper.GetString("JWT_ACCESS_TOKEN_SECRET"),
		JWTRefreshTokenSecret: viper.GetString("JWT_REFRESH_TOKEN_SECRET"),
		AllowedDomains:        strings.Split(viper.GetString("ALLOWED_DOMAINS"), ","),
		AccessTokenLife:       viper.GetInt64("ACCESS_TOKEN_LIFE"),
		RefreshTokenLife:      viper.GetInt64("REFRESH_TOKEN_LIFE"),
		MaxLoginAttemps:       int8(viper.GetInt("MAX_LOGIN_ATTEMPT")),
		LoginAgainAllowLife:   viper.GetInt64("LOGIN_AGAIN_ALLOW_LIFE"),
		SignupLinkLife:        viper.GetInt("SIGNUP_LINK_LIFE"),
		JWTSignupTokenSecret:  viper.GetString("JWT_SIGNUP_TOKEN_SECRET"),
		Database: &DatabaseConfig{
			PSQLConfig{
				Host:     viper.GetString("POSTGRES_HOST"),
				Port:     viper.GetInt("POSTGRES_PORT"),
				Username: viper.GetString("POSTGRES_USERNAME"),
				Password: viper.GetString("POSTGRES_PASSWORD"),
				DbName:   viper.GetString("POSTGRES_DBNAME"),
				SslMode:  viper.GetString("POSTGRES_SSL_MODE"),
			},
		},
		UserService: &UserServiceConfig{
			Host: viper.GetString("USERS_SERVICE_HOST"),
			Port: viper.GetString("USERS_SERVICE_PORT"),
		},
	}
}
