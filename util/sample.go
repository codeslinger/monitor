package util

import (
  "errors"
  "fmt"
  "time"
)

type SampleType int

type Sample struct {
  Timestamp time.Time
  Name      string
  Value     int64
  Type      SampleType
}

const (
  GAUGE SampleType = iota
  COUNTER
)

var typeToLabel = map[SampleType]string {
  GAUGE:   "g",
  COUNTER: "c",
}

var labelToType = map[string]SampleType {
  "g": GAUGE,
  "c": COUNTER,
}

var MalformedSampleSyntax = errors.New("malformed sample syntax")
var UnknownSampleType = errors.New("unknown sample type specified")

const sampFmt string = "%d|%s:%d|%s\n"
const ns_per_ms int64 = 1000 * 1000

// Create a new Sample record for a gauge sample.
func NewGauge(name string, value int64, t SampleType) *Sample {
  return &Sample{Name: name, Value: value, Type: GAUGE}
}

// Create a new Sample record for a counter sample.
func NewCounter(name string, value int64) *Sample {
  return &Sample{Name: name, Value: value, Type: COUNTER}
}

// Serialize this Sample to a string.
func (s *Sample) String() string {
  return fmt.Sprintf(sampFmt, NowMS() s.Name, s.Value, typeToLabel[s.Type])
}

// Deserialize a Sample record from a string.
func ParseSample(buf string) (*Sample, error) {
  s := &Sample{}
  t := ""
  ts := -1
  n, err := fmt.Scanf(sampFmt, ts, s.Name, s.Value, t)
  if err != nil {
    return nil, err
  }
  if n < 3 {
    return nil, MalformedSampleSyntax
  }
  if ts < 0 {
    return nil, InvalidTimestamp
  }
  it, ok := labelToType[t]
  if !ok {
    return nil, UnknownSampleType
  }
  s.Timestamp = MSToTime(ts)
  s.Type = it
  return s, nil
}

