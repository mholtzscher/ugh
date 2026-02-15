# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.1](https://github.com/mholtzscher/ugh/compare/v0.2.0...v0.2.1) (2026-02-14)


### Features

* add task events audit trail ([#99](https://github.com/mholtzscher/ugh/issues/99)) ([a336e56](https://github.com/mholtzscher/ugh/commit/a336e56f4c5ad1a9c31291713cae8c803eb8018e))
* **display:** add configurable datetime formatting ([#109](https://github.com/mholtzscher/ugh/issues/109)) ([63ed6d7](https://github.com/mholtzscher/ugh/commit/63ed6d7a5cbaaca7fecaa4b6bd5ebc1706dc8907))
* **filter:** support wildcard set predicates ([#87](https://github.com/mholtzscher/ugh/issues/87)) ([42ac354](https://github.com/mholtzscher/ugh/commit/42ac3546a40b407d34fd5b8cf6256ecd4e2e8bf6))
* **list:** add recent mode and list limits ([#106](https://github.com/mholtzscher/ugh/issues/106)) ([ef8d30c](https://github.com/mholtzscher/ugh/commit/ef8d30c9971b7c79523edc93742190fe4cecedf6))
* **log:** add task activity command and version diff output ([#103](https://github.com/mholtzscher/ugh/issues/103)) ([a72d27e](https://github.com/mholtzscher/ugh/commit/a72d27e30ec24870d8a47665132dc60dbb64ecf9))
* **seed:** add hidden temp-db seeding command ([#107](https://github.com/mholtzscher/ugh/issues/107)) ([b9ca562](https://github.com/mholtzscher/ugh/commit/b9ca5621ca129aa2684be2943642ae93e79f5654))


### Bug Fixes

* **release:** enable homebrew release ([#98](https://github.com/mholtzscher/ugh/issues/98)) ([691bb6e](https://github.com/mholtzscher/ugh/commit/691bb6e566e468992f3ca98019e4c24f7bf5b7c3))
* **shell:** make calendar view consistent with calendar command ([#85](https://github.com/mholtzscher/ugh/issues/85)) ([6f78f63](https://github.com/mholtzscher/ugh/commit/6f78f6378f6f0a6a25c1e73a81db1382ea061f86))

## [0.2.0](https://github.com/mholtzscher/ugh/compare/v0.1.1...v0.2.0) (2026-02-11)


### ⚠ BREAKING CHANGES

* **tui:** The interactive TUI has been removed. Use CLI commands directly.
* The legacy database schema and todo.txt file format support have been removed. This update requires a clean database migration that resets existing tasks to the new GTD-first structure.

### Features

* add interactive shell with REPL, command history and scripting support ([#16](https://github.com/mholtzscher/ugh/issues/16)) ([854cde7](https://github.com/mholtzscher/ugh/commit/854cde79e46a188572f85a4cc8a179aac8cea56e))
* **config:** auto-initialize config file when missing ([fffaf52](https://github.com/mholtzscher/ugh/commit/fffaf5299df9aba301c94945c2e99c292c90f4ef))
* **history:** add shell command history tracking and viewing ([#44](https://github.com/mholtzscher/ugh/issues/44)) ([56c3edf](https://github.com/mholtzscher/ugh/commit/56c3edf2ebb6dec32bf729fe62129dd62452660f))
* migrate terminal output to pterm with styled messages and inline form editing ([#14](https://github.com/mholtzscher/ugh/issues/14)) ([c186c78](https://github.com/mholtzscher/ugh/commit/c186c78dd5f0820c4c052769c3cba22c8a32012e))
* **nlp:** add due date filtering to task queries ([#47](https://github.com/mholtzscher/ugh/issues/47)) ([4c6426a](https://github.com/mholtzscher/ugh/commit/4c6426aba20aec4a6f844de157980c49c4a887ae))
* **nlp:** add due date filtering to task queries ([#48](https://github.com/mholtzscher/ugh/issues/48)) ([5dff5b0](https://github.com/mholtzscher/ugh/commit/5dff5b015c2ba3dde36b584204f464cb45f44973))
* **nlp:** add view and context command support ([#63](https://github.com/mholtzscher/ugh/issues/63)) ([f433b35](https://github.com/mholtzscher/ugh/commit/f433b353c72217ef1f82d22fe38d4238ee5812dd))
* **nlp:** support natural language due date parsing ([#82](https://github.com/mholtzscher/ugh/issues/82)) ([952a759](https://github.com/mholtzscher/ugh/commit/952a759b09755c1a7d53c0ac5d7ffbf6e8a77d52))
* redesign task model for GTD-first architecture ([#2](https://github.com/mholtzscher/ugh/issues/2)) ([5b7d79a](https://github.com/mholtzscher/ugh/commit/5b7d79a934e8d06e83d66610d77e7eddc18a5e01))
* **repl:** show current context when no args provided ([6145c24](https://github.com/mholtzscher/ugh/commit/6145c2465a4e9d8080ff2ab7e0c73475b42d6433))
* **shell:** add "last" keyword substitution for task references ([50c5b94](https://github.com/mholtzscher/ugh/commit/50c5b94dfe523112df84af5a4e6ea96095698547))
* **shell:** add color output support with pterm integration ([#45](https://github.com/mholtzscher/ugh/issues/45)) ([8e2a5bb](https://github.com/mholtzscher/ugh/commit/8e2a5bb25e05484ba64349ea801553e3b40f60d6))
* **shell:** add prompt autocomplete and syntax highlighting ([#61](https://github.com/mholtzscher/ugh/issues/61)) ([8309359](https://github.com/mholtzscher/ugh/commit/8309359c789da5b11f11c9f9aa8bfe9119a4f2ef))
* **shell:** add quick view shortcuts for common task filters ([#46](https://github.com/mholtzscher/ugh/issues/46)) ([1a4ac29](https://github.com/mholtzscher/ugh/commit/1a4ac291bd759794b45a2978fa8309020845e7df))
* **shell:** rename context commands from 'set' to 'context' syntax ([#41](https://github.com/mholtzscher/ugh/issues/41)) ([fa9272c](https://github.com/mholtzscher/ugh/commit/fa9272c5ede58555276efe87c35c48ba7a47bde9))
* **tui:** add natural-language command mode for tasks ([#13](https://github.com/mholtzscher/ugh/issues/13)) ([346fdad](https://github.com/mholtzscher/ugh/commit/346fdad106d4a8be86b84b647f05c7bd39cfb533))
* **tui:** componentize task panes and form interactions ([#12](https://github.com/mholtzscher/ugh/issues/12)) ([b148d7c](https://github.com/mholtzscher/ugh/commit/b148d7cb9395755d8064f0dd695d46bf7e9ba508))
* **tui:** hide redundant state column when filtered by state ([f466293](https://github.com/mholtzscher/ugh/commit/f466293208bf7b7a8573d2d84a0309979a91131f))
* **tui:** remove TUI functionality ([#15](https://github.com/mholtzscher/ugh/issues/15)) ([a913cc8](https://github.com/mholtzscher/ugh/commit/a913cc8c60002b0792a7cc06f9e027c2b0d7b3f0))


### Bug Fixes

* **nlp:** generate enum String methods for diagnostics ([#81](https://github.com/mholtzscher/ugh/issues/81)) ([62d2a48](https://github.com/mholtzscher/ugh/commit/62d2a48581c71269a43530f036af1a7c5ada5a54))
* **nlp:** require explicit command verbs in parser ([#17](https://github.com/mholtzscher/ugh/issues/17)) ([ba20af7](https://github.com/mholtzscher/ugh/commit/ba20af71ad2d46e1a1a8c5c94d1f763201abc805))
* **nlp:** support repeated filter predicates in query compilation ([#49](https://github.com/mholtzscher/ugh/issues/49)) ([6b922b6](https://github.com/mholtzscher/ugh/commit/6b922b6c3f319a240536737bffdbc7a59837b477))
* **shell:** move sticky context injection to AST phase ([#51](https://github.com/mholtzscher/ugh/issues/51)) ([a3985b9](https://github.com/mholtzscher/ugh/commit/a3985b9c6f5ef565767ef53cc351158891bf0283))
* **shell:** reject non-printable control characters in command input ([#80](https://github.com/mholtzscher/ugh/issues/80)) ([e68c2d7](https://github.com/mholtzscher/ugh/commit/e68c2d7587494c8ca6cba6f3134d4be01a164f90))
* **state:** normalize task state handling across parsing ([#79](https://github.com/mholtzscher/ugh/issues/79)) ([b457c6d](https://github.com/mholtzscher/ugh/commit/b457c6d00dd2b5a302ada10c0671c776a0478f56)), closes [#65](https://github.com/mholtzscher/ugh/issues/65)

## [0.1.1](https://github.com/mholtzscher/ugh/compare/v0.1.0...v0.1.1) (2026-02-04)


### Features

* **cli:** add command aliases for improved user experience ([dbe0986](https://github.com/mholtzscher/ugh/commit/dbe0986ed754f9cd7e0b252bdda6bd064e2fe267))
* **cli:** add projects and contexts commands to list available tags ([e31b1f7](https://github.com/mholtzscher/ugh/commit/e31b1f724eba7d2864a1b887c11d5af7bb2c2b96))
* **cli:** add short flags for command options ([bc4f756](https://github.com/mholtzscher/ugh/commit/bc4f7562858bea89370ba6c9ff45ff58f917f964))
* **config:** add configuration management commands ([bbd24b6](https://github.com/mholtzscher/ugh/commit/bbd24b6c50260b3c340528e9f2e794b0d234f790))
* **config:** add sync_on_write configuration option for automatic syncing ([99850bb](https://github.com/mholtzscher/ugh/commit/99850bb783ae9ff30398ce94c2d539e6a732d9eb))
* **config:** add TOML configuration file support with path resolution ([f35516b](https://github.com/mholtzscher/ugh/commit/f35516bec4b6740ca3b88b68f54b6b31f9a308b9))
* **daemon:** implement daemon service for background Turso sync with systemd/launchd support ([ce2a1f5](https://github.com/mholtzscher/ugh/commit/ce2a1f59b44c1932fb248ee5fad95f666dc1aa05))
* **edit:** add partial updates and interactive editor mode ([cba617a](https://github.com/mholtzscher/ugh/commit/cba617a7c3a6aeca91df4be483c5c02594cfa89e))
* ensure creation dates are set for all tasks ([79522f0](https://github.com/mholtzscher/ugh/commit/79522f0067ec74ba847546fef4b2e2d8f1dff188))
* implement initial task CLI with SQLite storage ([7a05af9](https://github.com/mholtzscher/ugh/commit/7a05af9f6ac9b4fcecdaa09aa56d00244f775393))
* migrate database backend from SQLite to Turso libSQL ([c02c8bd](https://github.com/mholtzscher/ugh/commit/c02c8bdbd808758486a08f4d2623328dba64b62b))
* migrate to sqlc for type-safe queries and goose for migrations ([a600bb9](https://github.com/mholtzscher/ugh/commit/a600bb976defc431004094d0ac019331ecba90e5))
* **output:** improve terminal table layout with dynamic column widths and responsive sizing ([be1a47d](https://github.com/mholtzscher/ugh/commit/be1a47d8e7acdd0286572657652c370437a1e55a))
* **root:** use OS-specific data directories for default database path ([a0ba545](https://github.com/mholtzscher/ugh/commit/a0ba5454ddd313ffa4c660cd4bd308696ae073ec))


### Bug Fixes

* suppress unused error from service Close calls ([e9616a8](https://github.com/mholtzscher/ugh/commit/e9616a8025c5eaa4b7007f2b63bb14eeb26b7551))

## [0.1.0](https://github.com/mholtzscher/ugh/releases/tag/v0.1.0) (YYYY-MM-DD)

### Features

- Initial release
- Cobra-based CLI for task workflows
- SQLite/Turso-backed storage
- Nix flake support
- GitHub Actions CI/CD
