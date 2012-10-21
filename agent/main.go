package main

import (
  "flag"
  "fmt"
  "os"
  "os/signal"
  "runtime"
  "time"
  "../util"
)

var sampleInterval int
var collectorAddr string

func init() {
  flag.IntVar(&sampleInterval, "sample", 10, "sampling interval in seconds")
  flag.StringVar(&collectorAddr, "collector", "", "address of collector service")
}

func main() {
  flag.Parse()
  if err := util.StartLogger(); err != nil {
    fmt.Println("could not start logger!")
    return
  }
  runtime.GOMAXPROCS(runtime.NumCPU())

  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, os.Interrupt, os.Kill)

  util.Log("agent started: sampling every %d seconds", sampleInterval)
  ticker := time.NewTicker(time.Duration(sampleInterval) * time.Second)
  for {
    select {
    case <-ticker.C:
      // TODO: do a sampling here
      break
    case s := <-signalChan:
      util.Log("caught signal %s: shutting down", s)
      return
    }
  }
}

