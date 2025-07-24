# v1.0.0

## Changelog

- Initial release of Alchemist CLI.
- Add main `alchemist` command with subcommands.
- Add `make` subcommand for scaffolding files and components.
- Add `migrate` subcommand to run and manage database migrations.
- Support configuration via `.env` file and environment variables using Viper.
- Add colored output for success, error, info, and warning messages.
- Add `--config` flag to specify configuration file.
- Add `--dir` flag to set migrations directory.
- Show version with `--version`.
- Organize project structure for maintainability and extensibility.
- Add `.gitignore` for build artifacts, dependencies, IDE files, databases, and configs.
- Example usage of Go `embed` package for future assets.
- Standardize error handling and output messages.