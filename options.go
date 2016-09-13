package thinkrus

import "time"

// Option configures the Hook
type Option func(*RethinkHook)

// WithBatchInterval sets the batchInterval
func WithBatchInterval(interval int) Option {
	return func(o *RethinkHook) {
		o.batchInterval = time.Duration(interval) * time.Second
	}
}

// WithBatchSize sets the batchSize
func WithBatchSize(count int) Option {
	return func(o *RethinkHook) {
		o.batchSize = count
	}
}
