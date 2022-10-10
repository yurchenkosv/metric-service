package finalizer

import "os"

func Shutdown(fn func()) {
	fn()
	os.Exit(0)
}
