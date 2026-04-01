import argparse
import socket
import logging
import time
import requests
import subprocess
import os

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

CENTRAL_SERVER_URL = os.environ.get("HIL_SERVER_URL", "http://localhost:8080")
CENTRAL_SSH_HOST = os.environ.get("HIL_SSH_HOST", "localhost")
CENTRAL_SSH_USER = os.environ.get("HIL_SSH_USER", os.getlogin())

# Mock Relay Plugin representing Phase 4 hardware integration
class MockRelay:
    def cycle_power(self, dut_name):
        logging.info(f"[MOCK_RELAY]: Power cycling DUT for {dut_name}...")
        time.sleep(1)
        logging.info(f"[MOCK_RELAY]: Successfully cycled power.")

def register_node(hostname):
    """Registers the node with the central Go backend and returns the assigned port."""
    url = f"{CENTRAL_SERVER_URL}/api/v1/nodes/register"
    payload = {"hostname": hostname}
    
    while True:
        try:
            logging.info(f"Attempting to register node '{hostname}' at {url}...")
            response = requests.post(url, json=payload, timeout=5)
            response.raise_for_status()
            data = response.json()
            assigned_port = data.get("assigned_ssh_port")
            logging.info(f"Successfully registered! Assigned Reverse SSH Port: {assigned_port}")
            return assigned_port
        except requests.exceptions.RequestException as e:
            logging.error(f"Failed to reach central server: {e}. Retrying in 5 seconds...")
            time.sleep(5)

def start_reverse_tunnel(assigned_port):
    """Starts the autossh reverse tunnel linking central server's port to local port 22."""
    logging.info(f"Starting autossh tunnel mapping {CENTRAL_SSH_HOST}:{assigned_port} -> localhost:22")
    
    # We use -N to not execute a remote command, just forward ports.
    # -R maps [remote_port]:[local_host]:[local_port]
    cmd = [
        "autossh", 
        "-f",       # run autossh in background instead of blocking
        "-M", "0",  # Don't use a monitoring port
        "-N",       # No command execution on remote host
        "-o", "BatchMode=yes", # Don't prompt for password! Fail if keys aren't loaded.
        "-o", "StrictHostKeyChecking=no",
        "-o", "ServerAliveInterval=30",
        "-o", "ServerAliveCountMax=3",
        "-R", f"{assigned_port}:localhost:22",
        f"{CENTRAL_SSH_USER}@{CENTRAL_SSH_HOST}"
    ]
    
    logging.info(f"Executing: {' '.join(cmd)}")
    try:
        subprocess.run(cmd, check=True)
        logging.info(f"Tunnel initiated! The Go Central Server can now securely route traffic back down through port {assigned_port}.")
    except FileNotFoundError:
        logging.error("CRITICAL: 'autossh' is not installed. Please install it (e.g., 'brew install autossh').")
        exit(1)
    except subprocess.CalledProcessError as e:
        logging.error(f"Failed to start autossh tunnel. Make sure passwordless SSH is configured for {CENTRAL_SSH_USER}@{CENTRAL_SSH_HOST}. Error: {e}")

def start_command_listener():
    """A loop mimicking Phase 4's RPC/WebSocket command listener."""
    logging.info("Agent is now running in the foreground, awaiting remote commands from the central server...")
    relay = MockRelay()
    try:
        while True:
            # We will implement real WebSockets here in Phase 4.
            time.sleep(10)
    except KeyboardInterrupt:
        logging.info("Agent shutting down.")

def main():
    parser = argparse.ArgumentParser(description="HIL Mock Node Agent")
    parser.add_argument("--name", type=str, default=socket.gethostname(), help="Mock hostname for the bench")
    args = parser.parse_args()

    # 1. Boot up and register with the Central API
    assigned_port = register_node(args.name)
    
    if assigned_port:
        # 2. Establish persistent, self-healing reverse tunnel via autossh
        start_reverse_tunnel(assigned_port)
        
        # 3. Enter main event loop (mocking Phase 4)
        start_command_listener()

if __name__ == "__main__":
    main()
