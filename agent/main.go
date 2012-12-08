package main

import (
  "flag"
  "log"
  "os"
  "os/signal"
  "runtime"
  "time"
  "./linux"
)

var sampleInterval int
var collectorAddr string

func init() {
  flag.IntVar(&sampleInterval, "t", 10, "sampling interval in seconds")
  flag.StringVar(&collectorAddr, "c", "", "address of collector service")
}

func main() {
  flag.Parse()
  runtime.GOMAXPROCS(runtime.NumCPU())
  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, os.Interrupt, os.Kill)
  log.Printf("agent started: sampling every %d seconds\n", sampleInterval)
  ticker := time.NewTicker(time.Duration(sampleInterval) * time.Second)
  for {
    select {
    case <-ticker.C:
      linux.StandardStats()
    case s := <-signalChan:
      log.Printf("caught signal %s: shutting down\n", s)
      return
    }
  }
}

