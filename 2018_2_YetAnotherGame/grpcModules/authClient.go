package grpcModules

import (
	"authorization/api"
	"golang.org/x/net/context"
	//"fmt"
	"google.golang.org/grpc"
	"log"
)

func SendCheckInfo(id string, conn *grpc.ClientConn) string {

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	c := api.NewPingClient(conn)
	ctx := context.Background()
	response, err := c.CheckSession(ctx, &api.PingMessage{Message: id})
	if err != nil {
		log.Fatalf("Error when calling CheckSession: %s", err)
	}
	log.Printf("Response from server: %s", response.Message)
	return response.Message
}
