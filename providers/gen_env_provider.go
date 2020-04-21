package providers

import "os"

type EnvGetter interface {
	Getenv(key string) string
}

type OsEnvGetter struct{}

func NewOsEnvGetter() *OsEnvGetter {
	return new(OsEnvGetter)
}

func (b *OsEnvGetter) Getenv(key string) string {
	return os.Getenv(key)
}
