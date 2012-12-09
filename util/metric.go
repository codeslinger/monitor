package util

import (
  "bufio"
  "fmt"
  "io"
  "os"
)

type Metadata struct {
  Name  string
  Value string
}

func NewMetadata(name, value string) *Metadata {
  return &Metadata{Name: name, Value: value}
}

func (m *Metadata) String() string {
  return fmt.Sprintf("M|%s|%s", m.Name, m.Value)
}

type Sample struct {
  Name  string
  Value uint64
}

func NewSample(name string, value uint64) *Sample {
  return &Sample{Name: name, Value: value}
}

func (s *Sample) String() string {
  return fmt.Sprintf("S|%s|%d", s.Name, s.Value)
}

type LineParser func(line string) ([]*Sample, error)

func StatsFromFile(path string, parser LineParser) ([]*Sample, error) {
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
    s, err := parser(line)
    if err != nil {
      return nil, err
    }
    if s != nil {
      for _, v := range s {
        if v != nil {
          samples = append(samples, v)
        }
      }
    }
  }
  return samples, nil
}

