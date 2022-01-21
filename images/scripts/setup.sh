#!/bin/bash

[ -z "$REPO_UID" ] && REPO_UID=1000

useradd -m -u "$REPO_UID" builder

chown -R "$REPO_UID" /work /repo

su builder build.sh