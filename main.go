package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

const protocolID = "/simple-chat/1.0.0"

func main() {
	ctx := context.Background()

	// Create a new libp2p host
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/192.168.7.109/tcp/0"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Display this node's multiaddress
	fmt.Println("Node ID:", host.ID())
	fmt.Println("Node Address:", host.Addrs()[0])

	// Set up stream handler for incoming connections
	host.SetStreamHandler(protocol.ID(protocolID), handleStream)

	// Connect to another peer if provided
	if len(os.Args) > 1 {
		peerAddr := os.Args[1]
		connectToPeer(ctx, host, peerAddr)
	} else {
		fmt.Println("Run this program with another peer's address as an argument to connect.")
	}

	// Keep the program running
	select {}
}

func handleStream(stream network.Stream) {
	fmt.Println("Connected to:", stream.Conn().RemotePeer())

	// Read incoming messages from the stream
	buf := bufio.NewReader(stream)
	for {
		msg, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed.")
			return
		}
		fmt.Printf("Received: %s", msg)
	}
}

func connectToPeer(ctx context.Context, host host.Host, peerAddr string) {
	// Parse the peer address
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Get peer info from the multiaddress
	peerinfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the peer
	if err := host.Connect(ctx, *peerinfo); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to peer:", peerinfo.ID)

	// Open a stream to the peer
	stream, err := host.NewStream(ctx, peerinfo.ID, protocol.ID(protocolID))
	if err != nil {
		log.Fatal(err)
	}

	// Send a message every few seconds
	go func() {
		for {
			_, err := stream.Write([]byte("Hello from " + host.ID() + "\n"))
			if err != nil {
				log.Println("Error sending message:", err)
				return
			}
			time.Sleep(2 * time.Second)
		}
	}()
}
