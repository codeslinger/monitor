package main

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "strings"
  "../util"
)

func CPUStats() ([]*util.Sample, error) {
  f, err := os.Open("/proc/stat")
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
    s, err := parseStatLine(line)
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

func parseStatLine(line string) ([]*util.Sample, error) {
  if !strings.HasPrefix(line, "cpu") {
    return nil, nil
  }

  var dev string
  var user, nice, sys, idle, iowait, hardirq uint64
  var softirq, steal, guest, guest_nice uint64

  _, err := fmt.Sscanf(line,
                       "%s %d %d %d %d %d %d %d %d %d %d",
                       &dev,
                       &user, &nice, &sys, &idle, &iowait, &hardirq,
                       &softirq, &steal, &guest, &guest_nice)
  if err != nil {
    return nil, err
  }
  if dev == "cpu" {
    // Overall CPU usage line; calculate machine uptime with this.
    // Don't include guest or guest_nice as user and nice already
    // account for these.
    s := util.NewSample("uptime",
                        (user + nice + sys + idle + iowait + hardirq +
                         softirq + steal) / util.HZ)
    return []*util.Sample{s}, nil
  }
  // Individual CPU usage line; calculate per-CPU metrics with this.
  idx := strings.Replace(dev, "cpu", "", 1)
  samples := []*util.Sample{
    util.NewSample(fmt.Sprintf("cpu.%s.user", idx),   user / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.nice", idx),   nice / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.sys", idx),    sys / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.iowait", idx), iowait / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.steal", idx),  steal / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.idle", idx),   idle / util.HZ),
  }
  return samples, nil
}

