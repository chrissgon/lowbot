## 🤖 LowBot Hello Example

This example demonstrates how to utilize LowBot with webhook action.

This bot will ask `👋 Hi. What is your name?` and will answer `Nice to meet you <USERNAME>` with the webhook

### Build

This assumes you already have a working Go environment setup and that
LowBot is correctly installed on your system.

From within the hello example folder, run the below command to compile the
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