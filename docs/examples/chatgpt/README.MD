## 🤖 LowBot Chat GPT Example

### Build

This assumes you already have a working Go environment setup and that
LowBot is correctly installed on your system.

From within the dungeon example folder, run the below command to compile the
example.

```sh
go build main.go
```

### Usage

In this example, the channel used is Telegram, so you need to enter the Telegram Token.
The consumer used is ChatGPT, so you need to enter the ChatGPT Token.

```sh
TELEGRAM_TOKEN=YOUR_TELEGRAM_TOKEN CHATGPT_TOKEN=YOUR_CHATGPT_TOKEN ./main
Bot is now running.  Press CTRL-C to exit.
```