package linux

import (
  "fmt"
  "github.com/bmizerany/assert"
  "io"
  "strings"
  "testing"
)

var (
  testHZ uint64 = 100
  procStatsOutput =
`cpu  1377723 12309 425558 92572282 176914 102 11966 0 0 0
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
softirq 37177398 6 8570721 477 1028639 715899 6 16037399 4556202 71761 6196288
`
  procLoadavgOutput = "0.00 0.02 0.05 1/406 16439"
  procMeminfoOutput =
`MemTotal:        3353936 kB
MemFree:         1071244 kB
Buffers:          149540 kB
Cached:           770616 kB
SwapCached:            0 kB
Active:          1462484 kB
Inactive:         609376 kB
Active(anon):    1153200 kB
Inactive(anon):    70800 kB
Active(file):     309284 kB
Inactive(file):   538576 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:       3487740 kB
SwapFree:        3487740 kB
Dirty:                44 kB
Writeback:             0 kB
AnonPages:       1151620 kB
Mapped:           171312 kB
Shmem:             72300 kB
Slab:             103716 kB
SReclaimable:      76384 kB
SUnreclaim:        27332 kB
KernelStack:        3176 kB
PageTables:        31652 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:     5164708 kB
Committed_AS:    3109576 kB
VmallocTotal:   34359738367 kB
VmallocUsed:      359696 kB
VmallocChunk:   34359368428 kB
HardwareCorrupted:     0 kB
AnonHugePages:         0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:       60396 kB
DirectMap2M:     3428352 kB
`
  procDiskstatsOutput =
`   1       0 ram0 0 0 0 0 0 0 0 0 0 0 0
   1       1 ram1 0 0 0 0 0 0 0 0 0 0 0
   1       2 ram2 0 0 0 0 0 0 0 0 0 0 0
   1       3 ram3 0 0 0 0 0 0 0 0 0 0 0
   1       4 ram4 0 0 0 0 0 0 0 0 0 0 0
   1       5 ram5 0 0 0 0 0 0 0 0 0 0 0
   1       6 ram6 0 0 0 0 0 0 0 0 0 0 0
   1       7 ram7 0 0 0 0 0 0 0 0 0 0 0
   1       8 ram8 0 0 0 0 0 0 0 0 0 0 0
   1       9 ram9 0 0 0 0 0 0 0 0 0 0 0
   1      10 ram10 0 0 0 0 0 0 0 0 0 0 0
   1      11 ram11 0 0 0 0 0 0 0 0 0 0 0
   1      12 ram12 0 0 0 0 0 0 0 0 0 0 0
   1      13 ram13 0 0 0 0 0 0 0 0 0 0 0
   1      14 ram14 0 0 0 0 0 0 0 0 0 0 0
   1      15 ram15 0 0 0 0 0 0 0 0 0 0 0
   7       0 loop0 0 0 0 0 0 0 0 0 0 0 0
   7       1 loop1 0 0 0 0 0 0 0 0 0 0 0
   7       2 loop2 0 0 0 0 0 0 0 0 0 0 0
   7       3 loop3 0 0 0 0 0 0 0 0 0 0 0
   7       4 loop4 0 0 0 0 0 0 0 0 0 0 0
   7       5 loop5 0 0 0 0 0 0 0 0 0 0 0
   7       6 loop6 0 0 0 0 0 0 0 0 0 0 0
   7       7 loop7 0 0 0 0 0 0 0 0 0 0 0
   8       0 sda 50762 6347 1674054 67360 23942 20742 563152 25820 0 14200 93132
   8       1 sda1 50432 6316 1671178 67232 23118 20742 563152 24696 0 13088 91876
   8       2 sda2 2 0 4 0 0 0 0 0 0 0 0
   8       5 sda5 161 31 1536 60 0 0 0 0 0 60 60
`
  etcMtabOutput =
`/dev/sda1 / ext4 rw,errors=remount-ro 0 0
proc /proc proc rw,noexec,nosuid,nodev 0 0
sysfs /sys sysfs rw,noexec,nosuid,nodev 0 0
none /sys/fs/fuse/connections fusectl rw 0 0
none /sys/kernel/debug debugfs rw 0 0
none /sys/kernel/security securityfs rw 0 0
udev /dev devtmpfs rw,mode=0755 0 0
devpts /dev/pts devpts rw,noexec,nosuid,gid=5,mode=0620 0 0
tmpfs /run tmpfs rw,noexec,nosuid,size=10%,mode=0755 0 0
none /run/lock tmpfs rw,noexec,nosuid,nodev,size=5242880 0 0
none /run/shm tmpfs rw,nosuid,nodev 0 0
gvfs-fuse-daemon /home/blorp/.gvfs fuse.gvfs-fuse-daemon rw,nosuid,nodev,user=blorp 0 0
`
  procNetDevOutput =
`Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo: 5303451   31684    0    0    0     0          0         0  5303451   31684    0    0    0     0       0          0
  eth0: 16651766   30158    0    0    0     0          0         0  4036294   22014    0    0    0     0       0          0
`
)

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

func Test_CPUSampler_should_parse_valid_proc_stats_file_properly(t *testing.T) {
  wr := NewBufferedSampleWriter()
  sampler := NewCPUSampler(NewStringOpener(procStatsOutput), wr)
  sampler.HZ = testHZ
  if err := sampler.Sample(); err != nil {
    t.Fatalf("Sample() failed: %s", err)
  }
  assert.Equal(t, "cpu 0 4450 1844 59 679 0 228028\n", wr.Lines[0])
  assert.Equal(t, "cpu 1 2538 573 7 109 0 233728\n", wr.Lines[1])
  assert.Equal(t, "cpu 2 4463 1275 46 809 0 230051\n", wr.Lines[2])
  assert.Equal(t, "cpu 3 2325 561 9 169 0 233913\n", wr.Lines[3])
}

func Test_LoadSampler_should_parse_valid_proc_loadavg_file_properly(t *testing.T) {
  wr := NewBufferedSampleWriter()
  sampler := NewLoadSampler(NewStringOpener(procLoadavgOutput), wr)
  if err := sampler.Sample(); err != nil {
    t.Fatalf("Sample() failed: %s", err)
  }
  assert.Equal(t, "load 0 0.02 0.05 406\n", wr.Lines[0])
}

func Test_MemorySampler_should_parse_valid_proc_meminfo_file_properly(t *testing.T) {
  wr := NewBufferedSampleWriter()
  sampler := NewMemorySampler(NewStringOpener(procMeminfoOutput), wr)
  if err := sampler.Sample(); err != nil {
    t.Fatalf("Sample() failed: %s", err)
  }
  assert.Equal(
    t,
    "memory 3353936 1071244 149540 770616 3487740 3487740\n",
    wr.Lines[0])
}

func Test_DiskIOSampler_should_parse_valid_proc_diskstats_file_properly(t *testing.T) {
  wr := NewBufferedSampleWriter()
  sampler := NewDiskIOSampler(NewStringOpener(procDiskstatsOutput), wr)
  if err := sampler.Sample(); err != nil {
    t.Fatalf("Sample() failed: %s", err)
  }
  assert.Equal(t, "disk sda 50762 837027 23942 281576\n", wr.Lines[0])
  assert.Equal(t, "disk sda1 50432 835589 23118 281576\n", wr.Lines[1])
  assert.Equal(t, "disk sda2 2 2 0 0\n", wr.Lines[2])
  assert.Equal(t, "disk sda5 161 768 0 0\n", wr.Lines[3])
}

func Test_NICSampler_should_parse_valid_proc_net_dev_file_properly(t *testing.T) {
  wr := NewBufferedSampleWriter()
  sampler := NewNICSampler(NewStringOpener(procNetDevOutput), wr)
  if err := sampler.Sample(); err != nil {
    t.Fatalf("Sample() failed: %s", err)
  }
  assert.Equal(
    t,
    "net eth0 16651766 30158 0 0 4036294 22014 0 0\n",
    wr.Lines[0])
}

