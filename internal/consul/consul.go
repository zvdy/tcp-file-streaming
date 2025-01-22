package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/zvdy/tcp-file-streaming/internal/utils"
)

func RegisterService(tcpPort, httpPort string) error {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	ip, err := utils.GetContainerIP()
	if err != nil {
		return err
	}

	tcpRegistration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("file-server-tcp-%s", ip),
		Name:    "file-server-tcp",
		Port:    8080,
		Tags:    []string{"tcp", "file", "server"},
		Address: ip,
	}

	httpRegistration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("file-server-http-%s", ip),
		Name:    "file-server-http",
		Port:    8081,
		Tags:    []string{"http", "file", "server"},
		Address: ip,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:8081/health", ip),
			Interval: "10s",
			Timeout:  "1s",
		},
	}

	if err := client.Agent().ServiceRegister(tcpRegistration); err != nil {
		return err
	}
	if err := client.Agent().ServiceRegister(httpRegistration); err != nil {
		return err
	}

	return nil
}
