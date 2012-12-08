package main

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "../util"
)

func DiskIOStats() ([]*util.Sample, error) {
  f, err := os.Open("/proc/diskstats")
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
    s, err := parseDiskstatLine(line)
    if err != nil {
      return nil, err
    }
    for _, v := range s {
      samples = append(samples, v)
    }
  }
  return samples, nil
}

func parseDiskstatLine(line string) ([]*util.Sample, error) {
  var major, minor uint
  var dev string
  var rd_ios, rd_merges, rd_sec, rd_ticks uint64
  var wr_ios, wr_merges, wr_sec, wr_ticks uint64
  var ios_in_progress, total_ticks, rq_ticks uint64

  _, err := fmt.Sscanf(line,
                       "%d %d %s %d %d %d %d %d %d %d %d %d %d %d",
                       &major, &minor, &dev,
                       &rd_ios, &rd_merges, &rd_sec, &rd_ticks,
                       &wr_ios, &wr_merges, &wr_sec, &wr_ticks,
                       &ios_in_progress, &total_ticks, &rq_ticks)
  if err != nil {
    return nil, err
  }
  samples := []*util.Sample{
    util.NewSample(fmt.Sprintf("disk.rop.%s", dev), rd_ios),
    util.NewSample(fmt.Sprintf("disk.wop.%s", dev), wr_ios),
    util.NewSample(fmt.Sprintf("disk.rkb.%s", dev), rd_sec / 2),
    util.NewSample(fmt.Sprintf("disk.wkb.%s", dev), wr_sec / 2),
  }
  return samples, nil
}

