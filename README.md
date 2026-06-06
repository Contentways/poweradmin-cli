# poweradmin-cli

A command-line interface for managing DNS zones, records and users via the [Poweradmin](https://www.poweradmin.org) REST API.

Built with [poweradmin-go](https://contentways.dev/contentways/poweradmin-go) — the Go SDK for Poweradmin.

## Requirements

- Poweradmin 4.3.0+ running in API-mode (`PA_DNS_BACKEND=api`)
- A valid Poweradmin API key

## Installation

### From release

Download the latest binary for your platform from the [releases page](https://github.com/Contentways/poweradmin-cli/releases):

```bash
# Linux (amd64)
curl -L https://github.com/Contentways/poweradmin-cli/releases/latest/download/poweradmin-cli_Linux_x86_64.tar.gz | tar xz
sudo mv poweradmin /usr/local/bin/
```

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

### Users

```bash
# List all users
poweradmin users list
poweradmin users list --output json

# Get a user by name or ID
poweradmin users get --name max
poweradmin users get --id 1
poweradmin users get --name max --output json

# Create a user (password will be prompted)
poweradmin users create \
  --username max \
  --email max@example.com \
  --fullname "Max Mustermann

# Create a user with password via flag (not recommended — visible in shell history)
poweradmin users create \
  --username max \
  --password secret123 \
  --email max@example.com

# Update a user
poweradmin users update --name max --email new@example.com
poweradmin users update --id 1 --active=false

# Delete a user
poweradmin users delete --name max
poweradmin users delete --id 1

# Assign a permission template
poweradmin users set-permission-template --name max --template-id 2
```

## Output Formats

| Flag | Description |
|------|-------------|
| `--output table` | Human-readable aligned table, long content truncated (default) |
| `--output full` | Human-readable aligned table, full content |
| `--output json` | JSON output, suitable for scripting and piping into `jq` |

## License

MIT — Copyright (c) 2026 Contentways
