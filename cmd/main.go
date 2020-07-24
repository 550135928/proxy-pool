package main

import (
	"proxy-pool/api"
	"proxy-pool/check"
	"proxy-pool/config"
	"proxy-pool/databases"
	"proxy-pool/fetch"

	"github.com/jasonlvhit/gocron"
)

// SchedulerRun 定时拉取服务
func SchedulerRun(conf *config.Config, db *databases.ORM) {
	checker := check.NewChecker(db, conf)
	fetcher := fetch.NewFetcher(db, conf, checker)
	// checker.CheckAll()
	fetcher.FetchAllAndCheck()

	gocron.Every(conf.CheckProxy.CheckAllInterval).Seconds().Do(checker.CheckAll)
	// 定时拉取
	gocron.Every(conf.FetchProxy.FetchProxyInterval).Seconds().Do(fetcher.FetchAllAndCheck)

	// pending
	<-gocron.Start()
}

// APIRun 启动http服务
func APIRun(conf *config.Config, db *databases.ORM) {

	srv := api.NewService(db, conf)
	router := api.InitRouter(srv)
	router.Run(conf.HTTP.Port)
}

func main() {
	config := config.New()
	db := databases.New(config)
	defer db.Close()
	go APIRun(config, db)
	go SchedulerRun(config, db)
	select {}
}
