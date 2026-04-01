# HIL Implementation Task Tracker & Test Plan

This document tracks the phased implementation of the HIL Infrastructure. It includes explicit intermediate testing milestones to ensure every component functions correctly in isolation before we weave them together.

## 💻 Local Testing Strategy
**All components (Backend, Database, Agent, CLI, and Dashboard) are designed to be run and tested simultaneously on a single development machine (e.g., your laptop) before deploying to physical Linux test benches.**
* **Mock Nodes**: You can simulate multiple benches by running multiple instances of the Python agent locally (e.g., `python agent.py --name mock-bench-01`).
* **Mock Hardware**: The Python agent uses a `MockRelay` plugin that logs hardware commands to the terminal (e.g., `[MOCK_RELAY]: Power Cycled DUT`) instead of requiring physical USB/Serial controllers.
* **Local Tunneling**: The `autossh` daemon and `hilcli connect` will route locally over `localhost` to prove the networking topology works end-to-end.

---

## Phase 1: Foundation – Go Backend & Database
*Goal: Establish the central registry that manages the node topology and port assignments.*

- [x] **1.1 Database Setup**
  - [x] Initialize `docker-compose.yml` with a PostgreSQL service.
  - [x] **Test/Milestone:** Run `docker-compose up` and successfully connect to the database locally using `psql`, verifying the container is healthy.

- [x] **1.2 Go Server Initialization**
  - [x] Run `go mod init backend`.
  - [x] Setup a fast web framework (e.g., Fiber or Gin) and Postgres driver (`gorm`).
  - [x] Implement the Database Models (`Node`, `PortAssignment`, `AccessLog`).
  - [x] Create `/api/v1/nodes/register` endpoint (handles MAC address and capability payloads).
  - [x] Create `/api/v1/nodes` (GET) endpoint to query live machines.
  - [x] Create `/api/v1/nodes/:id/access_log` endpoint to track CLI usage.
  - [x] **Test/Milestone:** Send a manual JSON payload via `curl` to the registration endpoint. Query Postgres to verify the data was saved and an SSH port was uniquely assigned. Send a request to the `access_log` endpoint to verify auditing works.

## Phase 2: Python Agent & Connectivity
*Goal: Connect a physical node to the central registry and establish the reverse tunnel.*

- [x] **2.1 Python Agent Skeleton**
  - [x] Setup Python environment (`requirements.txt` or `pyproject.toml`).
  - [x] Write a script that dynamically fetches the local host's hardware ID/name.
  - [x] Issue an HTTP POST to the Go Server's `/register` endpoint on startup.
  - [x] **Test/Milestone:** Run the Python script on your laptop. Check the Go Backend's logs to see the incoming connection, and check the DB for the new entry.

- [x] **2.2 Reverse SSH Implementation**
  - [x] Python agent parses the `assigned_ssh_port` from the Go server's response.
  - [x] Python agent uses the `subprocess` module to launch and monitor `autossh` (Targeting the central server).
  - [ ] **Test/Milestone:** Start the server and the agent side-by-side. Run `lsof -i -P -n | grep LISTEN` on the server to prove the reverse SSH port (e.g., `22005`) has actually opened. Use a normal SSH test (`ssh -p 22005 localhost`) to verify the tunnel drops you into the agent machine.

## Phase 3: The CLI Wrapper (`hilcli`)
*Goal: Abstract the reverse SSH tunnel complexity away from the developer.*

- [ ] **3.1 CLI Initialization**
  - [ ] Initialize `hilcli` Go project using Cobra for routing commands.
  - [ ] Build `hilcli list` mapping to the Go Backend API.
  - [ ] **Test/Milestone:** Run `go run main.go list` and verify it renders a clean CLI table showing the Python agent from Phase 2.

- [ ] **3.2 Seamless SSH Invocation**
  - [ ] Implement `hilcli connect <hostname>`.
  - [ ] Make the Go CLI fetch the port assignment from the API, construct the `ProxyJump` arguments, and invoke `os/exec` to run native ssh.
  - [ ] **Test/Milestone:** Run `hilcli connect <your-laptop-hostname>`. Verify the CLI drops you instantly into an SSH shell on the target without ever showing you the port number.

## Phase 4: Power Control & Script Execution
*Goal: Enable remote hardware reset functionality triggered globally.*

- [ ] **4.1 Hardware Abstraction (Node Agent)**
  - [ ] Add a mock hardware relay class in the Python agent (which writes "Relay ON" to a file instead of engaging actual USB hardware).
  - [ ] Expose an internal RPC or WebSocket listener to capture commands from the central server.
  - [ ] **Test/Milestone:** Inject a manual command from the Go backend into the stream. Verify the Python agent receives it and executes the mock file write.

- [ ] **4.2 CLI Power Management**
  - [ ] Implement `hilcli power <host> [on|off|cycle]`.
  - [ ] Process: CLI -> Go Backend API -> Python Agent -> Relay.
  - [ ] **Test/Milestone:** Run the full end-to-end stack: Trigger the CLI command and verify the final hardware mock script fires exactly as intended.

## Phase 5: Web Dashboard
*Goal: Provide a bird's eye visual representation of all active benches.*

- [ ] **5.1 React Foundation**
  - [ ] Initialize Next.js / TypeScript project with TailwindCSS.
  - [ ] Build a component grid representing "Active", "In-Use", and "Offline" nodes.
  - [ ] **Test/Milestone:** Navigate to `localhost:3000` inside your browser. Verify the React grid populates with the same live DB data as `hilcli list`.

- [ ] **5.2 Real-Time Updates**
  - [ ] Add WebSockets from the Go Server to the React frontend.
  - [ ] **Test/Milestone:** Keep the browser open. Stop the Python agent script. Watch the Next.js UI automatically change the node status from "Online" to "Offline" within 5 seconds without refreshing the page.
