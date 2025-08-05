<!--
  README.md – Silva Suite
-->

<h1 align="center">Silva Suite 🌲</h1>
<p align="center">
  An open-source Go ecosystem that <strong>reads calendars</strong>,
  <strong>decides when to book</strong>, and <strong>automatically reserves classes or services</strong>.
</p>

<p align="center">
  <a href="https://github.com/silvasuite/silva/actions"><img alt="CI" src="https://github.com/silvasuite/silva/actions/workflows/ci.yml/badge.svg"></a>
  <a href="https://pkg.go.dev/github.com/silvasuite/silva"><img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/silvasuite/silva.svg"></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/github/license/silvasuite/silva"></a>
</p>

---

## ✨ Key Features

* **Hexagonal architecture** – clear separation between core logic and adapters.
* **Pluggable providers** – swap in any calendar or booking service by implementing tiny interfaces.
* **Multiple UIs** – run head-less from cron, via CLI, or through a Telegram bot.
* **Pure Go** – single binary, cross-platform, no CGO.

---

## Repository Layout

| Path        | Status    | Description                                   |
| ----------- | --------- | --------------------------------------------- |
| `trunk/`    | ✅ stable  | Core domain, ports, and scheduling engine     |
| `internal/` | ✅ stable  | Helper packages *not* importable from outside |
| `canopy/`   | 🛠 planned | Calendar adapters (Google, CalDAV, …)         |
| `hermit/`   | 🛠 planned | Telegram bot UI                               |
| `roots/`    | 🛠 planned | CLI / daemon runner                           |

> The tree is concise on purpose — each directory is a self-contained Go package.

---

## Quick Start (core only)

```bash
go get github.com/silvasuite/silva@latest
```

```go
package main

import (
    "context"
    "log"

    "github.com/silvasuite/silva/trunk/usecase"
)

func main() {
    svc := trunk.Service{
        Calendar: myCalendarProvider, // implement ports.CalendarProvider
        Booking:  myBookingProvider,  // implement ports.BookingProvider
    }
    if err := svc.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

Run unit tests:

```bash
go test ./...
```

---

## Roadmap

- [ ] **Silva Canopy** – Google Calendar & CalDAV adapters
- [ ] **Silva Roots** – YAML-driven CLI / Daemon runner
- [ ] **Silva Hermit** – Telegram bot with inline commands
- [ ] Rule-based slot selection (time windows, max frequency)
- [ ] Plug-in system for custom booking providers

---

## Contributing

1. Fork, create a feature branch, and code away.  
2. Make sure `go test ./...` and `golangci-lint run` pass.  
3. Open a PR - we do squash & merge.

See `CONTRIBUTING.md` (coming soon) for details.


---

## License

Silva Suite is released under the **Apache License 2.0**.  
See [LICENSE](LICENSE) for the full text.
