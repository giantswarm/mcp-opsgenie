# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [Unreleased]

### Added

- Adopt `github.com/giantswarm/mcp-toolkit` v0.1.0 for cross-cutting plumbing. `cmd/serve.go` now imports `logging.New`, `tracing.Init`, `health.New`, `httpx.Run`, `responsecap.New`, and `timeout.New`. `responsecap` (default 128 KiB) and `timeout` (30s) middleware are wired on every tool call. SSE and streamable-HTTP transports host their own `*http.Server` and run via `httpx.Run`; `/healthz` and `/readyz` are mounted on the same mux. OTEL tracing is best-effort: when `OTEL_EXPORTER_OTLP_*` is unset, `tracing.Init` returns a no-op shutdown but still installs the W3C propagator so inbound `traceparent` headers chain.

### Changed

- Logger initialization moved to the toolkit's `logging.New`. For HTTP transports the default destination is now stderr (auto JSON in-cluster, text locally). Stdio still defaults to discard so it cannot pollute the protocol stream. `--log-file` continues to override both — the file path becomes the logger's destination regardless of transport.

[Unreleased]: https://github.com/giantswarm/mcp-opsgenie/tree/main
