package ssh

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"

	"github.com/meanii/sshm/config"
)

type DeafultShhParam struct {
	Host string
	User string
	Port int
}

type SSH struct {
	Host            string
	User            string
	Port            int
	Config          *ssh_config.Config
	PrivateKey      ssh.Signer
	SshClientConfig *ssh.ClientConfig
	ShhClient       *ssh.Client
}

func DefaultSSH(sconfig DeafultShhParam) *SSH {
	sshc := &SSH{}
	if sconfig.Port != 0 {
		sshc.Port = sconfig.Port
	}
	if sconfig.User == "" {
		defaultUser := sconfig.User
		if defaultUser == "" {
			defaultUser = os.Getenv("USER")
		}
		sshc.User = defaultUser
	}

	ssh_config_file := config.Config.SSHConfigFile
	ssh_private_key_file := config.Config.PrivateKeyFile
	fmt.Printf("ssh private key file: %s\n", ssh_private_key_file)
	sshc.LoadConfig(ssh_config_file)
	sshc.LoadPrivateKey(ssh_private_key_file)
	sshc.User = config.Config.SSHUser
	sshc.Port = config.Config.SSHPort
	sshc.GetClientConfig()
	sshc.Host = sconfig.Host
	return sshc
}

func (s *SSH) LoadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	sshConfig, err := ssh_config.Decode(f)
	if err != nil {
		return err
	}
	s.Config = sshConfig
	return nil
}

func (s *SSH) LoadPrivateKey(path string) (ssh.Signer, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("couldn't open private key! Error: %s\n path: %s", err, path)
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Fatalf("couldn't load private key! Error: %s", err)
		return nil, err
	}
	s.PrivateKey = key
	return key, nil
}

func (s *SSH) GetClientConfig() *ssh.ClientConfig {
	if s.SshClientConfig != nil {
		return s.SshClientConfig
	}
	config := &ssh.ClientConfig{
		User:            s.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(s.PrivateKey)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	s.SshClientConfig = config
	return config
}

func (s *SSH) GetHostConfig(host string) string {
	hostname, err := s.Config.Get(host, "hostname")
	if err != nil {
		log.Fatalf("Error reading hostname from ssh config: %s", err)
	}
	s.Host = hostname
	return hostname
}

func (s *SSH) Dial() (*ssh.Client, error) {
	address := net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
	client, err := ssh.Dial("tcp", address, s.GetClientConfig())
	if err != nil {
		log.Fatalf("Error dialing ssh: %s", err)
		return nil, err
	}
	s.ShhClient = client
	return client, nil
}
