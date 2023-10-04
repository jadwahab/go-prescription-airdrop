# go-prescription-airdrop

## Overview

This Go application reads a JSON file containing a list of prescriptions, performs an airdrop of these prescriptions, and then writes the results to two separate files: one for successful airdrops and another for unsuccessful ones.

## Requirements

- Go version 1.16 or higher

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/go-prescription-airdrop.git
```

Navigate to the project directory:

```bash
cd go-prescription-airdrop
```

## Usage

To run the program, execute:

```bash
go run main.go
```

If you want to only run a few at a time (for example only 100) instead of all 5000 or so, just change line 25 in main.go to
```go
	for _, persc := range perscList[:100] {
```

## Output

Two files will be generated:

- `perscListSuccess.json`: Contains the list of successful airdrops.
- `perscList.json`: Contains the list of unsuccessful airdrops that still need to be done.

## License

MIT License
