import argparse
import socket
import logging
import time
import requests
import subprocess
import os
import getpass

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# Core Central Server Configuration
CENTRAL_SERVER_URL = os.environ.get("HIL_SERVER_URL", "http://localhost:8080")
CENTRAL_SSH_HOST = os.environ.get("HIL_SSH_HOST", "localhost")

# Note for Local Testing vs Production:
# Locally, getpass.getuser() allows testing reverse tunneling natively without setting up new users.
# In production, this will be explicitly configured (e.g. HIL_SSH_USER=hiluser).
CENTRAL_SSH_USER = os.environ.get("HIL_SSH_USER", getpass.getuser())

def register_node(hostname):
    """Registers the physical node with the central Go backend and returns the assigned port."""
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
        logging.info(f"SSH Tunnel initiated! The Go Central Server can now securely route traffic back down through port {assigned_port}.")
    except FileNotFoundError:
        logging.error("CRITICAL: 'autossh' is not installed. Please install it (e.g., 'brew install autossh' or 'apt install autossh').")
        exit(1)
    except subprocess.CalledProcessError as e:
        logging.error(f"Failed to start autossh tunnel. Make sure passwordless SSH is configured for {CENTRAL_SSH_USER}@{CENTRAL_SSH_HOST}. Error: {e}")

def keep_alive():
    """Maintains the python process. Future phases will implement Power Control WebSockets here."""
    logging.info("Agent is fully operational! Remote SSH access is active.")
    try:
        while True:
            # Keeps the daemon alive so systemd doesn't restart it
            time.sleep(60)
    except KeyboardInterrupt:
        logging.info("Agent shutting down.")

def main():
    parser = argparse.ArgumentParser(description="HIL Physical Node Agent")
    parser.add_argument("--name", type=str, default=socket.gethostname(), help="Hostname for the bench")
    args = parser.parse_args()

    # 1. Boot up and register with the Central API
    assigned_port = register_node(args.name)
    
    if assigned_port:
        # 2. Establish persistent, self-healing reverse tunnel via autossh
        start_reverse_tunnel(assigned_port)
        
        # 3. Enter main event loop (Future: Power control listening)
        keep_alive()

if __name__ == "__main__":
    main()
