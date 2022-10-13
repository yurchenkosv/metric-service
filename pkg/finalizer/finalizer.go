package finalizer

import "os"

// Shutdown is function to make main methods compliant to analyzer method.
// It wraps function (most common some final function) and then call os.Exit to close program.
func Shutdown(fn func()) {
	fn()
	os.Exit(0)
}
