# 🎓 AIUB Companion

[![Go Version](https://img.shields.io/badge/Go-1.26.0-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![React Version](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react&logoColor=white)](https://react.dev/)
[![Vite Version](https://img.shields.io/badge/Vite-8-646CFF?style=flat&logo=vite&logoColor=white)](https://vite.dev/)
[![Wails Version](https://img.shields.io/badge/Wails-v3--alpha-red?style=flat&logo=wails)](https://v3.wails.io/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-v4-38B2AC?style=flat&logo=tailwindcss)](https://tailwindcss.com/)
[![SQLite](https://img.shields.io/badge/SQLite-DB-003B57?style=flat&logo=sqlite&logoColor=white)](https://sqlite.org/)
[![License](https://img.shields.io/badge/License-Apache_2.0-Green.svg)](https://opensource.org/licenses/Apache-2.0)

**AIUB Companion** is an offline-first desktop application for students of the American International University-Bangladesh (AIUB). It scrapes official university notices, manages class schedules from Excel sheets, and provides academic utilities — all in a native, lightweight desktop app.

Built with **Wails v3**, a pure-Go backend, and a React 19 SPA frontend styled with Tailwind CSS v4.

---

## 🚀 Features

### 📡 Notice Board

- **Background Sync** — Automatically fetches the latest AIUB notices at a configurable interval.
- **Fuzzy Full-Text Search** — JaroWinkler-based in-memory fuzzy search.
- **Attachment Viewer** — Browse and download PDF/image attachments associated with each notice.
- **Pin & Read State** — Mark notices as read or pin important ones; state persists locally.
- **Native Notifications** — System-level OS notifications on new notice arrival.

### 📅 Class Routine

- **Excel Importer** — Import offered courses directly from official AIUB spreadsheet files (`excelize`).
- **Fuzzy Course Search** — In-memory JaroWinkler fuzzy search across course title, code, section, faculty, and department — powered by an in-memory cache for zero DB round-trips while typing.
- **Week Timeline View** — Visual day-by-day schedule with live "Ongoing" and "Up Next" status indicators.

### 📆 Academic Calendar

- Scrapes and caches the official AIUB academic calendar locally.

### ⚙️ System Integration

- **System Tray** — Minimize to tray, close-to-tray, and tray-click to restore.
- **Single Instance Lock** — Prevents multiple instances; focuses existing window instead.
- **Autostart** — Optional Windows startup shortcut via `WScript.Shell` COM automation.
- **Auto-Update** — Powered by Wails v3's built-in updater with GitHub Releases as the provider. Supports daily, weekly, or monthly scheduled checks with persistent last-check state.
- **Deep Configuration** — Theme, log level, sync interval, notification toggle, launch behaviour — all persisted to `config.json` with JSON Schema validation.

---

## 🛠️ Technology Stack

| Backend (Go) | Frontend (React & TS) | Tooling |
| --- | --- | --- |
| **Go 1.26.0** | **React 19** + **TypeScript** | **SQLite** (Pure Go, Zero-CGO via `modernc.org/sqlite`) |
| **Wails v3** (`alpha2.117`) | **Vite 8** | **SQLC** (type-safe SQL → Go codegen) |
| **excelize v2** (Excel parsing) | **Tailwind CSS v4** | **Task** (task runner) |
| **go-edlib** (JaroWinkler fuzzy) | **TanStack Query v5** | **Docker** (cross-compile & server mode) |
| **bluemonday** (HTML sanitizer) | **Radix UI** + **Lucide React** | **PNPM** |
| **slog** + **tint** (structured logging) | **Sonner** (toasts) | **Prettier** + **ESLint** |
| **go-ole** (Windows COM) | | |

---

## 📂 Project Structure

```text
aiub-companion/
├── build/                   # Platform build assets (icons, Taskfiles, Dockerfiles, installers)
├── frontend/                # React SPA
│   └── src/
│       ├── components/      # Reusable UI components
│       ├── hooks/           # Custom React hooks (use-notices, use-debounce, use-updater, ...)
│       ├── lib/             # Shared utilities (logger, router, utils)
│       └── pages/           # Route-level pages (Notices, Routine, Calendar, Settings, ...)
├── internal/
│   ├── autostart/           # Manages start-on-boot shortcut creation/removal
│   ├── calendar/            # Academic calendar scraper & cache
│   ├── config/              # JSON config with schema validation; persists to config.json
│   ├── console/             # Wails devtools console bridge (dev-only)
│   ├── database/            # SQLite connection, migrations
│   │   ├── db/              # SQLC-generated Go models & queries
│   │   └── sql/
│   │       ├── queries/     # Raw SQL query files
│   │       └── schemas/     # Schema migration files (01_schema.sql, 02_indexes.sql)
│   ├── desktop/             # Window management, tray icon, single-instance
│   ├── event/               # Shared event name constants
│   ├── log/                 # Structured file+console logger (slog + tint)
│   ├── meta/                # App constants (name, version, repo)
│   ├── notice/              # Notice scraper, repository, service, fuzzy search
│   ├── routine/             # Course Excel importer, in-memory cache, fuzzy search
│   ├── search/              # Generic JaroWinkler fuzzy search (FuzzySearch[T])
│   ├── updater/             # Wails v3 updater wrapper; scheduled checks; state persistence
│   └── worker/              # Background sync ticker; notification dispatch
├── main.go                  # App entry point & service wiring
├── sqlc.yaml                # SQLC compiler config
├── Taskfile.yml             # Task runner definitions
└── go.mod
```

---

## 🚀 Getting Started

### Prerequisites

- **Go** 1.26.0+
- **Node.js** 26+ and **PNPM**
- **Wails v3 CLI** — `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
- **Task** — `go install github.com/go-task/task/v3/cmd/task@latest`
- **SQLC** (only needed if modifying SQL) — `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

### Install & Run

```bash
git clone https://github.com/atifchy/aiub-companion.git
cd aiub-companion

# Install frontend dependencies
cd frontend && pnpm install && cd ..

# Run in dev mode (hot-reload for both Go and React)
task dev
```

### Build

```bash
# Production binary for current OS
task build

# Packaged installer (NSIS on Windows, DMG on macOS, AppImage/deb on Linux)
task package
```

Executables are written to `bin/`.

---

## 🗃️ Database

The app uses a local SQLite database via `modernc.org/sqlite` (pure Go, zero CGO).

- **Schemas** live in `internal/database/sql/schemas/` and are embedded + auto-applied on startup.
- **Queries** live in `internal/database/sql/queries/` and are compiled to type-safe Go by SQLC.

After editing any `.sql` file, regenerate:

```bash
sqlc generate
```

---

## 🔄 Auto-Update

Updates are delivered via GitHub Releases using the Wails v3 built-in updater. The updater:

- Filters assets to skip installer packages, selecting only the portable binary.
- Supports configurable check intervals: **daily**, **weekly**, **monthly**, or **never**.
- Persists the last check timestamp so schedules survive restarts.
- Can be triggered manually from the Settings page.

---

## 🖥️ Server Mode

AIUB Companion can be compiled without a GUI as a headless HTTP server (for scraping/deployment):

```bash
task build:server
task run:server

# Or via Docker
task build:docker
task run:docker
```

---

## 🛡️ License

Licensed under the **Apache License 2.0**. See [LICENSE](LICENSE) for details.

---

_AIUB Companion is an unofficial fan-made tool and is not affiliated with or endorsed by the American International University-Bangladesh (AIUB)._
