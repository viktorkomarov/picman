package executor

import "github.com/viktorkomarov/picman/internal/api/telegram"

type ExecutorWatchdog struct {
	inMessages chan telegram.Message
	terminate  chan struct{}

	executor    *UseCaseExecution
	outMessage  chan telegram.Message
	fsmProvider interface{}
}

func NewExecutorWatchdog(fsmProvider interface{}) *ExecutorWatchdog {
	watchdog := &ExecutorWatchdog{
		inMessages:  make(chan telegram.Message),
		terminate:   make(chan struct{}),
		executor:    nil,
		fsmProvider: fsmProvider,
	}
	go watchdog.loop()

	return watchdog
}

func (e *ExecutorWatchdog) RecieveMessage(msg telegram.Message) {
	e.inMessages <- msg
}

func (e *ExecutorWatchdog) Terminate() {
	e.terminate <- struct{}{}
}

func (e *ExecutorWatchdog) loop() {
	defer func() {
		close(e.terminate)
		close(e.outMessage)
	}()

	for {
		select {
		case <-e.terminate:
		case msg, ok := <-e.inMessages:
			if !ok {
				return
			}

			if msg.IsCommand() {
				e.outMessage = make(chan telegram.Message)
				e.executor = NewUserExecution(nil, nil, e.outMessage)
			} else {
				e.outMessage <- msg
			}
		}
	}
}
