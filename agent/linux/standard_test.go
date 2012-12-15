package linux

import (
  "fmt"
  "io"
  "strings"
  "testing"
)

var testHZ uint64 = 100
var procStatsOutput = `cpu  1377723 12309 425558 92572282 176914 102 11966 0 0 0
cpu0 445076 5965 184472 22802862 67989 101 11569 0 0 0
cpu1 253806 742 57397 23372889 10999 0 55 0 0 0
cpu2 446324 4621 127575 23005164 80953 0 286 0 0 0
cpu3 232515 981 56113 23391366 16971 0 55 0 0 0
intr 95815165 832 265035 0 0 0 0 0 0 1 339121 0 0 2331321 0 0 0 416 0 0 0 0 0 0 2739 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 799772 15 0 0 0 0 22176 13598064 6555331 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 162083424
btime 1355344417
processes 44685
procs_running 1
procs_blocked 0
softirq 37177398 6 8570721 477 1028639 715899 6 16037399 4556202 71761 6196288`

type StringReadCloser struct {
  rd io.Reader
}

func NewStringReadCloser(rd io.Reader) *StringReadCloser {
  return &StringReadCloser{rd: rd}
}

func (s *StringReadCloser) Read(p []byte) (int, error) {
  return s.rd.Read(p)
}

func (s *StringReadCloser) Close() error {
  return nil
}

type StringOpener struct {
  data string
}

func NewStringOpener(data string) *StringOpener {
  return &StringOpener{data: data}
}

func (s *StringOpener) Open(path string) (io.ReadCloser, error) {
  return NewStringReadCloser(strings.NewReader(s.data)), nil
}

type BufferedSampleWriter struct {
  Lines []string
}

func NewBufferedSampleWriter() *BufferedSampleWriter {
  return &BufferedSampleWriter{Lines: make([]string, 0)}
}

func (b *BufferedSampleWriter) Init() error {
  return nil
}

func (b *BufferedSampleWriter) Write(v ...interface{}) {
  b.Lines = append(b.Lines, fmt.Sprintln(v...))
}

func AssertStringEquals(t *testing.T, expected, actual string) {
  if expected != actual {
    t.Fatalf("failed expectation: '%s' != '%s'",
             strings.Trim(expected, "\r\n"),
             strings.Trim(actual, "\r\n"))
  }
}

func AssertUint64Equals(t *testing.T, expected, actual uint64) {
  if expected != actual {
    t.Fatalf("failed expectation: %d != %d", expected, actual)
  }
}

func Test_CPUSampler_should_parse_valid_proc_stats_file_properly(t *testing.T) {
  wr := NewBufferedSampleWriter()
  sampler := NewCPUSampler(NewStringOpener(procStatsOutput), wr)
  sampler.HZ = testHZ
  if err := sampler.Sample(); err != nil {
    t.Fatalf("sampler failed: %s", err)
  }
  AssertStringEquals(t, "uptime 945768\n", wr.Lines[0])
  AssertStringEquals(t, "cpu 0 4450 1844 59 679 0 228028\n", wr.Lines[1])
  AssertStringEquals(t, "cpu 1", wr.Lines[2])
  AssertStringEquals(t, "cpu 2", wr.Lines[3])
  AssertStringEquals(t, "cpu 3", wr.Lines[4])
}

