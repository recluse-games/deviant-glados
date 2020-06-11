package main

import (
	"context"
	"io"
	"log"
	"time"

	channels "github.com/eapache/channels"

	"github.com/recluse-games/deviant-glados/hunting"
	deviant "github.com/recluse-games/deviant-protobuf/genproto/go"
	"google.golang.org/grpc"
)

// HACK: REMOVE THIS ONCE WE HAVE PROPER REGISTRATION.
func createEncounter() *deviant.EncounterRequest {
	encounterRequest := &deviant.EncounterRequest{}
	encounterRequest.EncounterCreateAction = &deviant.EncounterCreateAction{}
	encounterRequest.PlayerId = "0002"

	return encounterRequest
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := deviant.NewEncounterServiceClient(conn)
	ch := channels.NewRingChannel(0) // yes this is rather silly, but it should work

	stream, err := client.UpdateEncounter(context.Background())
	// HACK Used to add the bot to the players table.
	if err := stream.Send(createEncounter()); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				return
			}

			ch.In() <- in

			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
		}
	}()
	go func() {
		for {
			singleEncounterRes := <-ch.Out()
			if singleEncounterRes != nil {
				if singleEncounterRes.(*deviant.EncounterResponse).Encounter.ActiveEntity.OwnerId == "0002" {
					log.Printf("Current Active Entity %v", singleEncounterRes.(*deviant.EncounterResponse).Encounter.ActiveEntity.Id)

					for _, request := range hunting.TakeTurn(singleEncounterRes.(*deviant.EncounterResponse)) {
						log.Printf("Sending Request: %v", request)
						time.Sleep(500 * time.Millisecond)
						if err := stream.Send(request); err != nil {
							log.Fatalf("Failed to send a note: %v", err)
						}
					}
				}
			}
		}
	}()

	<-waitc
}
