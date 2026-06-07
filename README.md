# poweradmin-cli

A command-line interface for managing DNS zones, records, users and groups via the [Poweradmin](https://www.poweradmin.org) REST API.

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
| CLI flags | `-u`, `-k` |
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
poweradmin zones list -u https://dns.example.com -k pwa_...
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
poweradmin zones list -o json

# Filter zones
poweradmin zones list --type NATIVE
poweradmin zones list --name-filter contentways
poweradmin zones list --type NATIVE --name-filter contentways

# Get a zone by name or ID
poweradmin zones get --name example.com
poweradmin zones get --id 42
poweradmin zones get --name example.com -o json

# Create a zone
poweradmin zones create example.com
poweradmin zones create example.com --type NATIVE
poweradmin zones create example.com \
  --nameserver ns1.example.com \
  --nameserver ns2.example.com

# Delete a zone
poweradmin zones delete --name example.com
poweradmin zones delete --id 42
```

### Records

```bash
# List all records in a zone
poweradmin records list --zone-name example.com
poweradmin records list --zone-name example.com -o full
poweradmin records list --zone-name example.com -o json

# Create a record
poweradmin records create \
  --zone-name example.com \
  --name www.example.com \
  --type A \
  --content 1.2.3.4 \
  --ttl 3600

# Update a record
poweradmin records update \
  --zone-name example.com \
  --id eyJ6IjoiZXhhbXBsZS5jb20i... \
  --content 5.6.7.8

# Delete a record
poweradmin records delete \
  --zone-name example.com \
  --id eyJ6IjoiZXhhbXBsZS5jb20i...
```

### Users

```bash
# List all users
poweradmin users list
poweradmin users list -o json

# Get a user by name or ID
poweradmin users get --name patrick
poweradmin users get --id 1

# Create a user (password will be prompted)
poweradmin users create \
  --username patrick \
  --email patrick@example.com \
  --fullname "Patrick Omland"

# Create a user with password via flag (not recommended — visible in shell history)
poweradmin users create \
  --username patrick \
  --password secret123 \
  --email patrick@example.com

# Update a user
poweradmin users update --name patrick --email new@example.com
poweradmin users update --id 1 --active=false

# Delete a user
poweradmin users delete --name patrick
poweradmin users delete --id 1

# Assign a permission template
poweradmin users set-permission-template --name patrick --template-id 2
```

### Groups

```bash
# List all groups
poweradmin groups list
poweradmin groups list -o json

# Get a group by name or ID
poweradmin groups get --name Administrators
poweradmin groups get --id 1

# Create a group
poweradmin groups create --name "Zone Editors" --description "Can edit zone records"

# Update a group
poweradmin groups update --name "Zone Editors" --new-name "DNS Editors"

# Delete a group
poweradmin groups delete --name "DNS Editors"

# Manage members
poweradmin groups members --name Administrators
poweradmin groups member-add --group-id 1 --user-id 2
poweradmin groups member-remove --group-id 1 --user-id 2

# Manage zones
poweradmin groups zones --name Administrators
poweradmin groups zone-add --group-id 1 --zone-id 78
poweradmin groups zone-remove --group-id 1 --zone-id 78
```

### Version

```bash
poweradmin version
```

## Output Formats

| Flag | Description |
|------|-------------|
| `-o table` | Human-readable aligned table, long content truncated (default) |
| `-o full` | Human-readable aligned table, full content |
| `-o json` | JSON output, suitable for scripting and piping into `jq` |

## Documentation

Full reference documentation is available in [docs/reference](docs/reference/).

## License

MIT — Copyright (c) 2026 Contentways
