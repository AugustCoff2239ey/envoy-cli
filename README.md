# envoy-cli

> A CLI tool for managing and syncing environment variable sets across local, staging, and production configs.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envoy-cli/releases).

---

## Usage

```bash
# Initialize a new envoy config in the current directory
envoy init

# Add an environment variable to a specific environment
envoy set DATABASE_URL="postgres://localhost:5432/mydb" --env local

# Sync variables from local to staging
envoy sync --from local --to staging

# List all variables for an environment
envoy list --env production

# Pull the latest config from a remote source
envoy pull --env production
```

Configs are stored in `.envoy/` at the root of your project. Each environment (e.g., `local`, `staging`, `production`) is tracked in its own file and can be version-controlled or kept private via `.gitignore`.

---

## Configuration

By default, `envoy-cli` looks for a `.envoy/config.yaml` file. You can specify a custom config path using the `--config` flag:

```bash
envoy list --config /path/to/config.yaml --env staging
```

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes and open a PR

---

## License

This project is licensed under the [MIT License](LICENSE).