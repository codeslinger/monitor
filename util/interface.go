package util

import (
  "io"
)

// Interface for objects that can open resources (e.g. files, sockets, etc).
type Opener interface {
  Open(path string) (io.ReadCloser, error)
}

// Interface for objects that sample metrics from their host system or 
// applications periodically.
type Sampler interface {
  Init()   error
  Sample() error
}

// Interface for objects that can write samples to a sink (e.g. file,
// socket, etc).
type SampleWriter interface {
  Write(v ...interface{})
}

