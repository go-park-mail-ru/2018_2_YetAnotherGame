package grpcModules

import (
	"2018_2_YetAnotherGame/authorization/api"

	"golang.org/x/net/context"
	//"fmt"
	"log"

	"google.golang.org/grpc"
)

func SendCheckInfo(id string, conn *grpc.ClientConn) string {

	c := api.NewPingClient(conn)
	ctx := context.Background()
	response, err := c.CheckSession(ctx, &api.PingMessage{Message: id})
	if err != nil {
		log.Fatalf("Error when calling CheckSession: %s", err)
	}
	log.Printf("Response from server: %s", response.Message)
	return response.Message
}
