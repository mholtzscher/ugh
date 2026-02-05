# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
