package godash

import (
	"time"
)

type DebounceConfig struct {
	wait    time.Duration
	maxWait time.Duration
}

type Option func(conf *DebounceConfig)

type Debounced func()

func DebounceMaxWaitOption(maxWait time.Duration) Option {
	return func(conf *DebounceConfig) {
		conf.maxWait = maxWait
	}
}

func Debounce(cb func(), wait time.Duration, options ...Option) (debounced Debounced, quitChan chan<- struct{}) {
	var conf DebounceConfig
	conf.wait = wait
	for _, option := range options {
		option(&conf)
	}
	callChan := make(chan struct{}, 5)
	qChan := make(chan struct{})
	quitChan = qChan
	lastInvokeAt := time.Now()
	var lastCallAt *time.Time
	var canInvoke bool
	var invokeCB func()
	invokeCB = func() {
		canInvoke = false
		lastInvokeAt = time.Now()
		cb()
	}
	var invokeInterval func() time.Duration
	invokeInterval = func() time.Duration {
		du := conf.wait
		now := time.Now()
		if lastCallAt != nil {
			callElapsed := now.Sub(*lastCallAt)
			if callElapsed >= conf.wait {
				return 0
			}
			du -= callElapsed
		}
		if conf.maxWait > 0 {
			invokeElapsed := now.Sub(lastInvokeAt)
			maxWaitExp := conf.maxWait - invokeElapsed
			if maxWaitExp <= 0 {
				return 0
			}
			if maxWaitExp < du {
				du = maxWaitExp
			}
		}
		return du
	}
	var touch func()
	touch = func() {
		canInvoke = true
		now := time.Now()
		lastCallAt = &now
	}
	go func() {
	LOOP:
		for {
			if canInvoke {
				select {
				case <-qChan:
					break LOOP
				case <-time.After(invokeInterval()):
					invokeCB()
				case <-callChan:
					touch()
				}
			} else {
				select {
				case <-qChan:
					break LOOP
				case <-callChan:
					touch()
				}
			}
		}
	}()
	debounced = func() {
		callChan <- struct{}{}
	}
	return
}
