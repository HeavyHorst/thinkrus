package thinkrus

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	r "gopkg.in/dancannon/gorethink.v2"
)

type RethinkHook struct {
	session       *r.Session
	table         string
	batchInterval time.Duration
	batchSize     int
	batchChan     chan interface{}
	flushChan     chan struct{}
	flushed       chan struct{}
	err           error
	errLock       sync.RWMutex
}

func New(url, db, table string, opts ...Option) (*RethinkHook, error) {
	var err error

	hook := &RethinkHook{
		table:         table,
		flushChan:     make(chan struct{}),
		flushed:       make(chan struct{}),
		err:           nil,
		errLock:       sync.RWMutex{},
		batchSize:     200,
		batchInterval: 5 * time.Second,
	}

	for _, o := range opts {
		o(hook)
	}

	if hook.session == nil {
		hook.session, err = r.Connect(r.ConnectOpts{
			Address:  url,
			Database: db,
		})
		if err != nil {
			return nil, err
		}
	}

	hook.batchChan = make(chan interface{}, hook.batchSize)

	go func() {
		batch := make([]interface{}, 0, hook.batchSize)
		ticker := time.NewTicker(hook.batchInterval)

		flushAndClear := func() {
			_, err := r.Table(hook.table).Insert(batch).RunWrite(hook.session)
			hook.errLock.Lock()
			hook.err = err
			hook.errLock.Unlock()

			// only clear the buffer if all data is written to the server
			if err == nil {
				batch = batch[0:0]
			}
		}

		for {
			select {
			case <-ticker.C:
				flushAndClear()
			case b := <-hook.batchChan:
				batch = append(batch, b)
				if len(batch) >= hook.batchSize {
					flushAndClear()
				}
			case <-hook.flushChan:
				for p := range hook.batchChan {
					batch = append(batch, p)
				}
				flushAndClear()
				hook.flushed <- struct{}{}
			}
		}
	}()

	return hook, nil
}

func (h *RethinkHook) Close() {
	close(h.batchChan)
	h.flushChan <- struct{}{}
	<-h.flushed
}

func (h *RethinkHook) Fire(entry *logrus.Entry) error {
	if errData, ok := entry.Data[logrus.ErrorKey]; ok {
		if err, ok := errData.(error); ok && entry.Data[logrus.ErrorKey] != nil {
			entry.Data[logrus.ErrorKey] = err.Error()
		}
	}

	data := make(map[string]interface{})
	for k, v := range entry.Data {
		data[k] = v
	}

	data["level"] = entry.Level.String()
	data["time"] = entry.Time
	data["message"] = entry.Message
	h.batchChan <- data

	h.errLock.RLock()
	err := h.err
	h.errLock.RUnlock()
	return err
}

func (h *RethinkHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
