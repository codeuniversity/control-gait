package gait

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/alexmorten/mhist/models"
	mhist "github.com/alexmorten/mhist/proto"
)

// Executor handles movements
type Executor struct {
	mhistClient  mhist.MhistClient
	actionStream mhist.Mhist_StoreStreamClient

	currentMovement Movement
	movementQueue   []Movement
	events          chan struct{}
	*sync.Mutex
}

// NewExecutor ...
func NewExecutor(client mhist.MhistClient) *Executor {
	stream, err := client.StoreStream(context.Background())
	if err != nil {
		panic(err)
	}

	return &Executor{
		mhistClient:  client,
		actionStream: stream,
		Mutex:        &sync.Mutex{},
		events:       make(chan struct{}, 1),
	}
}

// AddNextMovement ...
func (e *Executor) AddNextMovement(movement Movement) {
	e.Lock()
	defer e.Unlock()
	log.Println("adding new movement")
	e.movementQueue = append(e.movementQueue, movement)
	e.events <- struct{}{}
}

// MarkActionDoneFor leg
func (e *Executor) MarkActionDoneFor(legName string, doneCommand string) {
	e.Lock()
	defer e.Unlock()
	for _, leg := range e.currentMovement.legs {
		if leg.name == legName && e.currentMovement.command == doneCommand {
			leg.done = true
			log.Println("marked action done for leg", legName, "and for command", doneCommand)

			e.events <- struct{}{}
			return
		}
	}

	log.Println("couldn't mark action for leg", legName, "and for command", doneCommand, "as done")
}

// Run should be called in a new goroutine
func (e *Executor) Run() {
	for {
		<-e.events

		e.checkState()
	}
}

func (e *Executor) checkState() {
	e.Lock()
	defer e.Mutex.Unlock()

	if e.currentMovement == nil {
		if len(e.movementQueue) > 0 {
			e.currentMovement = e.movementQueue[0]

			e.movementQueue = e.movementQueue[1:len(e.movementQueue)]
		} else {
			return
		}
	}

	e.checkCurrentMovementState()
}

func (e *Executor) checkCurrentMovementState() {
	if e.currentMovement == nil {
		return
	}

	var currentAction *Action = e.currentMovement

	if !currentAction.started {
		log.Println("current Action not started, executing it")
		e.executeAction(currentAction)
		return
	}

	if currentAction.Done() {
		log.Println("current Action done, taking next one")
		e.currentMovement = currentAction.next
		e.checkCurrentMovementState()
	}
}

func (e *Executor) executeAction(action *Action) {
	action.started = true

	for _, leg := range action.legs {
		write(e.actionStream, leg.name, action.command)
	}
}

func write(c mhist.Mhist_StoreStreamClient, legName, message string) {
	log.Println("writing action", message, "to", legName)
	err := c.Send(&mhist.MeasurementMessage{Name: "gait_actions", Measurement: mhist.MeasurementFromModel(&models.Raw{Value: []byte(fmt.Sprintf("%s %s", legName, message))})})

	if err != nil {
		panic(err)
	}
}
