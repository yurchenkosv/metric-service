package config

// Config interface.
type Config interface {
	Parse() error // method that must fulfill fields of struct, implementing interface
}
