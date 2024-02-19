package main

import (
	"fmt"

	"github.com/meanii/sshm/config"
	"github.com/meanii/sshm/ssh"
)

func main() {
	config.InitConfig()
	defaultSSh := ssh.DefaultSSH(
		ssh.DeafultShhParam{
			Host: "meanii.dev",
			User: "root",
			Port: 22,
		},
	)
	client, err := defaultSSh.Dial()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	fmt.Printf("Connected to %s\n", defaultSSh.Host)
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	output, err := session.CombinedOutput("ls -l")
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}
