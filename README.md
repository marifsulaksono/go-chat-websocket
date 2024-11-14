# WebSocket Chat by Muhammad Arif Sulaksono

A simple WebSocket-based chat application built with HTML, CSS, and JavaScript for the frontend, and Go for the backend.

## Features

* Real-time messaging between users
* Connection establishment and disconnection handling
* Message broadcasting and receiving
* Chat history display and clearing

## Getting Started

1. Clone the repository to your local machine. `git clone https://github.com/marifsulaksono/go-chat-websocket.git`.
2. Open project directory `cd go-chat-websocket`.
3. Run the Go server by executing `go run main.go` in the terminal.
4. Open multiple browser tabs and navigate to `index.html` file.
5. Enter sender and receiver IDs, and connect to the chat.
6. Send messages to other users and see them appear in real-time.

## Project Structure

* `index.html`: The main HTML file for the chat application.
* `main.go`: The Go server entry point.
* `chat/socket.go`: The Go file that defines the WebSocketManager and Client structs.

## Dependencies

* Go (version 1.17 or higher)
* Gorilla WebSocket (version 1.4.2 or higher)
* LabStack Echo (version 4.6.1 or higher)