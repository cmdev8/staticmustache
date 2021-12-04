package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		panic(pwdErr)
	}

	inputDir := filepath.Join(pwd, "input")
	outputDir := filepath.Join(pwd, "out")
	layout := filepath.Join(pwd, "input", "layout.mustache")

	if err := build(inputDir, outputDir, layout); err != nil {
		panic(err)
	}

	go startWatcher(inputDir, outputDir, layout)
	go startWebServer(outputDir)
	select {}
}

func startWatcher(inputDir, outputDir, layout string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)

					if err := build(inputDir, outputDir, layout); err != nil {
						panic(err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func startWebServer(outputDir string) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Static("/", outputDir)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
