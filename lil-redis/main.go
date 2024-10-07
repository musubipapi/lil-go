package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Cache struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Set(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Optional: Check if key already exists
	if _, exists := c.data[key]; exists {
		return fmt.Errorf("key '%s' already exists", key)
	}

	c.data[key] = value
	return nil
}

func (c *Cache) Get(key string) (*string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.data[key]

	if !exists {
		return nil, false
	}
	return &value, true
}

func handleConnection(conn net.Conn, cache *Cache) {
	defer conn.Close()

	fmt.Println("New connection established") // Add this line

	reader := bufio.NewReader(conn)

	for {
		fmt.Println("Waiting for input...") // Add this line
		input, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				fmt.Printf("Raw input received: %q\n", input) // Add this line
			} else {
				fmt.Printf("Error reading from connection: %s\n", err)
			}
			return
		}

		command := strings.Fields(strings.TrimSpace(input))
		fmt.Printf("Parsed command: %v\n", command) // Updated this line for better logging

		switch command[0] {
		case "SET":
			if len(command) == 3 {
				err := cache.Set(command[1], command[2])
				if err != nil {
					fmt.Fprintf(conn, "ERROR: %s\n", err)
				} else {
					fmt.Fprintf(conn, "OK\n")
				}
			} else {
				fmt.Fprintf(conn, "ERROR: Wrong number of arguments for SET command\n")
			}
		case "GET":
			if len(command) == 2 {
				value, ok := cache.Get(command[1])
				if ok {
					fmt.Fprintln(conn, *value)
				} else {
					fmt.Fprintf(conn, "%s not found\n", command[1])
				}
			} else {
				fmt.Fprintf(conn, "ERROR: Wrong number of arguments for GET command\n")
			}
		default:
			fmt.Fprintln(conn, "ERROR; Unknown command")
		}
	}
}

func main() {
	cache := NewCache()
	port := ":8080"
	listener, err := net.Listen("tcp", port)

	fmt.Printf("Cache listening on %s\n", port) // Add newline for better formatting

	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}
	defer listener.Close()

	for {
		fmt.Println("Waiting for new connection...") // Add this line
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			continue
		}
		fmt.Println("New connection accepted, starting handler") // Add this line
		go handleConnection(conn, cache)
	}
}
