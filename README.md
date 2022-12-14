# lynx-benchmark

- [lynx-benchmark](#lynx-benchmark)
  - [Architecture](#architecture)
  - [Development](#development)
      - [version 0.1:](#version-01)

Runner is a microservice handling code execution requests.

## Architecture

The idea is to make the service both safe and easy to tune to various needs,
so we'd like to dispatch user code into docker containers and run it in a
safe and consistent environment. Also, we want to keep it as simple as possible.

## Running in container
In order to run `lynx-benchmark` use:

    podman pull ghcr.io/group-project-gut/lynx-benchmark:0.1
    podman run lynx-benchmark:0.1 [url] [options]

## Development

#### version 0.1:

- [X] create N sessions at desired address of service
- [X] run some code on the sessions and measure latency

#### version 0.2:

- [X] parsing arguments

### version 0.3:
- [X]  Compatibility with lynx:0.3
