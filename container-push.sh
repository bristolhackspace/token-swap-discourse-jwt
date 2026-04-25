#!/bin/bash

CONTAINER_CMD=$(which podman || which docker)

$CONTAINER_CMD tag localhost/token-swap-discourse-jwt:latest registry.bristolhackspace.org/token-swap-discourse-jwt:latest

$CONTAINER_CMD push registry.bristolhackspace.org/token-swap-discourse-jwt:latest

