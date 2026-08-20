# 🎓 AIUB Companion

[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react&logoColor=white)](https://react.dev/)
[![Wails](https://img.shields.io/badge/Wails-v3--beta-red?style=flat&logo=wails)](https://v3.wails.io/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-v4-38B2AC?style=flat&logo=tailwindcss)](https://tailwindcss.com/)
[![License](https://img.shields.io/badge/License-Apache_2.0-Green.svg)](https://opensource.org/licenses/Apache-2.0)

A desktop app for AIUB students that pulls in university notices, lets you manage class schedules from Excel files, and keeps your academic calendar handy. Everything works offline and stays lightweight.

Built with Go + Wails v3 on the backend and React 19 + Tailwind CSS v4 on the frontend.

## Features

### 📢 Notices

- Syncs official AIUB notices automatically in the background
- Fuzzy search that still finds what you're looking for even with typos
- View and download PDF/image attachments right from the app
- Pin important notices and track what you've read
- Desktop notifications when new notices show up

### 📅 Class Routine

- Import course offerings directly from AIUB's Excel spreadsheets
- Search courses by name, code, section, faculty, or department
- Weekly timetable view with live "Ongoing" and "Up Next" indicators

### 📆 Academic Calendar

- Scrapes and caches the official AIUB academic calendar so you can check it anytime

### ⚙️ Settings and System

- Sits quietly in your system tray when minimized
- Only one instance runs at a time (opening it again just brings back the existing window)
- Optional start-on-boot
- Auto-updates from GitHub Releases on a schedule you pick (daily, weekly, monthly, or manual)
- Configurable themes, sync intervals, notification preferences, and more

## Tech Stack

**Backend (Go 1.27)**

- Wails v3 (beta) for the desktop shell and frontend/backend bridge
- SQLite via `modernc.org/sqlite` (pure Go, no CGO needed)
- SQLC for type-safe SQL query generation
- `excelize` for parsing Excel files
- `go-edlib` for Jaro-Winkler fuzzy matching
- `bluemonday` for HTML sanitization
- `goquery` for scraping
- Structured logging with `slog` + `tint`

**Frontend (React 19 + TypeScript)**

- Vite 7 for bundling and dev server
- Tailwind CSS v4 for styling
- shadcn/ui components (Radix primitives + Base UI)
- React Router v8 (hash-based routing)
- TanStack Query v5 for data fetching
- Zustand for client-side state
- Lucide React icons, Sonner toasts, Motion for animations
- `marked` + `DOMPurify` for rendering notice content

**Tooling**

- Task (task runner)
- oxlint + oxfmt (linting and formatting)
- Docker support for headless server builds

## Project Structure

```
aiub-companion/
├── build/                   # Build assets, icons, Docker configs, platform manifests
├── frontend/
│   └── src/
│       ├── components/      # UI components (notices, routine, settings, sidebar, etc.)
│       ├── hooks/           # React hooks (notices, config, debounce, stores)
│       ├── lib/             # Routing, utilities, lazy loading helpers
│       └── pages/           # Notices, Routine, Semester, Settings, Help, About
├── internal/
│   ├── calendar/            # Academic calendar scraper and cache
│   ├── config/              # JSON config with schema validation
│   ├── database/            # SQLite connection, schema embedding, SQLC output
│   ├── desktop/             # Window management, tray, single-instance lock
│   ├── fetcher/             # Shared HTTP client helpers
│   ├── log/                 # File and console logger
│   ├── meta/                # App ID, version, repo constants
│   ├── notice/              # Notice scraper, repository, fuzzy search
│   ├── persist/             # Shared persistence utilities
│   ├── routine/             # Excel importer, course cache, search
│   ├── search/              # Generic fuzzy search (Jaro-Winkler)
│   ├── tz/                  # Timezone helpers
│   ├── updater/             # Auto-updater using Wails v3
│   └── worker/              # Background sync and notification dispatch
├── main.go                  # Entry point, service wiring
├── Taskfile.yml             # Build and dev task definitions
└── sqlc.yaml                # SQLC code generation config
```

## Getting Started

### You'll Need

- [Go 1.27+](https://go.dev/)
- [Bun](https://bun.sh/)
- [Wails v3 CLI](https://v3.wails.io/): `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
- [Task](https://taskfile.dev/): `go install github.com/go-task/task/v3/cmd/task@latest`

### Development

```bash
git clone https://github.com/atifchy/aiub-companion.git
cd aiub-companion

# Start dev mode with hot reload
task dev
```

### Production Build

```bash
# Build the desktop binary
task build

# Build and package an installer (NSIS on Windows)
task package
```

Output goes to `bin/`.

## License

Apache License 2.0. See [LICENSE](LICENSE) for the full text.

> *AIUB Companion is a community project and is not affiliated with or endorsed by American International University-Bangladesh (AIUB).*
