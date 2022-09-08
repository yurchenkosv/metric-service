package main

import (
	"os"
	"time"
)

func Cleanup(mainLoop *time.Ticker, pushLoop *time.Ticker, mainLoopStop chan bool) {
	mainLoop.Stop()
	pushLoop.Stop()
	mainLoopStop <- true
	println("Program exit")
	os.Exit(0)
}
