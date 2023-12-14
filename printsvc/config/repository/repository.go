package repository

import (
	"log"

	"printsvc/config"
	"printsvc/repository"
	"printsvc/repository/hellosvc"
)

var repo *repository.Repository

func LoadRepository() {
	repoList, err := repository.NewRepository([]repository.RepoConf{
		hellosvc.NewHelloSvcConfig(
			config.GetConfig("hellosvc_host").GetString(),
			int(config.GetConfig("hellosvc_port").GetInt()),
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
