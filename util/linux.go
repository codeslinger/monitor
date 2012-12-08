package util

// #include <unistd.h>
// #include <errno.h>
import "C"

var HZ = getHZ()

func getHZ() uint {
  var ticks C.long

  ticks, err := C.sysconf(C._SC_CLK_TCK)
  if err != nil {
    panic(err)
  }
  return uint(ticks)
}

