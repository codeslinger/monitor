package main

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "strings"
  "syscall"
  "../util"
)

func DiskUsageStats() ([]*util.Sample, error) {
  f, err := os.Open("/etc/mtab")
  if err != nil {
    return nil, err
  }
  defer f.Close()
  rd := bufio.NewReader(f)
  samples := make([]*util.Sample, 0)
  for {
    line, err := rd.ReadString('\n')
    if err != nil {
      if err == io.EOF { break }
      return nil, err
    }
    s, err := parseMtabLine(line)
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

