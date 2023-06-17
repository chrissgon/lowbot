# Documentation

- [Understanding the Flows 🛣️](#understanding-the-flows-🛣️)
  - [Actions and Parameters](#actions-and-parameters)
  - [Next](#next)
  - [Custom Actions](#actions)
- [Implementing the Channels 🚛]()
  - [Telegram]()
  - [Discord]()
- [Implementing the Persists 🗃️]()
- [Examples](https://github.com/chrissgon/lowbot/tree/main/docs/examples)

# Understanding the Flows 🛣️

The `flow` defines the user's path through the bot. 🚶🏻

It has `steps` that specify the details of each step.

Each step has mandatory values, and they are:

- `action`: Defines the action taken;
- `parameters`: Additional information about `action`;
- `next`: Informs the next step.

The flow must be written in a YAML file, here is an example: 👇🏻

```yaml
name: flow

steps:
  mystep:
    # Send a text
    action: Text
    parameters:
      # Defines the text that will be send to the action Text
      texts:
        - Some text
    # Informs the next step, the value "default" is required
    next:
      default: nextstep

  nextstep: ...
```

Now let's understand about each of the values mentioned. 🤩

## Actions and Parameters

The `action` are functions that contain business rules.

LowBot implements standard `actions` with common rules.

In addition to the business rule, each action needs to inform whether it should continue to the next step, or wait for an interaction from the user.

Each `action` has additional data that is entered in `parameters`.

| Action                               | Description            | Parameter  | Type       | Next ? |
| ------------------------------------ | ---------------------- | ---------- | ---------- | ------ |
| [Text](#example-action-text)         | Send texts             | `texts`    | `[]string` | ✅     |
| [Image](#example-action-image)       | Send image             | `image`    | `string`   | ✅     |
| [Video](#example-action-video)       | Send video             | `video`    | `string`   | ✅     |
| [Document](#example-action-document) | Send document          | `document` | `string`   | ✅     |
| [Input](#example-action-input)       | Send texts             | `input`    | `string`   | ❌     |
| [Button](#example-action-button)     | Send texts and buttons | `buttons`  | `[]string` | ❌     |
| [Wait](#example-action-wait)         | Only wait              | -          | -          | ❌     |

### Example Action Text

```yaml
action: Text
parameters:
  texts:
    - "Hello"
    - "Welcome to LowBot"
```

### Example Action Image

```yaml
action: Image
parameters:
  # Image supports local path or url
  image: image.png | https://myimage.com/image.png
```

### Example Action Video

```yaml
action: Video
parameters:
  # Video supports local path or url
  video: video.mp4 | https://myvideo.com/video.mp4
```

### Example Action Document

```yaml
action: Document
parameters:
  # Document supports local path or url
  document: document.txt | https://mydoc.com/document.txt
```

### Example Action Input

```yaml
action: Input
parameters:
  texts:
    - "Hello"
    - "Are you fine?"
```

### Example Action Button

```yaml
action: Button
parameters:
  texts:
    - "Hello"
    - "Are you fine?"
  buttons:
    - Yes
    - No
```

### Example Action Wait

```yaml
action: Wait
```

## Parameters
