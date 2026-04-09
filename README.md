# videogen

CLI tool for AI video generation via [defapi.org](https://defapi.org).

## Installation

**Pre-built binary** (Linux/macOS):

```sh
# Detect OS and arch, download latest release
OS=$(uname -s | tr '[:upper:]' '[:lower:]') ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') && \
curl -fsSL "https://github.com/jhgundersen/videogen/releases/latest/download/videogen-${OS}-${ARCH}" \
  -o ~/.local/bin/videogen && chmod +x ~/.local/bin/videogen
```

**Using Go:**

```sh
go install github.com/jhgundersen/videogen@latest
```

**From source:**

```sh
make install
```

Installs to `~/.local/bin/videogen`.

## Requirements

Set your API key:

```sh
export DEFAPI_API_KEY=your_key_here
```

## Usage

```sh
videogen <model> [flags] <prompt>
```

### Models

#### `seedance` — ByteDance Seedance 2.0

```sh
videogen seedance "a fox running through a snowy forest"
videogen seedance "aerial ocean waves" --duration 10 --ratio 9:16
```

| Flag | Default | Options |
|------|---------|---------|
| `--duration` | `5` | `5`, `10`, `15` |
| `--ratio` | `16:9` | `16:9`, `9:16`, `1:1`, `4:3`, `3:4`, `21:9` |

#### `grok` — xAI Grok Imagine Video

```sh
videogen grok "timelapse of a city at night"
videogen grok "a surreal dreamscape" --duration 15 --ratio 1:1
```

| Flag | Default | Options |
|------|---------|---------|
| `--duration` | `10` | `10`, `15` |
| `--ratio` | `16:9` | `16:9`, `9:16`, `1:1`, `2:3`, `3:2` |

#### `sora` — OpenAI Sora 2 Stable

```sh
videogen sora "epic cinematic mountain landscape"
videogen sora "slow motion rain on a window" --variant sora-2-hd --ratio 9:16
videogen sora "a full short film scene" --duration 25  # auto-selects sora-2-pro
```

| Flag | Default | Options |
|------|---------|---------|
| `--duration` | `10` | `10`, `15`, `25` (25s requires `sora-2-pro`) |
| `--ratio` | `16:9` | `16:9`, `9:16` |
| `--variant` | `sora-2` | `sora-2`, `sora-2-hd`, `sora-2-pro` |

## Output

Generated videos are downloaded to `~/Downloads/videogen_<task_id>.mp4`. The path is printed as a clickable link in terminals that support OSC 8 hyperlinks (Kitty, Alacritty, WezTerm, GNOME Terminal 3.26+).
