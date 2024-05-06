# Replay
A shell command composer. Allows commands from your shell history to be composed and replayed in the order of your choosing.

## Installation
### From Binary
Download and install the [latest release](https://github.com/lupinelab/replay/releases/latest).

### From Source
```bash
go install github.com/lupinelab/replay/cmd/replay@latest
```

### Build from source
```bash
git clone https://github.com/lupinelab/replay.git
cd ./replay/cmd
go build replay.go
```

## Usage
### Command Order
By default commands are run in the order they appear in the history. Commands can be given a specific order to run in using `Alt+[1-9]` keys. Any commands with a specified order will be run before commands without, regardless of the number.

### Search
Any non-control key (control keys are listed in the footer) will be passed to the search pane. Search is a simple linear search with no glob or regex support.

### Running A Replay
Run the replay as show in the preview pane using `Alt+Enter`.