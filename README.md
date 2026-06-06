# poweradmin-cli

A command-line interface for managing DNS zones and records via the [Poweradmin](https://www.poweradmin.org) REST API.

Built with [poweradmin-go](https://contentways.dev/contentways/poweradmin-go) — the official Go SDK for Poweradmin.

## Requirements

- Poweradmin 4.3.0+ running in API-mode (`PA_DNS_BACKEND=api`)
- A valid Poweradmin API key

## Installation

### From source

```bash
git clone git@git.contentways.dev:contentways/poweradmin-cli.git
cd poweradmin-cli
make build
```

The binary is placed at `build/poweradmin`.

## Configuration

Credentials are resolved in the following order of precedence:

| Source | Example |
|--------|---------|
| CLI flags | `--url`, `--api-key` |
| Environment variables | `POWERADMIN_URL`, `POWERADMIN_API_KEY` |
| Config file | `~/.config/poweradmin/config.yaml` |

### Config file

```yaml
url: https://dns.example.com
api_key: pwa_your_api_key_here
```

Create the config directory and file:

```bash
mkdir -p ~/.config/poweradmin
cat > ~/.config/poweradmin/config.yaml << 'EOF'
url: https://dns.example.com
api_key: pwa_your_api_key_here
EOF
```

### Environment variables

```bash
export POWERADMIN_URL=https://dns.example.com
export POWERADMIN_API_KEY=pwa_your_api_key_here
```

### CLI flags

```bash
poweradmin zones list --url https://dns.example.com --api-key pwa_...
```

## Shell Completion

### zsh (Oh My Zsh)

```bash
poweradmin completion zsh > ~/.oh-my-zsh/completions/_poweradmin
source ~/.zshrc
```

### bash

```bash
poweradmin completion bash > ~/.bash_completion.d/poweradmin
source ~/.bash_completion.d/poweradmin
```

## Usage

### Zones

```bash
# List all zones
poweradmin zones list
poweradmin zones list --output json

# Get a zone by name or ID
poweradmin zones get --name example.com
poweradmin zones get --id 42
poweradmin zones get --name example.com --output json

# Create a zone
poweradmin zones create example.com
poweradmin zones create example.com --type NATIVE
poweradmin zones create example.com --output json

# Delete a zone
poweradmin zones delete --name example.com
poweradmin zones delete --id 42
```

### Records

```bash
# List all records in a zone
poweradmin records list --zone-name example.com
poweradmin records list --zone-id 42
poweradmin records list --zone-name example.com --output full
poweradmin records list --zone-name example.com --output json

# Create a record
poweradmin records create \
  --zone-name example.com \
  --name www.example.com \
  --type A \
  --content 1.2.3.4 \
  --ttl 3600

# Create a record and capture the ID
poweradmin records create \
  --zone-name example.com \
  --name www.example.com \
  --type A \
  --content 1.2.3.4 \
  --output json | jq '.id'

# Delete a record
poweradmin records delete \
  --zone-name example.com \
  --id eyJ6IjoiZXhhbXBsZS5jb20i...
```

## Output Formats

| Flag | Description |
|------|-------------|
| `--output table` | Human-readable aligned table, long content truncated (default) |
| `--output full` | Human-readable aligned table, full content |
| `--output json` | JSON output, suitable for scripting and piping into `jq` |

## License

MIT — Copyright (c) 2026 Contentways
