# Datetime Display Architecture

## Overview

Datetime display formatting is configured via `[display]` config and applied in the output layer.
`output.Writer` owns `TimeFormatter`, and shell code receives a prebuilt writer instead of threading config/formatter values.

## Component Relationship Diagram

```mermaid
flowchart LR
  subgraph CMD[CLI Wiring]
    root["cmd/root.go\noutputWriter()"]
    shellCmd["cmd/shell.go"]
  end

  subgraph SH[Shell Runtime]
    repl["internal/shell/repl.go\nOptions.Writer"]
    display["internal/shell/display.go\nuses output.Writer"]
  end

  subgraph OUT[Output Package]
    outCore["internal/output/output.go\nNewWriter + Writer"]
    human["internal/output/human.go\nWriter methods"]
    tf["internal/output/time_format.go\nTimeFormatter"]
  end

  cfg["internal/config/config.go\n[display]"]

  cfg -->|display config| outCore
  root -->|build output.Writer| outCore
  root -->|inject prebuilt writer| shellCmd
  shellCmd -->|opts.Writer| repl
  repl -->|render| display
  display -->|Writer.Write*| human
  outCore -->|Writer owns formatter| tf
  human -->|uses w.formatter| tf
```

## Sequence Diagram (Init + Render)

```mermaid
sequenceDiagram
  participant Config as internal/config/config.go
  participant Root as cmd/root.go
  participant Out as internal/output/output.go
  participant ShellCmd as cmd/shell.go
  participant REPL as internal/shell/repl.go
  participant Display as internal/shell/display.go
  participant Human as internal/output/human.go
  participant TF as internal/output/time_format.go

  Note over Config,Root: Startup / wiring
  Root->>Config: loadConfig()
  Config-->>Root: Config{Display}
  Root->>Out: output.NewWriter(jsonMode, cfg.Display)
  Out->>TF: NewTimeFormatter(cfg.Display)
  TF-->>Out: formatter
  Out-->>Root: Writer{formatter}

  Root->>ShellCmd: invoke shell command
  ShellCmd->>REPL: NewREPL(opts{Writer})
  REPL->>Display: NewDisplay(tty, writer)

  Note over REPL,Human: Runtime rendering
  REPL->>Display: ShowResult(result)
  Display->>Out: writer.WriteTask/WriteTasks/WriteHistory
  Out->>Human: w.writeHuman*()
  Human->>TF: w.formatter.Format(...)
  TF-->>Human: formatted datetime
  Human-->>Out: rendered lines
```
