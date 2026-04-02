# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go package for communicating with Prologix GPIB-USB and GPIB-ETHERNET controllers, which bridge a computer to GPIB (IEEE 488) test instruments. Part of the [gotmc](https://github.com/gotmc) ecosystem and designed to work with the [ivi](https://github.com/gotmc/ivi) package for standardized instrument control.

## Build & Test Commands

```bash
# Format and vet (runs automatically before tests)
just check

# Run unit tests with coverage and race detection
just unit

# Lint with golangci-lint
just lint

# Run a single test
go test -run TestIsPrimaryAddressValid ./...

# HTML coverage report
just cover
```

`just` (Justfile) is the task runner; there is no Makefile.

## Architecture

The package uses a layered design separating transport from GPIB protocol logic:

- **Transport layer** (`driver/`): Provides `io.ReadWriter` implementations for different connection types. Currently only `driver/vcp/` (Virtual COM Port via serial) is implemented. The `driver/usb/` and `driver/net/` directories are placeholders for future D2XX and Ethernet support.
- **Controller** (`controller.go`): Core type wrapping any `io.ReadWriter`. Constructed via `NewController()` with functional options (`WithSecondaryAddress`, `WithDebug`, `WithAR488`). Handles GPIB address setup, controller initialization commands, and the `++` command prefix for Prologix-specific commands.
- **Commands** (`commands.go`): Methods on `Controller` for Prologix `++` commands (e.g., `ClearDevice`, `SetReadTimeout`, `Version`, `InstrumentAddress`).

Key design points:
- `Controller` distinguishes between instrument communication (`Command`, `Query`, `Write`) and Prologix controller commands (`CommandController`, `QueryController`) which prepend `++`.
- Context-aware I/O: `Controller` checks if the underlying `io.ReadWriter` implements optional `ContextReader`/`ContextWriter` interfaces. If so, it delegates directly; otherwise it falls back to goroutine-based context handling.
- The `WithAR488()` option skips `verbose` and `savecfg` commands for Arduino AR488 compatibility.
- `Query` automatically sends `++read eoi` when auto read-after-write is disabled.
- Binary data methods (`Write`/`Read`) vs string methods (`WriteString`, `Command`, `Query`) handle GPIB character escaping differently.
- VCP serial config is hardcoded to Prologix defaults: 115200 baud, 7 data bits, even parity, 1 stop bit.

## Dependencies

- `github.com/gotmc/query` — SCPI query helpers
- `go.bug.st/serial` — cross-platform serial port access
- `go.uber.org/multierr` — error combining
