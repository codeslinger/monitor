package main

import (
  "log"
  "log/syslog"
)

var LOG *log.Logger

func StartLogger() (err error) {
  LOG, err = syslog.NewLogger(syslog.LOG_NOTICE, log.Lshortfile)
  return
}

