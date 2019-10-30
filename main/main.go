package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/codeuniversity/nervo/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(os.Args[1], grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := proto.NewNervoServiceClient(conn)

	response, err := c.ListControllers(
		context.Background(),
		&proto.ControllerListRequest{},
	)

	if err != nil {
		panic(err)
	}

	for index, info := range response.ControllerInfos {
		fmt.Println(index, info.Name, info.PortName)
	}

	portName := findLeg1PortName(response.ControllerInfos)
	if portName == "" {
		panic("Oh fuck, no port name")
	}
	moveLeg(c, portName)
}

func findLeg1PortName(infos []*proto.ControllerInfo) string {
	for _, info := range infos {
		if info.Name == "leg1" {
			return info.PortName
		}
	}

	return ""
}

func moveLeg(c proto.NervoServiceClient, portName string) {
	write(c, portName, "lift up \n")
	time.Sleep(1 * time.Second)
	write(c, portName, "move forward \n")
	time.Sleep(1 * time.Second)
	write(c, portName, "lift down \n")
	time.Sleep(1 * time.Second)
	write(c, portName, "move back")
}
func write(c proto.NervoServiceClient, portName, message string) {
	_, err := c.WriteToController(context.Background(), &proto.WriteToControllerRequest{
		ControllerPortName: portName, Message: []byte(message),
	})
	if err != nil {
		panic(err)
	}
}
