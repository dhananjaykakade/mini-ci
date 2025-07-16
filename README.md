# ⚡️ Mini-CI

A lightweight, self-hosted **continuous-integration & deployment dashboard**. Point Mini-CI at any public Git repository and it will:

1. Clone the repo.
2. Auto-generate an opinionated Dockerfile (or use the one already present).
3. Build a Docker image while streaming the logs in real-time to the browser.
4. Run the container on a random free port.
5. Show you the live URL and container details – ready for manual smoke-testing or to be wired behind a reverse-proxy.

Mini-CI is designed for hack-day projects and personal side-projects where you want Heroku-like convenience without the monthly bill.



[▶ Watch the demo (MP4)](./demo.mp4)
---

## ✨ Features

- 🔸 **One-click deploy** from any Git URL
- 🧠 **Smart defaults** – common stacks pre-filled
- ⚙️ **Environment variables** editor
- 📜 **Real-time logs** streamed via SSE
- ♻️ **Re-deploy** button for fast iteration
- 🖤 **Beautiful UI** (ShadCN + Tailwind)

---

## 🛠️ Tech Stack

- **Frontend**: Vite + React + TypeScript + Tailwind
- **UI Components**: ShadCN UI (Radix + Tailwind)
- **Icons**: Lucide React
- **Notifications**: Sonner
- **Backend**: Go net/http + Docker CLI
- **CI / Deployment**: Docker

---

## 📋 Prerequisites

- Docker ≥ 24
- Go ≥ 1.22
- Node 18+ and pnpm (or npm / yarn)

---

## 🚀 Getting Started

### 1. Clone the repo

```bash
$ git clone https://github.com/your-username/mini-ci.git
$ cd mini-ci
```

### 2. Start the backend

```bash
$ go run ./backend
```

### 3. Start the frontend

```bash
$ cd frontend
$ pnpm install       # or npm install / yarn
$ pnpm dev
```

Visit `http://localhost:5173` to open the dashboard.

---

## 🔄 Deployment Flow

```text
[Dashboard] ──► POST /build-stream ──► [Backend]
                                  │
                                  └─► git clone ➜ docker build ➜ docker run
                                             │              │
                                          logs ◄────────────┘ (SSE)
```

---

## 📂 Project Structure

```
mini-ci/
├─ backend/         # Go server & CI engine
│  ├─ api/          # HTTP handlers (SSE, health, etc.)
│  └─ runner/       # Docker helpers, git clone, generator, cleaner
├─ frontend/        # React dashboard (Vite)
│  ├─ src/
│  └─ public/
├─ demo.mp4         # Quick demo recording
├─ go.mod / sum     # Backend dependencies
└─ README.md        # ← you are here
```

---

## 📑 API Endpoints

| Method | Path            | Description                           |
| ------ | --------------- | ------------------------------------- |
| GET    | `/health`       | Liveness probe                        |
| GET    | `/test-stream`  | Example SSE stream                    |
| POST   | `/build-stream` | Deploy a repo, returns SSE build logs |
| POST   | `/deploy`       | (Future) production deploy            |
| GET    | `/logs/{id}`    | Fetch stored logs for a past run      |
| GET    | `/ping/{host}`  | Simple TCP ping helper                |

---

## 🤝 Contributing

1. Fork the repo & create a topic branch.
2. Ensure `go vet ./...` and `pnpm lint` pass.
3. Submit a PR explaining _why_ the change is needed.

---

## 📜 License

Released under the MIT License – see [LICENSE](LICENSE) for details.
