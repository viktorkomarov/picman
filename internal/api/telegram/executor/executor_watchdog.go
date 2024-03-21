package executor

import (
	"fmt"

	"github.com/viktorkomarov/picman/internal/api/telegram"
	"github.com/viktorkomarov/picman/internal/api/telegram/usecases/provider"
)

type ExecutorWatchdog struct {
	inMessages  chan telegram.UserEvent
	fsmProvider *provider.FSMBuilder

	executor *UseCaseExecution
	// shouldn't close this channel
	terminate <-chan error
	// should close this channel :с
	outMessage chan telegram.UserEvent
}

func NewExecutorWatchdog(fsmProvider *provider.FSMBuilder) *ExecutorWatchdog {
	watchdog := &ExecutorWatchdog{
		inMessages:  make(chan telegram.UserEvent),
		fsmProvider: fsmProvider,
	}
	go watchdog.loop()

	return watchdog
}

func (e *ExecutorWatchdog) RecieveMessage(msg telegram.UserEvent) {
	e.inMessages <- msg
}

func (e *ExecutorWatchdog) loop() {
	defer func() {
		close(e.outMessage)
	}()

	for {
		select {
		case err, ok := <-e.terminate:
			if ok {
				fmt.Println(err)
				return
			}
		case msg, ok := <-e.inMessages:
			if !ok {
				return
			}

			if msg.IsCommand() {
				e.spawnExecutionContext(msg)
			} else {
				e.outMessage <- msg
			}
		}
	}
}

func (e *ExecutorWatchdog) spawnExecutionContext(event telegram.UserEvent) {
	e.outMessage = make(chan telegram.UserEvent)
	e.executor = NewUserExecution(e.fsmProvider.GetFSMByCommandType(event.RawMessage.Text), e.outMessage, event.RawMessage.Chat.ID)
	e.terminate = e.executor.Run()
}
