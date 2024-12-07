package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var port string
var conn net.Conn

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Go-Zeroes",
	Short: "CLI application for Go-Zeroes",
	Long:  `Go-Zeroes is a CLI tool to connect to a Go-Zeroes server and interact with it.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Establish the connection to the server
		conn = connectToServer()
		// Start the interactive prompt for commands
		startCommandPrompt(conn)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd.Flags().StringVarP(&port, "port", "p", "8769", "Port to connect to the Go-Zeroes server")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func connectToServer() net.Conn {
	var err error
	conn, err = net.Dial("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	fmt.Printf("Connected to Go-Zeroes server at port %s\n", port)
	return conn
}

func startCommandPrompt(conn net.Conn) {
	fmt.Println("Type 'exit' to quit the CLI.")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Go-Zeroes> ")
		cmd, _ := reader.ReadString('\n')

		cmd = strings.TrimSpace(cmd)

		if cmd == "exit" {
			fmt.Println("Exiting Go-Zeroes CLI.")
			conn.Close()
			break
		}

		tokens := tokenize(cmd)
		fmt.Println(tokens)

		// Send the command to the server
		// sendCommandToServer(cmd, conn)

		// handleServerResponse(conn)
	}
}

func sendCommandToServer(cmd string, conn net.Conn) {

	_, err := conn.Write([]byte(cmd + "\r\n"))
	if err != nil {
		fmt.Println("Error sending command to the server:", err)
		return
	}
}

func handleServerResponse(conn net.Conn) {

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading response from server:", err)
		return
	}

	fmt.Println("Server Response:", string(buffer[:n]))
}

func tokenize(input string) []string {
	var tokens []string
	var currentToken strings.Builder
	inQuotes := false

	for _, ch := range input {

		if ch == '"' {
			inQuotes = !inQuotes
			if !inQuotes && currentToken.Len() > 0 {

				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		} else if ch == ' ' && !inQuotes {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		} else {

			currentToken.WriteRune(ch)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}
