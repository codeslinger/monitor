package linux

import (
  "fmt"
  "log"
  "strconv"
  "strings"
  "syscall"
  "../../util"
)

var linuxStatsFiles = map[string]util.LineParser{
  "/proc/stat":      parseStatLine,
  "/proc/loadavg":   parseLoadLine,
  "/proc/meminfo":   parseMeminfoLine,
  "/proc/diskstats": parseDiskstatLine,
  "/etc/mtab":       parseMtabLine,
  "/proc/net/dev":   parseNetdevLine,
}

func StandardStats() []*util.Sample {
  samples := make([]*util.Sample, 0)
  for path, parser := range linuxStatsFiles {
    s, err := util.StatsFromFile(path, parser)
    if err != nil {
      log.Printf("parser failure: %s: %s\n", path, err)
      continue
    }
    for _, v := range s {
      samples = append(samples, v)
    }
  }
  return samples
}

func parseStatLine(line string) ([]*util.Sample, error) {
  if !strings.HasPrefix(line, "cpu") {
    return nil, nil
  }

  var dev string
  var user, nice, sys, idle, iowait, hardirq uint64
  var softirq, steal, guest, guest_nice uint64

  _, err := fmt.Sscanf(line,
                       "%s %d %d %d %d %d %d %d %d %d %d",
                       &dev,
                       &user, &nice, &sys, &idle, &iowait, &hardirq,
                       &softirq, &steal, &guest, &guest_nice)
  if err != nil {
    return nil, err
  }
  if dev == "cpu" {
    // Overall CPU usage line; calculate machine uptime with this.
    // Don't include guest or guest_nice as user and nice already
    // account for these.
    s := util.NewSample("uptime",
                        (user + nice + sys + idle + iowait + hardirq +
                         softirq + steal) / util.HZ)
    return []*util.Sample{s}, nil
  }
  // Individual CPU usage line; calculate per-CPU metrics with this.
  idx := strings.Replace(dev, "cpu", "", 1)
  samples := []*util.Sample{
    util.NewSample(fmt.Sprintf("cpu.%s.user", idx),   user / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.nice", idx),   nice / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.sys", idx),    sys / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.iowait", idx), iowait / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.steal", idx),  steal / util.HZ),
    util.NewSample(fmt.Sprintf("cpu.%s.idle", idx),   idle / util.HZ),
  }
  return samples, nil
}

func parseLoadLine(line string) ([]*util.Sample, error) {
  var load1, load5, load15 float64
  var running, procs uint64
  var lastpid int

  _, err := fmt.Sscanf(line,
                       "%f %f %f %d/%d %d",
                       &load1, &load5, &load15,
                       &running, &procs, &lastpid)
  if err != nil {
    return nil, err
  }
  samples := []*util.Sample{
    util.NewSample("load.1m", uint64(load1 * 100)),
    util.NewSample("load.5m", uint64(load5 * 100)),
    util.NewSample("load.15m", uint64(load15 * 100)),
    util.NewSample("load.proc", procs),
  }
  return samples, nil
}

func parseMeminfoLine(line string) ([]*util.Sample, error) {
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
  return []*util.Sample{s}, nil
}

func parseDiskstatLine(line string) ([]*util.Sample, error) {
  var major, minor uint
  var dev string
  var rd_ios, rd_merges, rd_sec, rd_ticks uint64
  var wr_ios, wr_merges, wr_sec, wr_ticks uint64
  var ios_in_progress, total_ticks, rq_ticks uint64

  _, err := fmt.Sscanf(line,
                       "%d %d %s %d %d %d %d %d %d %d %d %d %d %d",
                       &major, &minor, &dev,
                       &rd_ios, &rd_merges, &rd_sec, &rd_ticks,
                       &wr_ios, &wr_merges, &wr_sec, &wr_ticks,
                       &ios_in_progress, &total_ticks, &rq_ticks)
  if err != nil {
    return nil, err
  }
  samples := []*util.Sample{
    util.NewSample(fmt.Sprintf("disk.rop.%s", dev), rd_ios),
    util.NewSample(fmt.Sprintf("disk.wop.%s", dev), wr_ios),
    util.NewSample(fmt.Sprintf("disk.rkb.%s", dev), rd_sec / 2),
    util.NewSample(fmt.Sprintf("disk.wkb.%s", dev), wr_sec / 2),
  }
  return samples, nil
}

func parseMtabLine(line string) ([]*util.Sample, error) {
  var dev, mount, fstype, options string
  var dump, fsck_order uint64
  var buf syscall.Statfs_t

  _, err := fmt.Sscanf(line,
                       "%s %s %s %s %d %d",
                       &dev, &mount, &fstype, &options,
                       &dump, &fsck_order)
  if err != nil {
    return nil, err
  }
  // only report on ext[234] filesystems for now
  if !strings.HasPrefix(fstype, "ext") {
    return nil, nil
  }
  if err := syscall.Statfs(mount, &buf); err != nil {
    return nil, err
  }
  if buf.Blocks == 0 {
    return nil, nil
  }
  to1K := uint64(buf.Bsize) / 1024
  samples := []*util.Sample{
    util.NewSample(fmt.Sprintf("fs.%s.total", mount), buf.Blocks * to1K),
    util.NewSample(fmt.Sprintf("fs.%s.free", mount), buf.Bfree * to1K),
    util.NewSample(fmt.Sprintf("fs.%s.avail", mount), buf.Bavail * to1K),
    util.NewSample(fmt.Sprintf("fs.%s.inode", mount), buf.Files),
    util.NewSample(fmt.Sprintf("fs.%s.ifree", mount), buf.Ffree),
  }
  return samples, nil
}

func parseNetdevLine(line string) ([]*util.Sample, error) {
  if strings.Contains(line, "|") {
    return nil, nil
  }

  var dev string
  var rx_byte, rx_pkt, rx_err, rx_drop uint64
  var rx_fifo, rx_frame, rx_comp, rx_mcast uint64
  var tx_byte, tx_pkt, tx_err, tx_drop uint64
  var tx_fifo, tx_coll, tx_carr, tx_comp uint64

  _, err := fmt.Sscanf(
    line,
    "%s: %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
    &dev,
    &rx_byte, &rx_pkt, &rx_err, &rx_drop,
    &rx_fifo, &rx_frame, &rx_comp, &rx_mcast,
    &tx_byte, &tx_pkt, &tx_err, &tx_drop,
    &tx_fifo, &tx_coll, &tx_carr, &tx_comp)
  if err != nil {
    return nil, err
  }
  samples := []*util.Sample{
    util.NewSample("net.%s.rx.bytes", rx_byte),
    util.NewSample("net.%s.rx.pkts", rx_pkt),
    util.NewSample("net.%s.rx.errs", rx_pkt),
    util.NewSample("net.%s.rx.drop", rx_pkt),
    util.NewSample("net.%s.tx.bytes", tx_byte),
    util.NewSample("net.%s.tx.pkts", tx_pkt),
    util.NewSample("net.%s.tx.errs", tx_err),
    util.NewSample("net.%s.tx.drop", tx_drop),
  }
  return samples, nil
}

