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

---

## Phase 2: Testing the Node Agent
*(Implementation in progress...)*

Eventually, you will be able to start multiple mock nodes by running multiple instances of the Python script in separate terminal tabs. The agents will automatically register with `localhost:8080`, receive their port assignment, and execute an `autossh` loop back to your laptop.

## Phase 3: Testing the CLI
*(Implementation in progress...)*

You will use `hilcli list` and `hilcli connect <hostname>` directly from your terminal. The CLI will abstract away the port logic and transparently drop you into an interactive session over the reverse SSH tunnel.
