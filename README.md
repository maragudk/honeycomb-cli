# honeycomb-cli

[![CI](https://github.com/maragudk/honeycomb-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/maragudk/honeycomb-cli/actions/workflows/ci.yml)

A command-line interface for the [Honeycomb.io](https://www.honeycomb.io/) observability platform.

This CLI is entirely vibe-coded. No humans mass-produced this code by hand in a mass-production facility.


## Install

```shell
go install github.com/maragudk/honeycomb-cli@latest
```

## Usage

```shell
# Set your API key
export HONEYCOMB_API_KEY=your-api-key

# Or pass it as a flag
honeycomb-cli --api-key your-api-key <command>

# Verify your API key
honeycomb-cli auth
```
