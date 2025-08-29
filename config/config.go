package config

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

type ApiRoute struct {
	Endpoint string
	Method   string
}

type Config struct {
	AccessTokenExpiresIn time.Duration
	AdminRoutes          []ApiRoute
	AllowHeaders         []string
	AllowMethods         []string
	AppEmail             string
	ApiUrl               string
	ClientUrl            string
	CloudinaryName       string
	CloudinaryKey        string
	CloudinarySecret     string
	CookieDomain         string
	CurrentUser          string
	CurrentUserId        string
	CurrentUserRole      string
	FacebookAuthId       string
	FacebookAuthSecret   string
	GinMode              string
	GoogleAuthId         string
	GoogleAuthSecret     string
	IsDevMode            bool
	JWTSecret            []byte
	MaxImageSize         int
	NonAuthRoutes        []ApiRoute
	Port                 string
	PostgresDbPort       string
	PostgresDbUrl        string
	RedisHost            string
	RedisPassword        string
	SmtpHost             string
	SmtpPassword         string
	SmtpPort             int
	SmtpUser             string
	Version              string
}

var AppConfig *Config

func InitializeConfig() {
	AppConfig = &Config{
		AccessTokenExpiresIn: time.Hour * 2,
		AppEmail:             os.Getenv("APP_EMAIL"),
		ApiUrl:               os.Getenv("API_URL"),
		ClientUrl:            os.Getenv("CLIENT_URL"),
		CloudinaryName:       os.Getenv("CLOUDINARY_NAME"),
		CloudinaryKey:        os.Getenv("CLOUDINARY_KEY"),
		CloudinarySecret:     os.Getenv("CLOUDINARY_SECRET"),
		CookieDomain:         os.Getenv("COOKIE_DOMAIN"),
		CurrentUser:          "CURRENT_USER",
		CurrentUserId:        "CURRENT_USER_ID",
		CurrentUserRole:      "CURRENT_USER_ROLE",
		FacebookAuthId:       os.Getenv("FACEBOOK_AUTH_ID"),
		FacebookAuthSecret:   os.Getenv("FACEBOOK_AUTH_SECRET"),
		GinMode:              os.Getenv("GIN_MODE"),
		GoogleAuthId:         os.Getenv("GOOGLE_AUTH_ID"),
		GoogleAuthSecret:     os.Getenv("GOOGLE_AUTH_SECRET"),
		IsDevMode:            os.Getenv("IS_DEV_MODE") == "true",
		JWTSecret:            []byte(os.Getenv("JWT_SECRET")),
		MaxImageSize:         1024 * 1024 * 5,
		Port:                 os.Getenv("PORT"),
		PostgresDbPort:       os.Getenv("POSTGRES_DB_PORT"),
		PostgresDbUrl:        os.Getenv("POSTGRES_DB_URL"),
		RedisHost:            os.Getenv("REDIS_URL"),
		RedisPassword:        os.Getenv("REDIS_PASSWORD"),
		SmtpHost:             os.Getenv("SMTP_HOST"),
		SmtpPassword:         os.Getenv("SMTP_PASSWORD"),
		SmtpPort:             func() int { port, _ := strconv.Atoi(os.Getenv("SMTP_PORT")); return port }(),
		SmtpUser:             os.Getenv("SMTP_USER"),
		Version:              os.Getenv("VERSION"),
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization", "X-RateLimit-Limit", "X-RateLimit-Reset",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		NonAuthRoutes: []ApiRoute{
			{Endpoint: "/api/v1", Method: http.MethodGet},
			{Endpoint: "/api/v1/health", Method: http.MethodGet},
			{Endpoint: "/api/v1/auth/signup", Method: http.MethodPost},
			{Endpoint: "/api/v1/auth/signin", Method: http.MethodPost},
			{Endpoint: "/api/v1/auth/google/callback", Method: http.MethodGet},
			{Endpoint: "/api/v1/auth/facebook/callback", Method: http.MethodGet},
			{Endpoint: "/api/v1/auth/refresh", Method: http.MethodPost},
			{Endpoint: "/api/v1/auth/verify-email", Method: http.MethodPost},
			{Endpoint: "/api/v1/auth/verify-phone", Method: http.MethodPost},
			{Endpoint: "/api/v1/auth/forgot-password", Method: http.MethodPost},
			{Endpoint: "/api/v1/auth/reset-password", Method: http.MethodPost},
			{Endpoint: "/api/v1/categories", Method: http.MethodGet},
			{Endpoint: "/api/v1/categories/:id", Method: http.MethodGet},
			{Endpoint: "/api/v1/products", Method: http.MethodGet},
			{Endpoint: "/api/v1/products/:id", Method: http.MethodGet},
			{Endpoint: "/api/v1/users/:id", Method: http.MethodGet},
			{Endpoint: "/api/v1/orders", Method: http.MethodPost},
			{Endpoint: "/api/v1/orders/:id", Method: http.MethodPut},
		},
		AdminRoutes: []ApiRoute{
			{Endpoint: "/api/v1/products", Method: http.MethodPost},
			{Endpoint: "/api/v1/products/:id", Method: http.MethodPut},
			{Endpoint: "/api/v1/products/:id", Method: http.MethodDelete},
			{Endpoint: "/api/v1/categories", Method: http.MethodPost},
			{Endpoint: "/api/v1/categories/:id", Method: http.MethodPut},
			{Endpoint: "/api/v1/categories/:id", Method: http.MethodDelete},
			{Endpoint: "/api/v1/orders", Method: http.MethodGet},
			{Endpoint: "/api/v1/users/:id", Method: http.MethodDelete},
		},
	}
}
