package util

import (
  "fmt"
  "io"
  "os"
)

// Opener that returns opened files.
type FileOpener struct {}

func NewFileOpener() *FileOpener {
  return &FileOpener{}
}

// Opens the given file path.
func (f *FileOpener) Open(path string) (io.ReadCloser, error) {
  return os.Open(path)
}

// Writes samples out to stdout.
type ConsoleSampleWriter struct {}

func NewConsoleSampleWriter() *ConsoleSampleWriter {
  return &ConsoleSampleWriter{}
}

// Write the given sample out to stdout.
func (c *ConsoleSampleWriter) Write(v ...interface{}) {
  fmt.Println(v...)
}

