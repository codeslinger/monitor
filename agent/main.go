package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "os/signal"
  "time"
  "./linux"
  "../util"
)

var sampleInterval int
var collectorAddr string

func init() {
  flag.IntVar(&sampleInterval, "t", 10, "sampling interval in seconds")
  flag.StringVar(&collectorAddr, "c", "", "address of collector service")
}

func main() {
  flag.Parse()
  //runtime.GOMAXPROCS(runtime.NumCPU())
  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, os.Interrupt, os.Kill)
  log.Printf("agent started: sampling every %d seconds\n", sampleInterval)
  ticker := time.NewTicker(time.Duration(sampleInterval) * time.Second)
  metadata, err := linux.Metadata()
  if err != nil {
    log.Fatalf("could not determine metadata: %s", err)
  }
  for {
    select {
    case <-ticker.C:
      collectAndSubmit(time.Now(), metadata)
    case s := <-signalChan:
      log.Printf("caught signal %s: shutting down\n", s)
      return
    }
  }
}

func collectAndSubmit(now time.Time, meta []*util.Metadata) {
  samples := linux.StandardStats()
  fmt.Printf("---Sample:%d---\n", now.Unix())
  for _, m := range meta {
    fmt.Println(m)
  }
  for _, s := range samples {
    fmt.Println(s)
  }
}

