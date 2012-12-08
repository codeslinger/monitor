package util

import (
  "bufio"
  "io"
  "os"
)

type Sample struct {
  Name  string
  Value uint64
}

type LineParser func(line string) ([]*Sample, error)

func NewSample(name string, value uint64) *Sample {
  return &Sample{Name: name, Value: value}
}

func StatsFromFile(path string, prsr LineParser) ([]*Sample, error) {
  f, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer f.Close()
  rd := bufio.NewReader(f)
  samples := make([]*Sample, 0)
  for {
    line, err := rd.ReadString('\n')
    if err != nil {
      if err == io.EOF { break }
      return nil, err
    }
    s, err := prsr(line)
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

