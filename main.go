package main

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

//go:embed dist/*
var FS embed.FS

func main() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	start := make(chan int)
	end := make(chan any)
	go Run(start, end)
	go func(start chan int, quit chan any) {
		port := <-start
		defer recoverFromError()
		ui, _ := lorca.New(fmt.Sprintf("http://127.0.0.1:%d/index.html", port), "", 1200, 800, "--disable-sync", " --disable-translate")
		defer ui.Close()
		quit <- <-ui.Done()
	}(start, end)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	select {
	case <-signalChannel:
		endWaiter.Done()
	case <-end:
		endWaiter.Done()
	}
	endWaiter.Wait()
}
func Run(start chan int, end chan any) {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()
	InitCors(router)
	staticFiles, _ := fs.Sub(FS, "dist")
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		reader, err := staticFiles.Open(path[1:])
		if err != nil {
			c.Status(http.StatusNotFound)
			log.Fatal(err)
			return
		} else {
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			contentType := mime.TypeByExtension(filepath.Ext(path))
			c.DataFromReader(http.StatusOK, stat.Size(), contentType, reader, nil)
		}
	})
	port := 27149
	start <- port
	runErr := router.Run(fmt.Sprintf(":%d", port))
	if runErr != nil {
		end <- runErr
		log.Fatal(runErr)
	}
}
func InitCors(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if origin == "http://127.0.0.1:3000" || origin == "http://localhost:3000" {
				return true
			} else {
				log.Printf("%v is now allowed", origin)
				return false
			}
		},
		MaxAge: 12 * time.Hour,
	}))
}
func recoverFromError() {
	if r := recover(); r != nil {
		fmt.Println("Recovering from panic:", r)
	}
}
