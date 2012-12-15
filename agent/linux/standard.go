// Linux-specific OS and machine-level statistics gathering.
package linux

import (
  "bufio"
  "fmt"
  "io"
  "strconv"
  "strings"
  "syscall"
  "../../util"
)

// #include <unistd.h>
// #include <errno.h>
import "C"

// Retrieve the clock ticks per second on this kernel.
func getHZ() (uint64, error) {
  ticks, err := C.sysconf(C._SC_CLK_TCK)
  if err != nil {
    return 0, err
  }
  return uint64(ticks), nil
}

// Sampler for standard OS- and machine-level metrics. This sampler is an
// aggregate of the CPUSampler, LoadSampler, MemorySampler, DiskIOSampler,
// FSUsageSampler and NICSampler samplers.
type StandardSampler struct {
  cpu  *CPUSampler
  load *LoadSampler
  mem  *MemorySampler
  disk *DiskIOSampler
  fs   *FSUsageSampler
  nic  *NICSampler
}

// Create a new standard sampler.
func NewStandardSampler(o util.Opener, s util.SampleWriter) *StandardSampler {
  return &StandardSampler{
    cpu: NewCPUSampler(o, s),
    load: NewLoadSampler(o, s),
    mem: NewMemorySampler(o, s),
    disk: NewDiskIOSampler(o, s),
    fs: NewFSUsageSampler(o, s),
    nic: NewNICSampler(o, s),
  }
}

// Initialize this sampler. Delegates initialization to its underlying
// samplers.
func (standard *StandardSampler) Init() (err error) {
  if err = standard.cpu.Init(); err != nil { return }
  if err = standard.load.Init(); err != nil { return }
  if err = standard.mem.Init(); err != nil { return }
  if err = standard.disk.Init(); err != nil { return }
  if err = standard.fs.Init(); err != nil { return }
  if err = standard.nic.Init(); err != nil { return }
  return
}

// Gather samples for all underlying samplers.
func (standard *StandardSampler) Sample() (err error) {
  if err = standard.cpu.Sample(); err != nil { return }
  if err = standard.load.Sample(); err != nil { return }
  if err = standard.mem.Sample(); err != nil { return }
  if err = standard.disk.Sample(); err != nil { return }
  if err = standard.fs.Sample(); err != nil { return }
  if err = standard.nic.Sample(); err != nil { return }
  return
}

// Sampler for CPU utilization metrics.
type CPUSampler struct {
  opener util.Opener
  sink   util.SampleWriter
  HZ     uint64
}

// Create a new CPU utilization sampler.
func NewCPUSampler(o util.Opener, s util.SampleWriter) *CPUSampler {
  return &CPUSampler{opener: o, sink: s}
}

// Initialize this sampler.
func (stats *CPUSampler) Init() (err error) {
  stats.HZ, err = getHZ()
  return
}

// Gather current CPU utilization sample.
func (stats *CPUSampler) Sample() (err error) {
  f, err := stats.opener.Open("/proc/stat")
  if err != nil {
    return
  }
  defer f.Close()
  rd := bufio.NewReader(f)
  for {
    var line string

    line, err = rd.ReadString('\n')
    if err == io.EOF {
      err = nil
      break
    } else if err != nil {
      return
    }
    if err = stats.parseLine(line); err != nil {
      return
    }
  }
  return
}

// Parse an individual line from Linux's /proc/stat.
func (stats *CPUSampler) parseLine(line string) (err error) {
  if !strings.HasPrefix(line, "cpu") {
    return
  }

  var dev string
  var user, nice, sys, idle, iowait, hardirq uint64
  var softirq, steal, guest, guest_nice uint64

  _, err = fmt.Sscanf(line,
                      "%s %d %d %d %d %d %d %d %d %d %d",
                      &dev,
                      &user, &nice, &sys, &idle, &iowait, &hardirq,
                      &softirq, &steal, &guest, &guest_nice)
  if err != nil {
    return
  }
  if dev == "cpu" {
    // Overall CPU usage line; calculate machine uptime with this.
    // Don't include guest or guest_nice as user and nice already
    // account for these.
    stats.sink.Write("uptime",
                     (user + nice + sys + idle + iowait + hardirq + softirq + steal) / stats.HZ)
    return
  }
  // Individual CPU usage line; calculate per-CPU metrics with this.
  idx := strings.Replace(dev, "cpu", "", 1)
  stats.sink.Write("cpu",
                   idx,
                   user / stats.HZ,
                   sys / stats.HZ,
                   nice / stats.HZ,
                   iowait / stats.HZ,
                   steal / stats.HZ,
                   idle / stats.HZ)
  return
}

// Sampler for system load statistics.
type LoadSampler struct {
  opener util.Opener
  sink   util.SampleWriter
}

// Create a new system load sampler.
func NewLoadSampler(o util.Opener, s util.SampleWriter) *LoadSampler {
  return &LoadSampler{opener: o, sink: s}
}

// Initialize this sampler.
func (load *LoadSampler) Init() (err error) {
  return
}

// Gather current system load statistics.
func (load *LoadSampler) Sample() (err error) {
  f, err := load.opener.Open("/proc/loadavg")
  if err != nil {
    return
  }
  defer f.Close()
  rd := bufio.NewReader(f)
  var line string
  line, err = rd.ReadString('\n')
  if err != nil {
    return
  }
  err = load.parseLine(line)
  return
}

// Parse Linux's /proc/loadavg.
func (load *LoadSampler) parseLine(line string) (err error) {
  var load1, load5, load15 float64
  var running, procs uint64
  var lastpid int

  _, err = fmt.Sscanf(line,
                      "%f %f %f %d/%d %d",
                      &load1, &load5, &load15,
                      &running, &procs, &lastpid)
  if err != nil {
    return
  }
  load.sink.Write("load", load1, load5, load15, procs)
  return
}

// Sampler for RAM usage statistics.
type MemorySampler struct {
  opener util.Opener
  sink   util.SampleWriter
}

// Create a new RAM usage sampler.
func NewMemorySampler(o util.Opener, s util.SampleWriter) *MemorySampler {
  return &MemorySampler{opener: o, sink: s}
}

// Initialize this sampler.
func (mem *MemorySampler) Init() (err error) {
  return
}

// Gather current RAM usage statistics.
func (mem *MemorySampler) Sample() (err error) {
  var total, free, buffer, cache, swapfree, swaptotal uint64

  f, err := mem.opener.Open("/proc/meminfo")
  if err != nil {
    return
  }
  defer f.Close()
  rd := bufio.NewReader(f)
  for {
    var line string

    line, err = rd.ReadString('\n')
    if err == io.EOF {
      break
    } else if err != nil {
      return
    }
    parts := strings.SplitN(line, ":", 2)
    switch parts[0] {
    case "MemTotal": total, err = mem.toUint(parts[1])
    case "MemFree": free, err = mem.toUint(parts[1])
    case "Buffers": buffer, err = mem.toUint(parts[1])
    case "Cached": cache, err = mem.toUint(parts[1])
    case "SwapFree": swapfree, err = mem.toUint(parts[1])
    case "SwapTotal": swaptotal, err = mem.toUint(parts[1])
    }
    if err != nil {
      return
    }
  }
  mem.sink.Write("memory", total, free, buffer, cache, swaptotal, swapfree)
  return
}

// Parse a uint out of a /proc/meminfo line.
func (mem *MemorySampler) toUint(raw string) (uint64, error) {
  return strconv.ParseUint(strings.Split(raw, " ")[0], 10, 64)
}

// Sampler for disk I/O statistics.
type DiskIOSampler struct {
  opener util.Opener
  sink   util.SampleWriter
}

// Create a new disk I/O sampler.
func NewDiskIOSampler(o util.Opener, s util.SampleWriter) *DiskIOSampler {
  return &DiskIOSampler{opener: o, sink: s}
}

// Initialize this sampler.
func (disk *DiskIOSampler) Init() (err error) {
  return
}

// Gather current disk I/O statistics.
func (disk *DiskIOSampler) Sample() (err error) {
  f, err := disk.opener.Open("/proc/diskstats")
  if err != nil {
    return
  }
  defer f.Close()
  rd := bufio.NewReader(f)
  for {
    var line string

    line, err = rd.ReadString('\n')
    if err == io.EOF {
      break
    } else if err != nil {
      return
    }
    if err = disk.parseLine(line); err != nil {
      return
    }
  }
  return
}

// Parse an individual line from Linux's /proc/diskstats.
func (disk *DiskIOSampler) parseLine(line string) (err error) {
  var major, minor uint
  var dev string
  var rd_ios, rd_merges, rd_sec, rd_ticks uint64
  var wr_ios, wr_merges, wr_sec, wr_ticks uint64
  var ios_in_progress, total_ticks, rq_ticks uint64

  _, err = fmt.Sscanf(line,
                      "%d %d %s %d %d %d %d %d %d %d %d %d %d %d",
                      &major, &minor, &dev,
                      &rd_ios, &rd_merges, &rd_sec, &rd_ticks,
                      &wr_ios, &wr_merges, &wr_sec, &wr_ticks,
                      &ios_in_progress, &total_ticks, &rq_ticks)
  if err != nil {
    return
  }
  // skip garbage devices
  if strings.HasPrefix(dev, "ram") || strings.HasPrefix(dev, "loop") {
    return
  }
  disk.sink.Write("disk", dev, rd_ios, rd_sec / 2, wr_ios, wr_sec / 2)
  return
}

// Sampler for filesystem usage statistics.
type FSUsageSampler struct {
  opener util.Opener
  sink   util.SampleWriter
}

// Create a new filesystem usage sampler.
func NewFSUsageSampler(o util.Opener, s util.SampleWriter) *FSUsageSampler {
  return &FSUsageSampler{opener: o, sink: s}
}

// Initialize this sampler.
func (fs *FSUsageSampler) Init() (err error) {
  return
}

// Gather current filesystem usage statistics.
func (fs *FSUsageSampler) Sample() (err error) {
  f, err := fs.opener.Open("/etc/mtab")
  if err != nil {
    return
  }
  defer f.Close()
  rd := bufio.NewReader(f)
  for {
    var line string

    line, err = rd.ReadString('\n')
    if err == io.EOF {
      break
    } else if err != nil {
      return
    }
    if err = fs.parseLine(line); err != nil {
      return
    }
  }
  return
}

// Parse an individual line from Linux's /etc/mtab.
func (fs *FSUsageSampler) parseLine(line string) (err error) {
  var dev, mount, fstype, options string
  var dump, fsck_order uint64
  var buf syscall.Statfs_t

  _, err = fmt.Sscanf(line,
                      "%s %s %s %s %d %d",
                      &dev, &mount, &fstype, &options,
                      &dump, &fsck_order)
  if err != nil {
    return
  }
  // only report on ext[234] filesystems for now
  if !strings.HasPrefix(fstype, "ext") {
    return
  }
  if err = syscall.Statfs(mount, &buf); err != nil {
    return
  }
  if buf.Blocks == 0 {
    return
  }
  to1K := uint64(buf.Bsize) / 1024
  fs.sink.Write("fs", mount, buf.Blocks * to1K, buf.Bfree * to1K,
                buf.Bavail * to1K, buf.Files, buf.Ffree)
  return
}

// Sampler for NIC utilization.
type NICSampler struct {
  opener util.Opener
  sink   util.SampleWriter
}

// Create a new NIC utilization sampler.
func NewNICSampler(o util.Opener, s util.SampleWriter) *NICSampler {
  return &NICSampler{opener: o, sink: s}
}

// Initialize this sampler.
func (nic *NICSampler) Init() (err error) {
  return
}

// Gather current NIC utilization statistics.
func (nic *NICSampler) Sample() (err error) {
  f, err := nic.opener.Open("/etc/mtab")
  if err != nil {
    return
  }
  defer f.Close()
  rd := bufio.NewReader(f)
  for {
    var line string

    line, err = rd.ReadString('\n')
    if err == io.EOF {
      break
    } else if err != nil {
      return
    }
    if err = nic.parseLine(line); err != nil {
      return
    }
  }
  return
}

// Parse an individual line from Linux's /proc/net/dev.
func (nic *NICSampler) parseLine(line string) (err error) {
  var dev string
  var rx_byte, rx_pkt, rx_err, rx_drop uint64
  var rx_fifo, rx_frame, rx_comp, rx_mcast uint64
  var tx_byte, tx_pkt, tx_err, tx_drop uint64
  var tx_fifo, tx_coll, tx_carr, tx_comp uint64

  if strings.Contains(line, "|") {
    return
  }
  _, err = fmt.Sscanln(line,
                       &dev,
                       &rx_byte, &rx_pkt, &rx_err, &rx_drop,
                       &rx_fifo, &rx_frame, &rx_comp, &rx_mcast,
                       &tx_byte, &tx_pkt, &tx_err, &tx_drop,
                       &tx_fifo, &tx_coll, &tx_carr, &tx_comp)
  if err != nil {
    return
  }
  dev = dev[0:len(dev)-1]
  if strings.HasPrefix(dev, "lo") {
    return
  }
  nic.sink.Write("net", dev,
                 rx_byte, rx_pkt, rx_err, rx_drop,
                 tx_byte, tx_pkt, tx_err, tx_drop)
  return
}

