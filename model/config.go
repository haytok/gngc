package model

type Config struct {
	GitHub GitHubConfig
	IFTTT  IFTTTConfig
}

type GitHubConfig struct {
	Token    string
	UserName string
}

type IFTTTConfig struct {
	EventName string
	Token     string
}
