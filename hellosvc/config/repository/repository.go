package repository

import (
	"log"

	"hellosvc/repository"
	"hellosvc/config"
	"hellosvc/repository/sampledata"
)

var repo *repository.Repository

func LoadRepository() {
	repoList, err := repository.NewRepository([]repository.RepoConf{
		// TODO: add repository initialization here
		sampledata.NewSampleDataConfig(
			config.GetConfig("sampledata_host").GetString(),
			int(config.GetConfig("sampledata_port").GetInt()),
			config.GetConfig("check_healthy_repo").GetBool(),
		),
	})
	if err != nil {
		log.Fatalf("cannot initiate repository, with error: %v", err)
	}
	repo = repoList
}

func GetRepo() *repository.Repository {
	return repo
}
