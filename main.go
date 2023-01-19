package main

import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"github.com/lightningnetwork/lnd/lnrpc"
)

func main() {

	//Call connectToServer function with url of the server every minute
	for {
		connectToServer("localhost:10001")
		time.Sleep(1 * time.Second)
	}
}

// func that receives a url and connects to a grpc server
func connectToServer(url string) {

	log.Printf("Connecting to %s", url)
	//Generate TLS credentials from directory
	creds, err := credentials.NewClientTLSFromFile("/Users/joseap/.polar/networks/1/volumes/lnd/alice/tls.cert", "")

	if err != nil{
		log.Fatalf("Failed to load credentials: %v", err)
	}

	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	lightningClient := lnrpc.NewLightningClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"macaroon": "0201036c6e6402f801030a108be5b2928f746a822b04a9b2848eb0321201301a160a0761646472657373120472656164120577726974651a130a04696e666f120472656164120577726974651a170a08696e766f69636573120472656164120577726974651a210a086d616361726f6f6e120867656e6572617465120472656164120577726974651a160a076d657373616765120472656164120577726974651a170a086f6666636861696e120472656164120577726974651a160a076f6e636861696e120472656164120577726974651a140a057065657273120472656164120577726974651a180a067369676e6572120867656e6572617465120472656164000006208e8b02d4bc0efd4f15a52946c5ef23f2954f8a07ed800733554a11a190cb71b4"})

	context := metadata.NewOutgoingContext(ctx, md)

	//Call ListChannels method of lightning client with metadata headers
	response, err := lightningClient.ListChannels(context, &lnrpc.ListChannelsRequest{
		ActiveOnly: false,
	})

	if err != nil {
		log.Fatalf("ListChannels fail", err)
	}

	//Iterate over response channels
	for _, channel := range response.Channels {
		//Print channel id
		log.Printf("Channel ID: %d", channel.ChanId)

		//Get the balance of the channel
		localBalance := channel.LocalBalance
		remoteBalance := channel.RemoteBalance

		//Print the balance of the channel
		log.Printf("Local Balance: %d", localBalance)
		log.Printf("Remote Balance: %d", remoteBalance)
	}


}
