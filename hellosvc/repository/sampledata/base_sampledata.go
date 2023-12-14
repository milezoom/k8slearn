package	sampledata

import (
	"context"
	"fmt"
	"hellosvc/repository"
	"net/url"
)

type sampleDataConfig struct {
	host        string
	port        int
	checkHealth bool
}

func NewSampleDataConfig(host string, port int, checkHealth bool) repository.RepoConf {
	return &sampleDataConfig{
		host:		host,
		port:		port,
		checkHealth: checkHealth,
	}
}

func (conf *sampleDataConfig) Init(r *repository.Repository) error {
	if conf.host == "" {
		return fmt.Errorf("sampledata repository host cannot be empty")
	}
	if conf.port == 0 {
		return fmt.Errorf("sampledata repository port cannot be empty")
	}
	_, err := url.Parse(fmt.Sprintf("%s:%d", conf.host, conf.port))
	if err != nil {
		return err
	}
	// TODO: uncomment code block below with imported client
	// cli, err := client.NewImplementor(client.ImplementorConfig{
	// 	Host: conf.host,
	// 	Port: conf.port,
	// })
	// if err != nil {
	// 	return err
	// }
	sampledataClient := &SampleDataClient{
		// TODO: uncomment line below if implementor already defined
		// cli: cli,
	}
	if conf.checkHealth {
		if err := sampledataClient.HealthCheck(context.Background()); err != nil {
			return fmt.Errorf("sampledata repository failed to be created: healthcheck error, %s", err.Error())
		}
	}
	r.SampleData = sampledataClient
	return nil
}

func (conf *sampleDataConfig) GetRepoName() string {
	return "SampleData"
}