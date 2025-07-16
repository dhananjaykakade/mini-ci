# 🚀 Deployment Dashboard

A modern and intuitive web-based dashboard that allows you to deploy applications directly from a Git repository. Built with **React + Vite + ShadCN UI**, this app supports various tech stacks like Node.js, React, Flask, Go, Java, and more.

![Demo video](./demo.mp4)

---

## ✨ Features

- 🌐 Deploy directly from a GitHub repo
- 🧠 Auto-filled commands based on selected app type
- 🧾 Add custom install, build, and start commands
- ⚙️ Environment variables support
- 🔁 Real-time deployment logs (streamed from backend)
- ✅ Deployment success message with URL
- 🔄 Re-deploy button to launch another deployment
- 🧱 Beautiful modern UI (Monochrome with gradient highlights)

---

## 🛠️ Tech Stack

- **Frontend**: Vite + React + TypeScript
- **UI Components**: ShadCN UI (Radix + Tailwind)
- **Icons**: Lucide React
- **Notifications**: Sonner
- **Backend**: Node.js + Express (with build-stream endpoint)
- **Deployment**: Docker (assumed for container-based deployment)

---

## 📦 Installation

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/deploy-dashboard.git
cd deploy-dashboard
