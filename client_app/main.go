package main

import (
	"flag"
	"fmt"

	"github.com/LibenHailu/peer_to_peer_file_share/peer-copy/client_app/server"
)

// var (
// 	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
// )

func main() {
	flag.Parse()
	fmt.Println("start the program")

	// go myServer()
	// go myClient()

	for {
		// start the app
		waitc := make(chan struct{}) // a wait lock

		// start the server thread
		go func() {
			fmt.Println("start the server")
			server.InitFileServer()
			defer close(waitc)
		}()

		// start the client thread
		// go func() {
		// 	// for {
		// 	serverAddr, server := filesource.SearchAddressForThefile("Liben.jpg")
		// 	fmt.Println(*serverAddr)
		// 	fmt.Println(*server)
		// 	client.InitFileClient(serverAddr, server)
		// 	client.DownloadFile("Liben.jpg")
		// 	// }
		// }()

		// start the input thread
		// go input()

		<-waitc
		// finished in this round restart the app
		fmt.Println("restart the app")
	}
}
