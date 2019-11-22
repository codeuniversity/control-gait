package gait

import (
	"context"
	"log"

	mhist "github.com/alexmorten/mhist/proto"
)

// SubscribeToHighLevelAction ...
func SubscribeToHighLevelAction(c mhist.MhistClient, handler func(movement Movement)) {
	stream, err := c.Subscribe(context.Background(), &mhist.Filter{Names: []string{"hlc_actions"}})
	if err != nil {
		panic(err)
	}

	for {
		_, err := stream.Recv()
		log.Println("bla")
		if err != nil {
			panic(err)
		}

		// handler(NewMovement([]string{"leg1", "leg3", "leg5"}))
		// handler(NewMovement([]string{"leg2", "leg4", "leg6"}))

		handler(NewMovement([]string{"leg1"}))
	}
}
