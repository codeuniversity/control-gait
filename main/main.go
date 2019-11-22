package main

import (
	"flag"

	mhist "github.com/alexmorten/mhist/proto"
	gait "github.com/codeuniversity/control-gait"
	"google.golang.org/grpc"
)

func main() {
	var mhistAddress string
	flag.StringVar(&mhistAddress, "mhist_address", "localhost:6666", "address to mhist including port")
	flag.Parse()
	mhistConn, err := grpc.Dial(mhistAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	mhistC := mhist.NewMhistClient(mhistConn)

	executor := gait.NewExecutor(mhistC)
	go gait.SubscribeToFeedback(mhistC, executor.MarkActionDoneFor)
	go gait.SubscribeToHighLevelAction(mhistC, executor.AddNextMovement)

	executor.Run()
}
