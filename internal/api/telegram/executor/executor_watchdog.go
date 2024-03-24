package executor

import (
	"fmt"
	"time"

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
	return &ExecutorWatchdog{
		inMessages:  make(chan telegram.UserEvent),
		fsmProvider: fsmProvider,
	}
}

func (e *ExecutorWatchdog) RecieveMessage(msg telegram.UserEvent) {
	e.inMessages <- msg
}

func (e *ExecutorWatchdog) Loop(keepAlive time.Duration) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer func() {
			close(done)
			e.closeExuctionContext()
		}()

		ticker := time.NewTicker(keepAlive)
		defer ticker.Stop()
		lastUpdated := time.Now()

		for {
			select {
			case <-ticker.C:
				if time.Since(lastUpdated) > keepAlive {
					return
				}
			case err, ok := <-e.terminate:
				lastUpdated = time.Now()
				if ok {
					fmt.Println(err)
				}
				e.closeExuctionContext()
			case msg, ok := <-e.inMessages:
				lastUpdated = time.Now()
				if !ok {
					return
				}

				if msg.IsCommand() {
					e.spawnExecutionContext(msg, msg.RawMessage.Text)
				} else if e.executor == nil {
					e.spawnExecutionContext(msg, "fallback")
				} else {
					e.outMessage <- msg
				}
			}
		}
	}()

	return done
}

func (e *ExecutorWatchdog) spawnExecutionContext(event telegram.UserEvent, command string) {
	e.closeExuctionContext()
	e.outMessage = make(chan telegram.UserEvent)
	e.executor = NewUserExecution(e.fsmProvider.GetFSMByCommandType(command), e.outMessage, event.RawMessage.Chat.ID)
	e.terminate = e.executor.Run()
}

func (e *ExecutorWatchdog) closeExuctionContext() {
	if e.executor == nil {
		return
	}
	close(e.outMessage)
	e.executor = nil
	e.terminate = nil
}
