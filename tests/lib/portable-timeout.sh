#!/usr/bin/env bash
# GNU timeout is not part of macOS, so fall back to gtimeout (coreutils) and
# then to perl, which ships with macOS.
# Usage: run_with_timeout SECONDS command [args...]
# Exit code 124 means the command was killed by the timeout, as with GNU timeout.
run_with_timeout() {
    local secs="$1"
    shift
    if command -v timeout >/dev/null 2>&1; then
        timeout "$secs" "$@"
    elif command -v gtimeout >/dev/null 2>&1; then
        gtimeout "$secs" "$@"
    else
        local rc=0
        perl -e 'alarm shift; exec @ARGV' "$secs" "$@" || rc=$?
        # SIGALRM surfaces as 142; report it as 124 like GNU timeout.
        if [ "$rc" -eq 142 ]; then
            return 124
        fi
        return "$rc"
    fi
}
