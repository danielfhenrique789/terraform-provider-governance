#!/bin/bash

case "$1" in
  set)
    export TF_CLI_CONFIG_FILE="$(pwd)/.terraformrc"
    echo "TF_CLI_CONFIG_FILE set to: $TF_CLI_CONFIG_FILE"
    ;;
  unset)
    unset TF_CLI_CONFIG_FILE
    echo "TF_CLI_CONFIG_FILE unset"
    ;;
  *)
    echo "Usage: source ./dev-env.sh {set|unset}"
    return 1 2>/dev/null || exit 1
    ;;
esac