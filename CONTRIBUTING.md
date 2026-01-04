# Contributing to Twitter CLI

Thank you for your interest in contributing to `twitter-cli`! This project is a learning exercise in building backend systems with Go, and we welcome improvements, bug fixes, and new features.

## Getting Started

1.  **Fork the repository** on GitHub.
2.  **Clone your fork** locally:
    ```bash
    git clone https://github.com/RazinShafayet2007/twitter-cli.git
    cd twitter-cli
    ```
3.  **Install Go**: Ensure you have Go 1.21 or later installed.

## Development

### Running Locally
To run the CLI from source:
```bash
go run main.go [command]
```

### Running Tests
Run all tests to ensure your changes didn't break anything:
```bash
go test ./...
```

### Code Style
We use standard Go formatting. Before committing, please run:
```bash
go fmt ./...
go vet ./...
```

## Release Workflow & Changesets

This project uses [Changesets](https://github.com/changesets/changesets) so that we can automate versioning and changelogs. **This is mandatory for all code changes.**

### How to add a Changeset

1.  Make your code changes.
2.  Run the changeset wizard:
    ```bash
    npx changeset
    ```
    *(Note: You need Node.js installed to run this. If you don't have it, please install it or ask a maintainer for help).*
3.  Select the type of change:
    -   **patch**: Bug fixes (0.0.x)
    -   **minor**: New features (0.x.0)
    -   **major**: Breaking changes (x.0.0)
4.  Write a brief summary of your change.
5.  This will create a new file in `.changeset/` (e.g., `.changeset/warm-clouds-sing.md`).
6.  **Commit this file** along with your code.

### Why do I need this?
If you don't include a changeset file, our CI checks will fail on your Pull Request. This system allows us to automatically update the version number and changelog when your code is merged.

## Pull Requests

1.  Create a new branch for your feature or fix.
2.  Push your branch to your fork.
3.  Open a Pull Request against the `main` branch.
4.  Ensure the `Require Changeset` CI check passes.
