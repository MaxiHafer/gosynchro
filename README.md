# gosynchro

gosynchro is a slim alternative to browser-sync, designed to synchronize browser windows for development. It utilizes server sent events for real-time communication between the server and the client.

## Installation
To install `gosynchro`, you can use the `go install` command.
```bash
go install github.com/maxihafer/gosynchro@latest
```

## Usage
`gosynchro` is a command-line tool with several commands and options. Here's a brief overview:

### Commands

- `proxy`: Starts a proxy server listening on `PORT`. It proxies all requests to the remote set using the `--remote` flag (default: `http://localhost:8080`).
- `reload`: Reloads the client listening on `PORT`.
- `help, h`: Shows a list of commands or help for one command.

### Global Options

- `--port value, -p value`: Specifies the port for the server to listen on. The default value is 3000.
- `--verbose, -v`: Enables verbose logging. This is disabled by default.
- `--json`: Outputs logs in JSON format. This is disabled by default.
- `--help, -h`: Shows help information.

Here's an example of how to use `gosynchro`:

```bash
# Start a proxy server on port 3000 proxying all requests to `http://localhost:8080`.
gosynchro proxy

# Reload the client listening on the default port (3000)
gosynchro reload
