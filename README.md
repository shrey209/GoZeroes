# GoZeroes

GoZeroes is a minimalist Redis-like in-memory database clone implemented in Go. It serves as a basic demonstration of event loop architecture and utilizes the Linux `epoll` API for efficient I/O operations. This project supports basic Redis-like commands and is designed to run on Linux-based systems or WSL.

## Features

### Supported Commands
- **String operations**
  - `GET` - Retrieve the value associated with a key
  - `SET` - Assign a value to a key
  - `DEL` - Delete a key and its associated value
  - `INCRBY` - Increment a key's value by a specified integer
  - `DECRBY` - Decrement a key's value by a specified integer
  - `INCR` - Increment a key's value by 1
  - `DECR` - Decrement a key's value by 1

- **Sorted Sets** (Work in progress)
  - Data structure and logic implemented
  - Commands for interacting with sorted sets are not fully completed yet

## Prerequisites

- **Operating System:** Ubuntu, any Linux-based OS, or WSL (Windows Subsystem for Linux)
- **Go Version:** Go 1.20+

## Getting Started

### 1. Clone the Repository
```bash
git clone <repository-url>
cd gozeroes
```

### 2. Run the Server
Navigate to the `server` folder and start the server:
```bash
cd server
go run main.go
```

### 3. Run the CLI
To interact with the server, you can use the CLI application. Navigate to the CLI folder and run the CLI binary or its source code:

#### Using the Binary
```bash
./go-zeroes-cli
```

#### Using the Source Code
```bash
cd cli
go run main.go
```

## Usage

Once the server and CLI are running, you can interact with the database using supported commands like `GET`, `SET`, `INCR`, `DEL`, etc. Examples:

```bash
> SET key1 value1
OK
> GET key1
value1
> INCR counter
1
> INCRBY counter 10
11
> DEL key1
OK
```

## Notes

- **Linux Dependency:** This project relies on the Linux `epoll` API. It must be run on a Linux-based OS or WSL to function correctly.
- **Work in Progress:** Sorted sets are partially implemented, with the data structure and logic added, but command functionality is not yet complete.

## Contribution
Feel free to fork the repository and contribute to the project. Bug reports, feature requests, and pull requests are welcome!

## License
This project is licensed under the MIT License.
