package main

import (
    "bufio"
    "fmt"
    "net"
    "sync"
)

// TcpListener starts a TCP server on a given address and port.
// It accepts incoming connections and passes them to the TcpHandler.
func TcpListener(address string) {
    ln, err := net.Listen("tcp", address)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        return
    }
    defer ln.Close()
    fmt.Printf("TCP is listening on %s\n", address)

    var wg sync.WaitGroup

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting:", err.Error())
            continue
        }
        wg.Add(1)
        fmt.Printf("New client accepted: %s\n", conn.RemoteAddr().String())

        // Each connection is handled in its own goroutine
        go func() {
            defer wg.Done()
            TcpHandler(conn)
        }()
    }

    // Wait for all handlers to finish
    wg.Wait()
}

// TcpHandler handles the individual client connection.
// It reads from the connection and writes back the same data.
func TcpHandler(conn net.Conn) {
    defer conn.Close()
    remoteAddr := conn.RemoteAddr().String()
    fmt.Printf("Handle Request from [%s]\n", remoteAddr)

    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            fmt.Printf("Client [%s] Error: %s\n", remoteAddr, err.Error())
            break
        }
        fmt.Printf("Received from [%s]: %s", remoteAddr, message)

        // Echoing back the message to the client
        _, err = conn.Write([]byte(message))
        if err != nil {
            fmt.Printf("Error sending to [%s]: %s\n", remoteAddr, err.Error())
            break
        }
    }
}

func main() {
    TcpListener("127.0.0.1:8080")
}
