# Hardware-in-the-Loop (HIL) Test Infrastructure

A scalable, fully open-source, self-hosted software architecture for managing geographically distributed Hardware-in-the-Loop test benches.

## 🎥 Demos

### 🌐 Full Stack Integration
This demo showcases the complete workflow: from registration of a mock HIL node to visualization in the Next.js dashboard and real-time status updates through the Go backend.

![Full Stack Integration](demos/full_stack_demo.mp4)

### 💻 CLI Access & Connectivity
See the `hilcli` tool in action, demonstrating how developers can interact with remote test benches and initiate secure reverse-proxy connections.

![CLI Access & Connectivity](demos/hilclidemo.mp4)

## 📋 Prerequisites

### Local Testing (Development Machine)
To run the entire infrastructure stack locally for development, you will need:
- **Docker**: (e.g., OrbStack, Docker Desktop, or Colima) to run the `docker-compose.yml` PostgreSQL database.
- **Go 1.20+**: To compile and run the central backend API server and the `hilcli` tool.
- **Python 3.9+**: To run the mock HIL node agent scripts locally.
- **Node.js**: (v18+) Required for running the Next.js React Web Dashboard.
- **OpenSSH Daemon (Remote Login)**: Since testing locally treats your laptop as both the Server and the HIL Node, you must enable "Remote Login" under your macOS Sharing settings, and ensure your own public SSH key is appended to your `~/.ssh/authorized_keys` file to support passwordless `autossh` routing.

### Production Deployment (Geographically Distributed)
- **Central Management Server**: A Linux VPS with Docker (for PostgreSQL) and a static public IP address. It will host the compiled Go binary and terminate all reverse SSH tunnels.
- **HIL Bench Nodes**: One physical Linux machine per test bench, connected to the DUT (Device Under Test). Must have Python 3, `autossh`, and the native OpenSSH server (`sshd`) installed.

## 🚀 Quick Start
For detailed instructions on testing each component locally without physical hardware, please refer to [TESTING.md](./TESTING.md).

To bootstrap the local control plane:
```bash
# 1. Start the PostgreSQL database
docker-compose up -d

# 2. Run the Go backend server (port 8080)
cd backend && go run main.go

# 3. Start the Next.js Web Dashboard (in a new terminal)
cd frontend
npm install
npm run dev
```
