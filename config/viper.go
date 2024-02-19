package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type SSHMConfig struct {
	SSHConfigFile  string `mapstructure:"ssh.config.file"`
	PrivateKeyFile string `mapstructure:"ssh.key.file"`
	PublicKeyFile  string `mapstructure:"ssh.pub.file"`
	SSHUser        string `mapstructure:"ssh.user"`
	SSHPort        int    `mapstructure:"ssh.port"`
	Timeout        int    `mapstructure:"ssh.timeout"`
}

var Config *SSHMConfig

func InitConfig() {
	Config = loadConfig()
}

func loadConfig() (config *SSHMConfig) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Default().Fatalf("Error getting user home directory: %s", err)
	}
	sshConfigPath := filepath.Join(home, ".ssh")
	sshmConfigPath := filepath.Join(home, ".config", "sshm")
	sshmConfigFile := filepath.Join(sshmConfigPath, "config")

	_ = os.MkdirAll(sshmConfigPath, 0755)
	f, err := os.Create(sshmConfigFile)
	if err != nil {
		log.Default().Fatalf("Error creating config file: %s", err)
	}
	defer f.Close()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(sshmConfigPath)

	viper.SetDefault("ssh.config.file", filepath.Join(sshConfigPath, "config"))
	viper.SetDefault("ssh.key.file", filepath.Join(sshConfigPath, "id_rsa"))
	viper.SetDefault("ssh.pub.file", filepath.Join(sshConfigPath, "id_rsa.pub"))
	viper.SetDefault("ssh.user", "root")
	viper.SetDefault("ssh.port", 22)
	viper.SetDefault("ssh.timeout", 30)

	viper.WriteConfig()
	err = viper.ReadInConfig()
	if err != nil {
		log.Default().Fatalf("Error reading config file: %s", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Default().Fatalf("Error unmarshaling config file: %s", err)
	}
	return config
}
