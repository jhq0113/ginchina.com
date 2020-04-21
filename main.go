package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := gin.New()
	router.Static("/web", "./web")

	server := &http.Server{
		Addr:           bootstrap.Config.Application.Addr,
		Handler:        router,
		ReadTimeout:    helper.Int64ToSecond(bootstrap.Config.Application.ReadTimeout),
		WriteTimeout:   helper.Int64ToSecond(bootstrap.Config.Application.WriteTimeout),
		MaxHeaderBytes: bootstrap.Config.Application.MaxHeaderBytes,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	//优雅退出
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-ch
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(cxt)
	if err != nil {
		fmt.Println(err)
	}
}
