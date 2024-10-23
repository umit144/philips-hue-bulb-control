# Philips Hue Bulb Control

A CLI application to control Philips Hue smart bulbs using the Hue Bridge API.

## Features

- List all connected bulbs
- Toggle bulbs on/off
- Interactive TUI mode
- Command line mode for scripting

## Prerequisites

- Go 1.21 or higher
- Philips Hue Bridge
- Bridge API username/token
- Network connectivity to your Hue Bridge

## Installation

```bash
# Clone the repository
git clone https://github.com/umit144/philips-hue-bulb-control.git

# Enter the directory
cd philips-hue-bulb-control

# Install dependencies
go mod download
```

## Configuration

Create a `.env` file in the project root:

```env
HUE_BRIDGE_IP=<your-bridge-ip>
HUE_USERNAME=<your-username>
```

## Usage

### Interactive Mode

Run the application in interactive mode:

```bash
go run main.go
```

### Command Line Mode

Toggle a specific light:

```bash
go run main.go -light=1 -state=on
go run main.go -light=1 -state=off
```

## Development

### Running Tests

```bash
make test
```

### Building

```bash
make build
```

## Project Structure

```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── app/
│   ├── config/
│   ├── light/
│   └── ui/
├── tests/
│   ├── app/
│   ├── config/
│   ├── light/
│   └── ui/
├── .env
├── .gitignore
├── go.mod
└── README.md
```