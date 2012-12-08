package util

// #include <unistd.h>
// #include <errno.h>
import "C"

var HZ uint64 = getHZ()

func getHZ() uint64 {
  var ticks C.long

  ticks, err := C.sysconf(C._SC_CLK_TCK)
  if err != nil {
    panic(err)
  }
  return uint64(ticks)
}

