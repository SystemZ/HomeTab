package main

// #cgo LDFLAGS: -lX11 -lXmu
// #include <activewindow.h>
import "C"
import "log"

func getActiveWindow() {
	log.Printf("%v", C.GoString(C.activeWindowName()))
}
