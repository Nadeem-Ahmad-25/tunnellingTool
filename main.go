package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

type TunnelServer struct {
	localPort   int
	publicPort  int
	connections sync.WaitGroup
	stopChan    chan struct{}
}

func NewTunnelServer(localPort int) *TunnelServer {
	return &TunnelServer{
		localPort: localPort,
		stopChan:  make(chan struct{}),
	}
}

func (ts *TunnelServer) Start() error {
	// Getting your local machine's IP address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return fmt.Errorf("failed to get network interfaces: %v", err)
	}

	var hostIP string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				hostIP = ipnet.IP.String()
				break
			}
		}
	}

	//This is like setting up a public phone line on a random available extension
	publicListener, err := net.Listen("tcp", ":0")
	if err != nil {
		return fmt.Errorf("failed to create public listener: %v", err)
	}

	_, portStr, err := net.SplitHostPort(publicListener.Addr().String())
	if err != nil {
		return fmt.Errorf("failed to parse listener address: %v", err)
	}

	ts.publicPort, _ = strconv.Atoi(portStr)

	log.Printf(" Tunnel created successfully!")
	log.Printf(" Local Server: localhost:%d", ts.localPort)
	log.Printf(" Public Endpoint: %s:%d", hostIP, ts.publicPort)

	go ts.acceptConnections(publicListener)
	return nil
}
func (ts *TunnelServer) acceptConnections(publicListener net.Listener) {
	for {
		select {
		case <-ts.stopChan:
			publicListener.Close()
			return
		default:
			// Accept incoming public connections
			publicConn, err := publicListener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}

			ts.connections.Add(1)
			go ts.handleConnection(publicConn)
		}
	}
}

func (ts *TunnelServer) handleConnection(publicConn net.Conn) {
	defer ts.connections.Done()
	defer publicConn.Close()

	// Establish connection to local server
	localConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", ts.localPort))
	if err != nil {
		log.Printf("Failed to connect to local server: %v", err)
		return
	}
	defer localConn.Close()

	// Bidirectional data transfer
	errChan := make(chan error, 2)

	go func() {
		_, err := io.Copy(localConn, publicConn)
		errChan <- err
	}()

	go func() {
		_, err := io.Copy(publicConn, localConn)
		errChan <- err
	}()

	// Wait for either copy operation to complete or fail
	<-errChan
}

func (ts *TunnelServer) Stop() {
	close(ts.stopChan)
	ts.connections.Wait()
	log.Println("🛑 Tunnel stopped")
}

func main() {
	// Get local port from command line or use default
	localPort := 3000
	if len(os.Args) > 1 {
		port, err := strconv.Atoi(os.Args[1])
		if err == nil {
			localPort = port
		}
	}

	// Create and start tunnel
	tunnel := NewTunnelServer(localPort)

	if err := tunnel.Start(); err != nil {
		log.Fatalf("Failed to start tunnel: %v", err)
	}

	// Keep the program running until interrupted
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	tunnel.Stop()
}
