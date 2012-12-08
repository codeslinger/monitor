package main

import (
  "bufio"
  "errors"
  "io"
  "os"
  "strconv"
  "strings"
  "../util"
)

var (
  missingFieldsErr = errors.New("some memory info fields missing")
)

func MemoryStats() ([]*util.Sample, error) {
  f, err := os.Open("/proc/meminfo")
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
    s, err := parseMeminfoLine(line)
    if err != nil {
      return nil, err
    }
    if s != nil {
      samples = append(samples, s)
    }
    if len(samples) == 6 {
      break
    }
  }
  if len(samples) < 6 {
    return nil, missingFieldsErr
  }
  return samples, nil
}

func parseMeminfoLine(line string) (*util.Sample, error) {
  var s *util.Sample

  p := strings.SplitN(strings.Trim(line, " \r\n"), ":", 2)
  if len(p) < 2 {
    return nil, nil
  }
  cmp := strings.Split(strings.Trim(p[1], " \r\n"), " ")[0]
  val, err := strconv.ParseUint(cmp, 10, 64)
  if err != nil {
    return nil, err
  }
  switch p[0] {
  case "MemTotal":  s = util.NewSample("mem.total", val)
  case "MemFree":   s = util.NewSample("mem.free", val)
  case "Buffers":   s = util.NewSample("mem.buffer", val)
  case "Cached":    s = util.NewSample("mem.cache", val)
  case "SwapFree":  s = util.NewSample("mem.swfree", val)
  case "SwapTotal": s = util.NewSample("mem.swtot", val)
  default:          s = nil
  }
  return s, nil
}

