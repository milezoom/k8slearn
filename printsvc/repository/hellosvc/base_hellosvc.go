package hellosvc

import (
	"fmt"
	client "hellocln"
	"net/url"
	"printsvc/repository"
)

type helloSvcConfig struct {
	host string
	port int
}

func NewHelloSvcConfig(host string, port int) repository.RepoConf {
	return &helloSvcConfig{
		host: host,
		port: port,
	}
}

func (conf *helloSvcConfig) Init(r *repository.Repository) error {
	if conf.host == "" {
		return fmt.Errorf("upggateway repository host cannot be empty")
	}
	if conf.port == 0 {
		return fmt.Errorf("upggateway repository port cannot be empty")
	}
	_, err := url.Parse(fmt.Sprintf("%s:%d", conf.host, conf.port))
	if err != nil {
		return err
	}
	cli, err := client.NewImplementor(client.ImplementorConfig{
		Host: conf.host,
		Port: conf.port,
	})
	if err != nil {
		return err
	}
	r.HelloSvc = &HelloSvcClient{cli}
	return nil
}

func (conf *helloSvcConfig) GetRepoName() string {
	return "HelloSvc"
}
