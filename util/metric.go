package util

type Sample struct {
  Name  string
  Value uint64
}

func NewSample(name string, value uint64) *Sample {
  return &Sample{Name: name, Value: value}
}

