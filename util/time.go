package util

import (
  "time"
)

const ns_per_ms int64 = 1000000

func NowMS() int64 {
  return time.Now().UnixNano() / ns_per_ms
}

func MSToTime(ms int64) time.Time {
  return time.Unix(ms / 1000, ms % 1000)
}

