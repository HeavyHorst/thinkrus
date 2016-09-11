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

// WithBatchInterval sets the batchCount
func WithBatchCount(count int) Option {
	return func(o *RethinkHook) {
		o.batchCount = count
	}
}
