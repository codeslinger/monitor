package linux

import (
  "fmt"
  "os"
  "net"
  "strings"
  "../../util"
)

func Metadata() ([]*util.Metadata, error) {
  host, err := os.Hostname()
  if err != nil {
    return nil, err
  }
  ips, err := ipv4ForInterfaces()
  if err != nil {
    return nil, err
  }
  m := make([]*util.Metadata, 0)
  m = append(m, util.NewMetadata("host", host))
  for ifname, ip := range ips {
    m = append(m, util.NewMetadata(fmt.Sprintf("nic.%s", ifname), ip.String()))
  }
  return m, nil
}

func ipv4ForInterfaces() (map[string]net.IP, error) {
  rv := make(map[string]net.IP)
  ifs, err := net.Interfaces()
  if err != nil {
    return nil, err
  }
  for _, i := range ifs {
    if strings.HasPrefix(i.Name, "lo") {
      continue
    }
    addrs, err := i.Addrs()
    if err != nil {
      return nil, err
    }
    for _, addr := range addrs {
      if ip := addr.(*net.IPNet).IP.To4(); ip != nil {
        rv[i.Name] = ip
        break
      }
    }
  }
  return rv, nil
}

