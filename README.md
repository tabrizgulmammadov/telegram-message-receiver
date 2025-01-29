
# Telegram Message Receiver

This project is a sample application for receiving messages from Telegram. It demonstrates the basic setup for interacting with the Telegram API and processing received messages.

## Project Structure

- `config/config.go`: Contains configuration settings for the application.
- `handler/handler.go`: Handles incoming messages and related logic.
- `logger/logger.go`: Implements logging functionality for the application.
- `storage/storage.go`: Manages storage and retrieval of data.
- `.env`: Environment variables for configuration (e.g., Telegram bot token).
- `.gitignore`: Specifies files to be ignored by Git.
- `Dockerfile`: Docker configuration to build the container image.
- `docker-compose.yaml`: Docker Compose configuration for multi-container setup.
- `go.mod`, `go.sum`: Go module files for dependency management.
- `main.go`: The entry point of the application.
- `README.md`: This file.

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/telegram-message-receiver.git
   cd telegram-message-receiver
   ```

2. Create a `.env` file and add your Telegram bot token:
   ```
   TELEGRAM_BOT_TOKEN=your-telegram-bot-token
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

Alternatively, you can use Docker to run the application:
```bash
docker-compose up --build
```

## License

This project is licensed under the MIT License.
