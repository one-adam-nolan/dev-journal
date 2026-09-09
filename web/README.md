# Web

Placeholder for the dev-journal web client.

The web app will consume the shared protobuf contract in `proto/dj/v1/` via generated TypeScript Connect clients (future `gen/ts/` output from Buf).

## Planned stack

- Connect-Web client generated from `proto/dj/v1/journal.proto`
- Backend served from `backend/cmd/server`
- Shared API versioning via Buf breaking-change checks

## Generate TypeScript clients (future)

When the web app is started, add a TypeScript plugin to `buf.gen.yaml` and generate into `gen/ts/`.
