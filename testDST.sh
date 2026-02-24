#!/bin/sh

# Test output for the day before, the day when DST changes and the day after

# Exit on error
set -e

for day in 2026-03-28 2026-03-29 2026-03-30 2025-10-25 2025-10-26 2025-10-27; do
    go run . -d $day helsinki toronto paris katmandu
done
