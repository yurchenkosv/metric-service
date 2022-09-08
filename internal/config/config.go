package config

type Config interface {
	Parse() error
}
