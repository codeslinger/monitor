package util

import (
  "log"
  "log/syslog"
)

var LOG *log.Logger

func StartLogger() (err error) {
  LOG, err = syslog.NewLogger(syslog.LOG_NOTICE, log.Lshortfile)
  return
}

func Log(fmt string, v ...interface{}) {
  LOG.Printf(fmt, v...)
}

