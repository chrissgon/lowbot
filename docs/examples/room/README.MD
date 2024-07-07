## 🤖 LowBot Ping Pong Example

This example demonstrates how to utilize LowBot to create a Ping Pong Bot.

This Bot will respond to "ping" with "Pong!" and "pong" with "Ping!".

### Build

This assumes you already have a working Go environment setup and that
LowBot is correctly installed on your system.

From within the pingpong example folder, run the below command to compile the
example.

```sh
go build main.go
```

### Usage

In this example, the channel used is Telegram, so you need to enter the Telegram Token

```sh
TELEGRAM_TOKEN=YOUR_TELEGRAM_TOKEN ./main
Bot is now running.  Press CTRL-C to exit.
```