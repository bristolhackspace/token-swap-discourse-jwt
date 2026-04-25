#!/bin/bash

CONTAINER_CMD=$(which podman || which docker)

$CONTAINER_CMD build -t localhost/token-swap-discourse-jwt:latest -f Dockerfile

