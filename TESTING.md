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

**Prerequisite (Local SSH Setup):**
Because the python script relies on `autossh` using your current Mac Username to tunnel back into your *own* Mac (simulating a remote connection), you must have two things configured locally:
1. **Remote Login** enabled in macOS (`System Settings -> General -> Sharing`).
2. **Passwordless SSH** configured for your own account. Run the following command once to authorize yourself to SSH into your own machine without a password:
   ```bash
   cat ~/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys
   chmod 600 ~/.ssh/authorized_keys
   ```

**Step 1:** In a new terminal tab, initialize a virtual environment and run the agent script:
```bash
cd agent
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python node_agent.py --name mock-bench-phase2
```
*Expected Output:* The python log will show it registered with the Go server, found its Assigned SSH Port (e.g., 22002), and successfully initiated an `autossh` tunnel!

**Step 2:** Verify the reverse tunnel is active:
In another tab, run:
```bash
lsof -i -P -n | grep LISTEN | grep 2200
```
*Expected Output:* You should see `IPv4` and `IPv6` listening ports for the port assigned to you by the Python script.

## Phase 3: Testing the CLI

**Step 1:** Build the CLI tool:
```bash
cd cli
go build -o hilcli .
```

**Step 2:** Query the Benches using the CLI Table:
```bash
./hilcli list
```
*(You should see a clean Terminal visualization of the PostgreSQL database, showing any active agent machines!)*

**Step 3:** Simulate the Developer Experience
If your Python mock agent is still running and listening from Phase 2, try to use the CLI abstraction wrapper to transparently reverse-jump into it:
```bash
./hilcli connect mock-bench-phase2
```
> **Note:** Just like the Phase 2 Testing, this will proxy-jump through your `localhost` back into yourself. Once inside the remote terminal, type `exit` to end the test session!
