# lowbot

🤖 LowBot is a Go project on GitHub that simplifies bot creation by automatically generating bot journeys from a YAML script. It offers extensive customization options, making it highly versatile and adaptable.

## 🤔 Why LowBot?

- ### 📚 Simple and Self-Explanatory Scripting:

  LowBot follows a straightforward and intuitive YAML script format, making it easy to understand and modify for creating bot journeys.

- ### 🧩 Extensibility for Custom Business Rules:

  LowBot offers extensive flexibility, allowing users to define their own custom actions and implement unique business rules tailored to their specific requirements.

- ### 📡 Expandable Bot Channels:
  With LowBot, users can seamlessly integrate and implement their own channels, empowering them to connect the bot with various platforms and communication channels of their choice.

## 📦 Install

```bash
go get github.com/chrissgon/lowbot
```

## 🚀 Quick Start

Start a bot with Telegram and local persist

```golang
import "github.com/chrissgon/lowbot"

flow, _ := lowbot.NewFlow("your_flow.yaml")
channel, _ := lowbot.NewTelegram("your_telegram_token")
persist, _ := lowbot.NewLocalPersist()

lowbot.StartBot(flow, channel, persist)
```

Start a bot with Discord and custom actions

```golang
import "github.com/chrissgon/lowbot"

myActions := lowbot.ActionsMap{
    "Custom": func(flow *lowbot.Flow, channel lowbot.IChannel) (bool, error) {
        // your rules
        wait := true
        return wait, nil
    },
}

lowbot.SetCustomActions(myActions)

flow, _ := lowbot.NewFlow("your_flow.yaml")
channel, _ := lowbot.NewDiscord()
persist, _ := lowbot.NewLocalPersist()

lowbot.StartBot(flow, channel, persist)
```

You can create your Channel or Persist by implementing their interfaces

```golang
import "github.com/chrissgon/lowbot"

func MyNewChannel () lowbot.IChannel {
  // implements Channel
}

func MyNewPersist () lowbot.Persist {
  // implements Persist
}

flow, _ := lowbot.NewFlow("your_flow.yaml")

// And you can pass in the StartBot
lowbot.StartBot(flow, MyNewChannel(), MyNewPersist())
```

## 📚 Documentation

Read all documentation in [docs](https://github.com/chrissgon/lowbot/tree/main/docs) folder.

## 🌎 Global Variables

When you run this project, you can adjusts the following environment variables.

`Debug`: Show debug => `default: true`

`EnableLocalPersist`: Enable local persist => `default: true`

## 🔒 Environment Variables

When you run this project, you need set the following environment variables.

- ### When you to use Telegram:
  `TELEGRAM_TOKEN`: Telegram bot token.
- ### When you to use Discord:
  `DISCORD_TOKEN`: Discord bot token.

## 📝 Anotations

Run all the tests with coverage.
```bash
go test ./... -coverprofile=coverage.out && ./coverage-ignore.sh && go tool cover -html=coverage.out
```

## 💪🏻 Contribution

This project is open source and welcomes community contributions. Feel free to fork, implement improvements, and submit a pull request. Every contribution is valued and appreciated!

We hope that lowbot proves useful to you and enhances your manga reading experience. Feel free to explore the source code, provide feedback, and report any issues you encounter.

## ❤️ Authors

- [@chrissgon](https://www.github.com/chrissgon)
