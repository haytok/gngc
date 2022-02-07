package model

type Config struct {
	GitHub GitHubConfig
	IFTTT  IFTTTConfig
}

type GitHubConfig struct {
	UserName string
	Token    string
}

type IFTTTConfig struct {
	EventName string
	Token     string
}
