package netx

import "time"

type options struct {
	Reconnect    bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Option func(*options)

func WithReconnect(reconnect bool) Option {
	return Option(func(os *options) {
		os.Reconnect = reconnect
	})
}

func WithReadTimeout(t time.Duration) Option {
	return Option(func(os *options) {
		os.ReadTimeout = t
	})
}

func WithWriteTimeout(t time.Duration) Option {
	return Option(func(os *options) {
		os.WriteTimeout = t
	})
}
