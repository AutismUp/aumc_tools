# AUMC Tools - Development Guide

## Installation & Setup
```bash
pip install -e .                              # Install in development mode
aumc                                          # Run CLI (requires AU_CONFIG_FILE env var)
```

## Testing
No automated test suite currently exists. Manual testing via Vagrant VM:
```bash
vagrant up                                    # Start test server
vagrant ssh                                   # SSH into test server
```

## Code Style

### Imports
- Standard library imports first, then third-party (Click), then local modules
- Absolute imports preferred: `from aumc import aumc`

### Formatting & Structure
- 4-space indentation
- Classes use PascalCase: `AuMc`, `MCConfig`, `EnterDir`
- Functions/methods use snake_case: `create_new_world`, `build_new_jar`
- Private methods: no leading underscore convention used
- String quotes: single quotes preferred

### Error Handling
- Custom exceptions for domain errors: `AuServerCreationException`, `AuRestoreException`
- Use subprocess.call() for non-critical commands, subprocess.run() with check=True for critical ones
- Print error messages to stdout for user feedback

### CLI Commands (Click framework)
- All CLI commands decorated with @cli.command()
- Use click.echo() for output, not print() in CLI commands
- Options use --kebab-case format
