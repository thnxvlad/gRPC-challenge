# Go Coding Challenge - J

## Testing (or getting) Knowledge

- Go-lang
- gRPC / gRPC-Streaming
- gRPC Metadata
- HTTP(REST)
- Go testing
- Docker
- Console and CLI tools

## Tasks

> Read through all the tasks and notes before solving.

### F | 12%pp

- [ ] Fork this repo
- [ ] Clone -> Now you can push your changes
- [ ] Install [buf](https://docs.buf.build/introduction) and generate code from [proto](pkg/proto/challenge.proto)
- [ ] Make gRPC Server base: `main.go` in `cmd/server`, server methods in `pkg/server`
- [ ] Comments and docs are appreciated :)

### D | 3%pp

- [ ] Read environment vars: [task description](#environment)
- [ ] Make Metadata reader method: [task description](#metadata)
- [ ] Make link shortener method: [task description](#short-links)

### C | 3%pp

- [ ] Make link metadata test(s): [task description](#metadata-bonus)
- [ ] Make Timer Streaming method: [task description](#online-timer-and-streaming)
- [ ] Make it running in Docker: [task description](#docker)

### B | 5%pp

- [ ] Make link shortener test(s): [task description](#short-links-bonus)
- [ ] Make `cmd/client` using [cobra](https://github.com/spf13/cobra)

### A | 5%pp

- [ ] Make Timer Streaming test: [task description](#online-timer-and-streaming-bonus)
- [ ] Gituhub - Publish, test and CI: [task description](#github)

## Important Notes to this tasks

1. Along with `cmd/server` - `cmd/client` has to be provided. This would be used for manual endpoints testing.
2. Protobuf messages structure(s) can be edited in anyway you'd need(for example to match external APIs responses)
3. All dependencies should be logged and provided in `go.mod`(see `go mod tidy`)

## Environment

`main.go` must read environment variables using [viper](https://github.com/spf13/viper) and store it accesible by server scope. As the variable to read, use `BITLY_OAUTH_LOGIN` and `BITLY_OAUTH_TOKEN`, you'll need them in [the next task](#short-links)

## Metadata

Each gRPC call has context. And Outgoing and Incoming context can be appended to it.

Implement `ChallengeService.ReadMetadata` method. It has to read `i-am-random-key` from context metadata and return it in the Response data as string.

### Metadata Bonus

Implement server method test using `*testing.T`. Test should make request with random string in metadata `i-am-random-key` key and ensure it's being returned so.

## Short-links

Using [Bitly API](https://dev.bitly.com), implement `ChallengeService.MakeShortLink` method. It has to receive long link and return short bit.ly link.

API calls has to be implemented with `net/http`.

> You need to register at bit.ly to get token and username

### Short-links Bonus

Implement server method test using `*testing.T` and `net/http`. Test should generate link using server method and ensure short-link actually leads to given link(see http statuses and follow-redirects).

Additional tests covering bit.ly API edge scenarious are welcomed.

> Note:
> Tests must be checking methods over network and read server location from environment

## Online Timer and Streaming

Using [this API](https://alestic.com/2015/07/timercheck-scheduled-events-monitoring/) or simillar implement streaming endpoint `ChallengeService.StartTimer`.

Endpoint accepts message with amount of Seconds timer should run, frequency of refresh requests in seconds. Endpoint should do refresh every `freqency` seconds and send amount of seconds left back to stream and name for the new Timer.

Every new client trying to start a timer with the name of existing timer should be automatically subscribed to existing one.

### Online Timer and Streaming Bonus

Implement server method test using `*testing.T`.

Test should cover:

 1. Client reconnect - client must be resubscribed to the same timer
 2. Additional client - new client must be subscribed to the same timer as the old client and receive same (amount of) messages

> Note:
> Tests must be checking methods over network and read server location from environment

## Docker

Write `Dockerfile` for this project.

Requirements:

    1. Each step of preparation and build has to be happening in a Docker-runtime
    2. Image has to be as lightweight as possible

## Github

Publish this project to Github repo.

Set Github Actions Workflow called `ci.yml`.
Workflow must:

1. Build Docker Image
2. Start Docker Container
3. Run tests with this container
4. Publish Docker Image to `ghcr.io`(Github Packages of your repo)

> You should store bit.ly credentials in Repo secrets
