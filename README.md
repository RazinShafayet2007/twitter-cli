# Twitter CLI

This is a command-line interface (CLI) application for a simplified Twitter-like service.

## Getting Started

To build and run the application, follow these steps:

1.  **Build the executable:**

    ```bash
    go build -o twt .
    ```

    This will create an executable named `twt` in the current directory.

2.  **Create a user:**

    ```bash
    ./twt user create <username>
    ```

    Replace `<username>` with your desired username (e.g., `alice`).

## Project Structure

-   `cmd/`: Contains the command-line interface commands.
-   `internal/db/`: Database related functionalities.
-   `internal/models/`: Data models for the application (user, post, social).
-   `internal/store/`: Data store interfaces and implementations.
-   `main.go`: The main entry point of the application.

