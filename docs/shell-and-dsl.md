# Shell & DSL Architecture

## Overview

`ugh` provides an interactive shell for task management using a natural language-like Domain Specific Language (DSL). This document explains how the shell and DSL work together to parse and execute commands.

## Architecture Diagram

```mermaid
flowchart TD
    subgraph Shell["Shell Layer (Interactive)"]
        A[User Input] --> B[REPL/Script Mode]
        B --> C[Preprocessing]
        C --> D[Pronoun Resolution]
        D --> E[Context Injection]
    end

    subgraph NLP["NLP Layer (Parsing)"]
        E --> F[Lexer]
        F --> G[Grammar Parser]
        G --> H[Convert to AST]
    end

    subgraph Compile["Compile Layer"]
        H --> I[Validate & Normalize]
        I --> J[Build Execution Plan]
    end

    subgraph Execute["Execution Layer"]
        J --> K[Service Calls]
        K --> L[Database]
        K --> M[Output Display]
    end

    style Shell fill:#e1f5fe
    style NLP fill:#fff3e0
    style Compile fill:#e8f5e9
    style Execute fill:#fce4ec
```

## Component Breakdown

### 1. Shell Layer (`internal/shell/`)

The shell provides the user interface and session management:

```mermaid
flowchart LR
    subgraph Modes["Operating Modes"]
        A1[Interactive REPL] --> A2[Script File]
        A1 --> A3[Stdin Pipe]
    end

    subgraph State["Session State"]
        B1[Selected Task ID] --> B4[Executor]
        B2[Last Task IDs] --> B4
        B3[Context Filters<br/>#project @context] --> B4
    end

    A1 --> B4
    A2 --> B4
    A3 --> B4

    style Modes fill:#e1f5fe
    style State fill:#fff3e0
```

**Key Files:**
- `repl.go` - Main REPL loop (interactive/scripting modes)
- `executor.go` - Bridges shell to NLP and service layers
- `prompt.go` - Interactive prompt with readline (history, editing)
- `scripting.go` - Script file processing
- `types.go` - Session state and execution types

**Features:**
- **Three modes**: Interactive REPL, script file execution, stdin pipe
- **Pronoun resolution**: `it`/`this` → last task, `that` → second-to-last, `selected` → selected task
- **Sticky context**: `context #project` and `context @context` apply to all subsequent commands
- **Session persistence**: Tracks recently accessed tasks and selected task

### 2. NLP Layer (`internal/nlp/`)

Parses natural language-like commands into executable structures:

```mermaid
flowchart TD
    A[Raw Input] --> B[Lexer<br/>participle]
    B --> C[Tokens]
    C --> D[Grammar Parser<br/>GInput struct]
    D --> E[Convert<br/>convert.go]
    E --> F[AST Commands]

    subgraph Tokens["Token Types"]
        T1[Quoted: "buy milk"]
        T2[ProjectTag: #groceries]
        T3[ContextTag: @store]
        T4[Verb: add/create/set/find]
        T5[Op: + - !]
        T6[Logic: and/or/not]
    end

    subgraph AST["AST Command Types"]
        A1[CreateCommand]
        A2[UpdateCommand]
        A3[FilterCommand]
    end

    C --> Tokens
    F --> AST

    style Tokens fill:#fff3e0
    style AST fill:#e8f5e9
```

**Key Files:**
- `lexer.go` - Tokenizer using participle library
- `grammar.go` - Grammar structures (auto-generated)
- `convert.go` - Converts grammar tokens to typed AST (706 lines)
- `ast.go` - AST type definitions
- `types.go` - Parse modes and intent types

**Operations Supported:**

```mermaid
flowchart LR
    subgraph Ops["Operation Types"]
        S[SetOp<br/>field:value]
        A[AddOp<br/>+field:value]
        R[RemoveOp<br/>-field:value]
        C[ClearOp<br/>!field]
        T[TagOp<br/>#project @context]
    end

    style Ops fill:#e8f5e9
```

**Fields:**
- `title`, `notes`, `due`, `waiting`/`waiting-for`, `state`
- `projects`, `contexts`, `meta` (list fields supporting Add/Remove)

### 3. Compile Layer (`internal/nlp/compile/`)

Transforms AST commands into service-ready execution plans:

```mermaid
flowchart LR
    subgraph Input["AST Input"]
        A1[CreateCommand]
        A2[UpdateCommand]
        A3[FilterCommand]
    end

    subgraph Compile["Compilation"]
        B[Normalize Dates<br/>today, tomorrow, next-week]
        C[Normalize States<br/>todo → inbox]
        D[Resolve Targets<br/>selected → ID]
        E[Build Plan]
    end

    subgraph Output["Service Requests"]
        F1[CreateTaskRequest]
        F2[UpdateTaskRequest]
        F3[ListTasksRequest]
    end

    A1 --> Compile --> F1
    A2 --> Compile --> F2
    A3 --> Compile --> F3

    style Input fill:#fff3e0
    style Compile fill:#e8f5e9
    style Output fill:#fce4ec
```

**Key File:**
- `plan.go` - Main compilation logic (464 lines)

**Responsibilities:**
- Date normalization: `today`, `tomorrow`, `next-week` → `YYYY-MM-DD`
- State normalization: `todo` → `inbox`
- Target resolution: `selected`, `it`, `that` → actual task IDs
- Filter compilation: Build SQL-compatible filter expressions

### 4. Execution Layer

The executor coordinates between shell state and service layer:

```mermaid
sequenceDiagram
    participant User
    participant Shell as Shell (REPL)
    participant NLP as NLP Parser
    participant Compile as Compile Plan
    participant Service as Task Service
    participant DB as SQLite

    User->>Shell: "add buy milk #groceries"
    Shell->>Shell: Preprocess (pronouns, context)
    Shell->>NLP: Parse input
    NLP->>NLP: Lex + Grammar parse
    NLP->>NLP: Convert to AST
    NLP->>Compile: AST Command
    Compile->>Compile: Normalize & validate
    Compile->>Shell: Execution Plan
    Shell->>Service: Execute plan
    Service->>DB: CRUD operations
    DB->>Service: Results
    Service->>Shell: Task result
    Shell->>Shell: Update session state
    Shell->>User: Display output
```

**Key File:**
- `internal/shell/executor.go`

## Data Flow Examples

### Example 1: Creating a Task

```mermaid
flowchart LR
    A["add buy milk #groceries @store due:tomorrow"] --> B["Lexer"]
    B --> C["Tokens:<br/>Verb: add<br/>Text: buy milk<br/>Tag: #groceries<br/>Tag: @store<br/>Set: due:tomorrow"]
    C --> D["Grammar → GInput"]
    D --> E["Convert → CreateCommand"]
    E --> F["Compile:<br/>due:tomorrow → 2026-02-09"]
    F --> G["CreateTaskRequest"]
    G --> H["Service.CreateTask"]
    H --> I["New task created<br/>ID: 123"]
```

### Example 2: Updating with Pronouns

```mermaid
flowchart LR
    A["set it state:done"] --> B["Preprocessor:<br/>it → 123"]
    B --> C["Lexer & Parse"]
    C --> D["UpdateCommand:<br/>Target: 123<br/>SetOp: state:done"]
    D --> E["Compile & Execute"]
    E --> F["Task 123 marked done"]
```

### Example 3: Filter with Context

```mermaid
flowchart LR
    subgraph Context["Sticky Context Set"]
        C1["context #work"]
        C2["context @urgent"]
    end

    A["list state:now"] --> B["Context Injection:<br/>#work + @urgent"]
    B --> C["Effective Query:<br/>state:now and project:work and context:urgent"]
    C --> D["Compile to FilterExpr"]
    D --> E["ListTasksRequest"]
    E --> F["Filtered results"]
```

## Command Syntax Reference

### Creating Tasks

```
add buy milk
add "buy milk and eggs" #groceries @store
create task due:tomorrow state:inbox +projects:personal
new "complex task" #work @urgent waiting-for:bob
```

### Updating Tasks

```
set 123 state:done
set selected title:"new title" notes:"details here"
set it +projects:work -projects:personal
edit that due:next-week
update 456 !notes          # clear notes
```

### Filtering/Querying

```
find state:now
show project:work and state:inbox
list due:tomorrow
filter @urgent
find state:now or state:waiting
show id:123
```

### Context Commands

```
context #project    # Set default project filter
context @context    # Set default context filter
context clear       # Remove all sticky filters
```

## Key Design Decisions

1. **Two-Phase Parsing**: Grammar parsing → AST conversion keeps the grammar simple while enabling rich semantics
2. **Stateful Shell**: Session tracks last/selected tasks for natural pronoun usage
3. **Sticky Context**: Project/context filters persist across commands for workflow efficiency
4. **Operations Model**: Consistent `+`/`-`/`!` syntax for list field modifications
5. **Participle Parser**: Uses `participle` library for maintainable grammar definitions
6. **Separate Compilation**: AST → Plan → Service request enables validation and normalization

## File Locations

```
internal/
├── shell/
│   ├── repl.go           # Interactive/scripting modes
│   ├── executor.go         # Shell-NLP bridge
│   ├── prompt.go           # Readline integration
│   ├── scripting.go        # File/stdin processing
│   ├── display.go          # Output formatting
│   └── types.go            # Session types
└── nlp/
    ├── lexer.go            # Tokenizer
    ├── grammar.go          # Grammar structures
    ├── parser.go           # Parser interface
    ├── convert.go          # Grammar → AST (706 lines)
    ├── ast.go              # AST types
    ├── types.go            # Parse types
    └── compile/
        └── plan.go         # AST → Execution plan (464 lines)
```
