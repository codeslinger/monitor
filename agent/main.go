// Main routines for metric collection agent.
package main

import (
  "flag"
  "log"
  "os"
  "os/signal"
  "time"
  "./linux"
  "../util"
)

// Number of seconds between samples.
var sampleInterval int

// Network address of collector to which to submit gathered samples.
var collectorAddr string

func init() {
  flag.IntVar(&sampleInterval, "t", 10, "sampling interval in seconds")
  flag.StringVar(&collectorAddr, "c", "", "address of collector service")
}

func main() {
  flag.Parse()
  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, os.Interrupt, os.Kill)
  log.Printf("agent started: sampling every %d seconds\n", sampleInterval)
  ticker := time.NewTicker(time.Duration(sampleInterval) * time.Second)
  sampler := linux.NewStandardSampler(util.NewFileOpener(),
                                      util.NewConsoleSampleWriter())
  if err := sampler.Init(); err != nil {
    log.Fatalf("could not initialize sampler: %s\n", err)
  }
  for {
    select {
    case <-ticker.C:
      sampler.Sample()
    case s := <-signalChan:
      log.Printf("caught signal %s: shutting down\n", s)
      return
    }
  }
}

