package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OAuthConfig struct {
	GoogleLoginConfig oauth2.Config
	GitHubLoginConfig oauth2.Config
}

var AppConfig OAuthConfig

func GoogleConfig() oauth2.Config {
	AppConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  viper.GetString("oauth2.google.redirect_url"),
		ClientID:     viper.GetString("oauth2.google.client_id"),
		ClientSecret: viper.GetString("oauth2.google.client_secret"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig
}

func GithubConfig() oauth2.Config {
	AppConfig.GitHubLoginConfig = oauth2.Config{
		RedirectURL:  viper.GetString("oauth2.github.redirect_url"),
		ClientID:     viper.GetString("oauth2.github.client_id"),
		ClientSecret: viper.GetString("oauth2.github.client_secret"),
		Scopes:       []string{"user", "repo"},
		Endpoint:     github.Endpoint,
	}

	return AppConfig.GitHubLoginConfig
}
