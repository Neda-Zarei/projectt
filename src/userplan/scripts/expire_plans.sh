#!/bin/bash

# script to expire user plans daily
# this should be run via cron at 00:00 every day

set -e

SERVICE_NAME="userplan"
LOG_FILE="/var/log/userplan/expire_plans.log"
PID_FILE="/var/run/userplan/expire_plans.pid"
USERPLAN_BINARY="/usr/local/bin/userplan"  # adjust path ?

#creating log directory if it doesn't exist
mkdir -p "$(dirname "$LOG_FILE")"

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

#if already running
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        log "Process already running with PID $PID"
        exit 1
    else
        log "Removing stale PID file"
        rm -f "$PID_FILE"
    fi
fi

echo $$ > "$PID_FILE"

cleanup() {
    rm -f "$PID_FILE"
    log "Expiration check completed"
}

trap cleanup EXIT

log "Starting plan expiration check"

#  the expiration command
if [ -f "$USERPLAN_BINARY" ]; then
    log "Running plan expiration via CLI"
    if "$USERPLAN_BINARY" -expire-plans; then
        log "Plan expiration completed successfully"
    else
        log "Plan expiration failed"
        exit 1
    fi
else
    log "UserPlan binary not found at $USERPLAN_BINARY"
    log "Falling back to database function call"
    # fallback to direct database call
    psql -d "$DB_NAME" -c "SELECT expire_user_plans();" || {
        log "Database function call failed"
        exit 1
    }
fi

log "Plan expiration check finished"
