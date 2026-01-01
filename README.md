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

## Commands

-   **`twt user create <username>`**: Creates a new user.
-   **`twt login <username>`**: Logs in as an existing user.
-   **`twt logout`**: Logs out the current user.
-   **`twt post <text>`**: Creates a new post for the currently logged-in user.
-   **`twt profile <username>`**: Views a user's posts.
-   **`twt delete <post_id>`**: Deletes your own post.
-   **`twt show <post_id>`**: Shows a single post with details.
-   **`twt follow <username>`**: Follows another user.
-   **`twt unfollow <username>`**: Unfollows a user.
-   **`twt feed`**: Views your personalized feed (posts from followed users and your own posts).

## Project Structure

-   `cmd/`: Contains the command-line interface commands.
-   `internal/db/`: Database related functionalities.
-   `internal/models/`: Data models for the application (user, post, social).
-   `internal/store/`: Data store interfaces and implementations.
-   `main.go`: The main entry point of the application.
