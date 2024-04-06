package application

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Closer struct {
	mutex sync.Mutex
	once  sync.Once
	done  chan struct{}

	gracefulCloseTimeout  time.Duration
	log                   *slog.Logger
	gracefulShutdownFuncs []func() error
}

func NewCloser(log *slog.Logger, gracefulCloseTimeout time.Duration, signals ...os.Signal) *Closer {
	closer := &Closer{
		mutex:                 sync.Mutex{},
		once:                  sync.Once{},
		done:                  make(chan struct{}),
		gracefulCloseTimeout:  gracefulCloseTimeout,
		log:                   log,
		gracefulShutdownFuncs: nil,
	}

	if len(signals) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, signals...)
			<-ch
			signal.Stop(ch)
			closer.CloseEverything()
		}()
	}
	return closer
}

func (c *Closer) Run(funcs ...func() error) {
	for _, f := range funcs {
		go func(f func() error) {
			if err := f(); err != nil {
				c.log.Error("closer functions error", "error", err)
				c.CloseEverything()
			}
		}(f)
	}
}

func (c *Closer) Add(f ...func() error) {
	c.mutex.Lock()
	c.gracefulShutdownFuncs = append(c.gracefulShutdownFuncs, f...)
	c.mutex.Unlock()
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) CloseEverything() {
	c.once.Do(func() {
		defer close(c.done)

		c.log.Info("Start graceful shutdowns")
		if ok := c.waitAllCloseFuncs(); !ok {
			c.log.Error("graceful shutdown error: timeout limit exceed")
		}
	})
}

func (c *Closer) waitAllCloseFuncs() bool {
	c.mutex.Lock()
	gracefulFuncs := c.gracefulShutdownFuncs
	c.mutex.Unlock()

	wg := &sync.WaitGroup{}
	for _, f := range gracefulFuncs {
		wg.Add(1)
		go func(f func() error) {
			if err := f(); err != nil {
				c.log.Error("close function error", "error", err)
			}
			wg.Done()
		}(f)
	}

	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	timer := time.NewTimer(c.gracefulCloseTimeout)
	defer timer.Stop()

	select {
	case <-ch:
		return true
	case <-timer.C:
		return false
	}
}
