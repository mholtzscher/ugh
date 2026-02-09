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

    subgraph NLP["NLP Layer (Participle DSL)"]
        E --> F[Lexer<br/>dslLexer]
        F --> G[Participle Parser<br/>dslParser]
        G --> H[Grammar Nodes<br/>dsl_nodes.go]
        H --> I[Post-Process<br/>dsl_postprocess.go]
        I --> J[AST Commands]
    end

    subgraph Compile["Compile Layer"]
        J --> K[Validate & Normalize]
        K --> L[Build Execution Plan]
    end

    subgraph Execute["Execution Layer"]
        L --> M[Service Calls]
        M --> N[Database]
        M --> O[Output Display]
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
- `display.go` - Output formatting

**Features:**
- **Three modes**: Interactive REPL, script file execution, stdin pipe
- **Pronoun resolution**: `it`/`this` → last task, `that` → second-to-last, `selected` → selected task
- **Sticky context**: `context #project` and `context @context` apply to all subsequent commands
- **Session persistence**: Tracks recently accessed tasks and selected task

### 2. NLP Layer (`internal/nlp/`)

Uses [participle](https://github.com/alecthomas/participle) to parse natural language-like commands:

```mermaid
flowchart TD
    A[Raw Input] --> B[Lexer<br/>dslLexer]
    B --> C[Tokens]
    C --> D[Participle Parser<br/>dslParser]
    D --> E[Grammar Nodes<br/>CreateCommand etc.]
    E --> F[Post-Process<br/>dsl_postprocess.go]
    F --> G[AST Commands]

    subgraph Tokens["Token Types"]
        T1[Quoted: "buy milk"]
        T2[ProjectTag: #groceries]
        T3[ContextTag: @store]
        T4[SetField: title:]<-->T5[Ident: buy]
        T6[AndOp: &&/and]<-->T7[OrOp: ||/or]
    end

    subgraph Grammar["Grammar Nodes"]
        G1[CreateCommand]
        G2[UpdateCommand]
        G3[FilterCommand]
    end

    subgraph AST["AST Command Types"]
        A1[CreateCommand]
        A2[UpdateCommand]
        A3[FilterCommand]
    end

    C --> Tokens
    D --> Grammar
    G --> AST

    style Tokens fill:#fff3e0
    style Grammar fill:#e8f5e9
    style AST fill:#e8f5e9
```

**Key Files:**
- `lexer.go` - Participle lexer with regex rules for tokens
- `dsl_parser.go` - Participle parser configuration with Union types
- `dsl_nodes.go` - Grammar node types with struct tags and custom Parse methods
- `dsl_parse.go` - Custom Parse implementations for verbs, targets, operators
- `dsl_symbols.go` - Token type symbol mapping
- `dsl_postprocess.go` - Normalization and validation of parsed commands
- `parser.go` - Public parser interface
- `ast.go` - Final AST type definitions
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

**Participle Grammar Features:**
- **Struct tags** for simple token matching: `` `parser:"@Ident"``
- **Custom Parse methods** for complex logic (synonyms, multi-token values)
- **Union types** for polymorphic nodes: `participle.Union[Command](&CreateCommand{}, ...)`
- **Custom Capture** methods for field normalization

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
- `plan_test.go` - Compilation tests

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
    participant Post as Post-Process
    participant Compile as Compile Plan
    participant Service as Task Service
    participant DB as SQLite

    User->>Shell: "add buy milk #groceries"
    Shell->>Shell: Preprocess (pronouns, context)
    Shell->>NLP: Parse input
    NLP->>NLP: Lex + Participle parse
    NLP->>Post: Grammar nodes
    Post->>Post: Normalize & validate
    Post->>Compile: AST Command
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
    B --> C["Tokens:<br/>Ident: add<br/>Ident: buy<br/>Ident: milk<br/>ProjectTag: #groceries<br/>ContextTag: @store<br/>SetField: due:<br/>Ident: tomorrow"]
    C --> D["Participle Parser → CreateCommand"]
    D --> E["Post-Process:<br/>Title: 'buy milk'<br/>Ops: [#groceries, @store, due:tomorrow]"]
    E --> F["Compile:<br/>due:tomorrow → 2026-02-09"]
    F --> G["CreateTaskRequest"]
    G --> H["Service.CreateTask"]
    H --> I["New task created<br/>ID: 123"]
```

### Example 2: Updating with Pronouns

```mermaid
flowchart LR
    A["set it state:done"] --> B["Preprocessor:<br/>it → 123"]
    B --> C["Lexer & Participle Parse"]
    C --> D["UpdateCommand:<br/>Target: 123<br/>Ops: [state:done]"]
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

1. **Participle Parser**: Uses `participle` library for maintainable grammar definitions with struct tags and custom Parse methods
2. **Grammar Nodes → AST**: Participle produces grammar-specific nodes that are post-processed into the final AST
3. **Custom Parse Methods**: Used for verb synonyms (add/create/new), target references (#123, selected, it), and multi-token filter values
4. **Struct Tags**: Simple token matching where possible, custom logic only where needed
5. **Stateful Shell**: Session tracks last/selected tasks for natural pronoun usage
6. **Sticky Context**: Project/context filters persist across commands for workflow efficiency
7. **Operations Model**: Consistent `+`/`-`/`!` syntax for list field modifications
8. **Separate Compilation**: AST → Plan → Service request enables validation and normalization

## File Locations

```
internal/
├── shell/
│   ├── repl.go              # Interactive/scripting modes
│   ├── executor.go            # Shell-NLP bridge
│   ├── prompt.go              # Readline integration
│   ├── scripting.go           # File/stdin processing
│   ├── display.go             # Output formatting
│   └── types.go               # Session types
└── nlp/
    ├── lexer.go               # Participle lexer rules
    ├── dsl_parser.go          # Participle parser config
    ├── dsl_nodes.go           # Grammar node types
    ├── dsl_parse.go           # Custom Parse methods
    ├── dsl_symbols.go         # Token type symbols
    ├── dsl_postprocess.go     # Grammar → AST normalization
    ├── parser.go              # Public parser interface
    ├── parser_test.go         # Parser tests
    ├── ast.go                 # Final AST types
    ├── types.go               # Parse types
    └── compile/
        ├── plan.go            # AST → Execution plan
        └── plan_test.go       # Compilation tests
```

## Participle Grammar Details

### Lexer Rules (lexer.go)

Tokens are defined with regex patterns in priority order:
- `Quoted`: `"..."` strings
- `HashNumber`: `#123` numeric IDs
- `ProjectTag`: `#word` project tags
- `ContextTag`: `@word` context tags
- `SetField`: `field:` field setters
- `AddField`: `+field:` field additions
- `RemoveField`: `-field:` field removals
- `ClearField`: `!field` field clearing
- `Ident`: words and identifiers
- `Whitespace`: spaces (elided)

### Union Types (dsl_parser.go)

Participle unions enable polymorphic parsing:

```go
participle.Union[Command](
    &CreateCommand{},
    &UpdateCommand{},
    &FilterCommand{},
)
```

### Custom Parse Methods (dsl_parse.go, dsl_nodes.go)

Used when struct tags aren't sufficient:

```go
func (v *CreateVerb) Parse(lex *lexer.PeekingLexer) error {
    // Handle synonyms: add, create, new
    // Return participle.NextMatch if not a match
}
```

### Post-Processing (dsl_postprocess.go)

After participle parsing:
- Combine text tokens into title
- Normalize operations
- Validate required fields
- Set default targets
