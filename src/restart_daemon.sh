#!/bin/bash

# Define color output (optional, for better log readability)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color (reset)

# Logging function
log_info() {
    echo -e "${GREEN}[INFO] $(date +'%Y-%m-%d %H:%M:%S') $1${NC}"
}

log_warn() {
    echo -e "${YELLOW}[WARN] $(date +'%Y-%m-%d %H:%M:%S') $1${NC}"
}

log_error() {
    echo -e "${RED}[ERROR] $(date +'%Y-%m-%d %H:%M:%S') $1${NC}"
}

# Step 1: Find target process (extract PID)
log_info "Start searching for processes containing 'init_daemon'..."
# Extract PID (only first column, filter out grep itself)
target_pid=$(ps -eo pid,ppid,pcpu,cmd | grep -E "init_daemon" | grep -v grep | awk '{print $1}')

# Step 2: Handle target process
if [ -z "$target_pid" ]; then
    log_warn "No process containing 'init_daemon' found, skip kill operation"
else
    log_info "Found target process PID: $target_pid, executing force kill..."
    # Force kill the process (-9 = SIGKILL)
    kill -9 "$target_pid"
    # Check if kill operation succeeded
    if [ $? -eq 0 ]; then
        log_info "Process $target_pid killed successfully"
    else
        log_error "Failed to kill process $target_pid, check permissions or process status"
        # Optional: Exit script on failure (adjust based on your needs)
        # exit 1
    fi
    # Short delay to ensure process exits completely
    sleep 1
fi

# Step 3: Start mini-docker daemon
log_info "Starting mini-docker daemon..."
# Switch to script directory (ensure ./mini-docker path is correct, modify as needed)
# cd /root/root/tiny-docker/src/ || { log_error "Failed to switch directory"; exit 1; }
./mini-docker daemon
# Check if startup succeeded
if [ $? -eq 0 ]; then
    log_info "mini-docker daemon started successfully"
else
    log_error "mini-docker daemon startup failed"
    exit 1
fi

exit 0
