# 🎓 AIUB Companion

[![Go Version](https://img.shields.io/badge/Go-1.26.0-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![React Version](https://img.shields.io/badge/React-19.2-61DAFB?style=flat&logo=react&logoColor=white)](https://react.dev/)
[![Vite Version](https://img.shields.io/badge/Vite-8.0-646CFF?style=flat&logo=vite&logoColor=white)](https://vite.dev/)
[![Wails Version](https://img.shields.io/badge/Wails-v3--alpha-red?style=flat)](https://v3.wails.io/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-v4.3-38B2AC?style=flat&logo=tailwindcss)](https://tailwindcss.com/)
[![SQLite](https://img.shields.io/badge/SQLite-FTS5-003B57?style=flat&logo=sqlite&logoColor=white)](https://sqlite.org/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**AIUB Companion** is a highly optimized, offline-first, modern desktop application designed for students and faculty of the American International University-Bangladesh (AIUB). It streamlines daily academic activities by scraping official university notices, managing routine class schedules from Excel sheets, and providing an array of interactive academic tools.

Built on the cutting-edge **Wails v3 framework**, AIUB Companion features a blazing-fast pure-Go backend integrated with a fluid React 19 single-page application, styled gracefully with Tailwind CSS v4.

---

## 🚀 Key Features

### 📡 Offline-First Notice Scraper & Sync

- **Automatic Background Synchronization:** Keeps you updated with the latest AIUB notices without manually visiting the website.
- **Fast Full-Text Search (FTS):** Search notices instantly using the SQLite FTS5 engine.
- **Attachment Downloader:** Easily browse and download associated notice PDF/image attachments locally.
- **Native OS Notifications:** Receive system notifications the moment a new notice is published.

### 📅 Academic Routine & Class Schedule

- **Excel Importer:** Import your class or offered course schedule directly from official AIUB spreadsheet files (using the pure-Go `excelize` library).
- **Flexible Searching & Filtering:** Search by Course Title, Course Code, Class ID, Faculty name, Day, Room, or Department.
- **Conflict Checker & Day View:** Quickly check what classes you have on any given day.

### 🧮 Academic Utilities

- **CGPA Calculator:** Dynamic and interactive tracker to calculate and forecast your Cumulative GPA.
- **GPA Trend Analyzer:** Beautifully visualize your academic progress and GPA trajectory across multiple semesters.
- **Exam Countdown:** Configurable timers and countdowns for midterms, finals, or other critical deadlines.

### ⚙️ Deep System Integration

- **System Tray/App Indicator Support:** Minimize to system tray, launch minimized, or close-to-tray to stay lightweight in the background.
- **Flexible Configuration:** Customizable background sync intervals, notification preferences, system startup behavior, logging details, and dark/light/system theme toggles.
- **Cross-Platform:** Built with cross-compilation in mind, supporting Windows, macOS, and Linux (with exploratory setups for mobile Android/iOS).
- **Headless Server Mode:** Supports compiling and running as a lightweight GUI-less HTTP server containerizable via Docker.

---

## 🛠️ Technology Stack

| Backend (Go)                       | Frontend (React & TS)                  | Database & Tooling                  |
| ---------------------------------- | -------------------------------------- | ----------------------------------- |
| **Go 1.26.0**                      | **React 19** & **TypeScript 6/7**      | **SQLite** (Pure Go, Zero-CGO)      |
| **Wails v3** (Alpha 2.106)         | **Vite 8** & **Rolldown**              | **SQLC** (Type-safe SQL generation) |
| **Excelize v2** (Excel parsing)    | **Tailwind CSS v4**                    | **Taskfile** (Advanced Task runner) |
| **Bluemonday** (HTML sanitization) | **TanStack React Query v5**            | **Docker** (Cross-compiling/Server) |
| **Slog** (Structured logging)      | **Base UI (Radix)** & **Lucide React** | **Prettier & ESLint**               |

---

## 📂 Project Structure

A guide to the main directories and components of the application:

```text
├── .task/                  # Taskfile execution cache
├── bin/                    # Compiled production executables
├── build/                  # Platform-specific build assets (icons, configs, installers, Dockerfiles)
├── frontend/               # React frontend application
│   ├── src/
│   │   ├── components/     # Reusable UI components (custom inputs, buttons, etc.)
│   │   ├── hooks/          # React hooks (use-notices, use-debounce, etc.)
│   │   ├── lib/            # Shared utilities, logging, routing, and configurations
│   │   ├── pages/          # Layout pages (Notices, Routine, CGPA, Exam Countdown, Settings, Help)
│   │   ├── App.tsx         # Main application component & routes
│   │   └── main.tsx        # React entry point
│   ├── package.json        # Frontend dependencies (React 19, Tailwind v4, etc.)
│   └── vite.config.ts      # Vite compilation configuration
├── internal/               # Go backend package modules (private domain logic)
│   ├── config/             # Config service, validation, schemas, and launch rules
│   ├── database/           # SQLite connection pool, migrations, and transactions
│   │   ├── db/             # Generated SQLC Go files (models, querier)
│   │   └── sql/            # Raw SQL schemas and queries
│   ├── desktop/            # Native window handling, single-instance lock, system tray
│   ├── log/                # Custom structured logging and logger middleware
│   ├── notice/             # Scraper client, database repositories, notice services
│   ├── routine/            # Offered course Excel importer, routine queries, service binds
│   └── worker/             # Background synchronization ticker & OS notification dispatch
├── main.go                 # Entry point of the Go backend
├── sqlc.yaml               # SQLC compiler configuration
└── Taskfile.yml            # Main automation task configuration
```

---

## 🚀 Getting Started

Ensure you have the following installed on your machine:

- **Go** (1.26.0 or higher)
- **Node.js** (v18 or higher) and **PNPM** (or npm/yarn/bun)
- **Wails v3 CLI** (Required for hot-reload dev mode)
- **Task** (Highly recommended task runner: `go install github.com/go-task/task/v3/cmd/task@latest`)

### 1. Installation & Setup

Clone the repository and install the frontend dependencies:

```bash
# Clone the repository
git clone https://github.com/atifchy/aiub-companion.git
cd aiub-companion

# Install frontend dependencies (automatically managed if using Task)
cd frontend
pnpm install
cd ..
```

### 2. Running in Development Mode

Run the application with real-time hot-reloading for both backend Go code and frontend React UI:

```bash
# Using Task (Recommended)
task dev

# Or directly using Wails3
wails3 dev
```

### 3. Compilation & Building

To generate optimized binary releases for your current Operating System:

```bash
# Production Build
task build

# Packaging an Installer (e.g. MSI on Windows, AppImage/Deb on Linux, DMG/App on macOS)
task package
```

Built executables are created in the `bin/` directory.

---

## 🌐 Headless Server Mode (Docker Support)

AIUB Companion can be run as a headless, UI-less server that scrapes notices and exposes endpoints via an HTTP API.

To build and deploy using Docker:

```bash
# Build the production server mode Docker container
task build:docker

# Run the Docker container
task run:docker
```

The server is configured using the Dockerfile located at `build/docker/Dockerfile.server`.

---

## 📝 Development Guide

- **Database Migrations & SQLC:** If you edit schema files in `internal/database/sql/schemas/` or queries in `internal/database/sql/queries/`, regenerate the SQLC bindings:

  ```bash
  sqlc generate
  ```

- **Logging:** Structured logging is implemented through Go's native `slog` library. Log files can be found locally in the application's storage path under the `logs` folder.
- **Configurations:** The application uses a local JSON configuration file parsed dynamically with JSonSchema validation (`internal/config/schema.json`).

---

## 🛡️ License

This project is licensed under the **Apache License 2.0**. See the [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing

Contributions are highly welcomed! Please feel free to open an issue or submit a pull request if you want to enhance the UI, add features, improve scraping robustness, or fix bugs.

_Disclaimer: AIUB Companion is an unofficial fan-made companion tool and is not affiliated, associated, authorized, endorsed by, or in any way officially connected with the American International University-Bangladesh (AIUB)._
