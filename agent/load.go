package main

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "../util"
)

func LoadStats() ([]*util.Sample, error) {
  f, err := os.Open("/proc/loadavg")
  if err != nil {
    return nil, err
  }
  defer f.Close()
  rd := bufio.NewReader(f)
  samples := make([]*util.Sample, 0)
  for {
    line, err := rd.ReadString('\n')
    if err != nil {
      if err == io.EOF { break }
      return nil, err
    }
    s, err := parseLoadStats(line)
    if err != nil {
      return nil, err
    }
    if s != nil {
      for _, v := range s {
        samples = append(samples, v)
      }
    }
  }
  return samples, nil
}

func parseLoadStats(line string) ([]*util.Sample, error) {
  var load1, load5, load15 float64
  var running, procs uint64
  var lastpid int

  _, err := fmt.Sscanf(line,
                       "%f %f %f %d/%d %d",
                       &load1, &load5, &load15,
                       &running, &procs, &lastpid)
  if err != nil {
    return nil, err
  }
  samples := []*util.Sample{
    util.NewSample("load.1m", uint64(load1 * 100)),
    util.NewSample("load.5m", uint64(load5 * 100)),
    util.NewSample("load.15m", uint64(load15 * 100)),
    util.NewSample("load.proc", procs),
  }
  return samples, nil
}

