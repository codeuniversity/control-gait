package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/alexmorten/mhist/models"
	mhist "github.com/alexmorten/mhist/proto"

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

	go moveLeg(mhistC, "leg1")
	go moveLeg(mhistC, "leg2")
	go moveLeg(mhistC, "leg3")
	go moveLeg(mhistC, "leg4")
	go moveLeg(mhistC, "leg5")
	go moveLeg(mhistC, "leg6")
}

func moveLeg(c mhist.MhistClient, portName string) {
	stream, err := c.StoreStream(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		write(stream, portName, "move up")
		time.Sleep(time.Second / 10)
		write(stream, portName, "rotate 30")
		time.Sleep(time.Second / 10)
		write(stream, portName, "move down")
		time.Sleep(time.Second / 10)
		write(stream, portName, "rotate -20")
		time.Sleep(time.Second / 10)
	}
}

func write(c mhist.Mhist_StoreStreamClient, legName, message string) {
	err := c.Send(
		&mhist.MeasurementMessage{Name: "gait_actions", Measurement: mhist.MeasurementFromModel(&models.Raw{
			Value: []byte(fmt.Sprintf("%s %s", legName, message)),
		})},
	)

	if err != nil {
		panic(err)
	}
}
