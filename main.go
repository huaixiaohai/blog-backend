package main

import (
	"my/blog-backend/conf"
	_ "my/blog-backend/dao"
	"my/blog-backend/lib/log"
	"my/blog-backend/lib/redis"
	"my/blog-backend/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initRedis()
	router.Init()
	s := &http.Server{
		Addr:    ":" + conf.C.ServerPort,
		Handler: router.Router,
	}
	startServer(s)
	waitExit()
}

func startServer(s *http.Server) {
	go func() {
		log.Info("服务启动")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Starting server failed, err=%v", err)
			log.ErrorWithFields("Starting server failed", log.Fields{
				"err": err,
			})
		}
	}()
}

func waitExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	i := <-c
	log.InfoWithFields("Received interrupt signal, shutting down...", log.Fields{
		"signal": i,
	})
	log.Info("服务退出")
}

func initRedis() {
	cr := conf.C.Redis
	redis.RegisterRedisPool(cr.Dsn, cr.Password, cr.MaxIdle, cr.CatchDB)
}
