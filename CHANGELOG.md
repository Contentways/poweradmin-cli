# Changelog

## [v1.1.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v1.1.0)

### Features

- add delete confirmation prompt with --yes/-y flag
- add --no-header flag to list commands, --quiet/-q to create/delete
- add color output for TTY terminals
- add records get command with client-side ID lookup
- add sort flags, records get, records/users filter, stderr error output, base package refactor
- add shell completion via API for zone, user and group names

## [v0.4.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v0.4.0)

### Features

- add version command with build-time version injection
- add records update command with tests
- add --type and --name-filter flags to zones list
- add record update, version command, zone filters, fix tests

## [v0.3.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v0.3.0)

### Features

- add groups commands (list, get, create, update, delete, members, zones)

### Bug Fixes

- update tests for new schema package JSON output

## [v0.2.0](https://github.com/Contentways/poweradmin-cli/releases/tag/v0.2.0)

### Features

- add short flags -u (--url), -k (--api-key), -o (--output)
- add schema package for clean JSON output, add root objects and omit empty fields
