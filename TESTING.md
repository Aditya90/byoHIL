# Local Testing Guide

Because this architecture separates concerns between a central control server and geographically distributed test nodes, **the entire stack is designed to be fully testable on a single `localhost` machine.**

## Phase 1: The Control Plane (Go Backend & DB)
The core of the infrastructure is the Go API server and its PostgreSQL database.

**Step 1:** Spin up the database background container.
```bash
docker-compose up -d
```

**Step 2:** Start the Go Backend Server.
```bash
cd backend
go run main.go
```

**Step 3:** Verify the API is awake and assigning ports correctly by mocking a node registration request from a new terminal tab.
```bash
curl -X POST -H 'Content-Type: application/json' -d '{"hostname": "mock-bench-test"}' http://localhost:8080/api/v1/nodes/register
```
*Expected Output:* `{"assigned_ssh_port":22000,"status":"registered"}`

**Step 4:** Query the Registered Agents (The foundation for the CLI and Web Dashboard).
```bash
curl -s http://localhost:8080/api/v1/nodes
```
*Expected Output:* A JSON array containing the `mock-bench-test` node you just registered.

**Step 5:** Simulate an Access Log (The foundation for tracking who connects to the bench).
```bash
curl -X POST -H 'Content-Type: application/json' -d '{"username": "aditya", "action": "cli_connect"}' http://localhost:8080/api/v1/nodes/mock-bench-test/access_log
```
*Expected Output:* `{"status":"logged"}`

---

## Phase 2: Testing the Node Agent

You will run a mock node agent on your laptop to verify the registration protocol and reverse SSH implementation.

**Step 1:** In a new terminal tab, initialize a virtual environment and run the mock script:
```bash
cd agent
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python mock_agent.py --name mock-bench-phase2
```
*Expected Output:* The python log will show it registered with the Go server, found its Assigned SSH Port (e.g., 22002), and successfully initiated an `autossh` tunnel!

> **Note:** Because the script uses your current username to SSH into your own laptop to test the tunnel locally, you must have "Remote Login" enabled in macOS `System Settings -> General -> Sharing` and optionally have run `ssh-copy-id localhost` to enable passwordless auth.

**Step 2:** Verify the reverse tunnel is active:
In another tab, run:
```bash
lsof -i -P -n | grep LISTEN | grep 2200
```
*Expected Output:* You should see `IPv4` and `IPv6` listening ports for the port assigned to you by the Python script.

## Phase 3: Testing the CLI
*(Implementation in progress...)*

You will use `hilcli list` and `hilcli connect <hostname>` directly from your terminal. The CLI will abstract away the port logic and transparently drop you into an interactive session over the reverse SSH tunnel.
