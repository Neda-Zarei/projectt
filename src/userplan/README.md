# User-Plan Service

## Table of Contents

- [How to Setup](#setup) – Installation, configuration, and environment setup.
- [How to Develop](#develop) – Development workflow, running locally, testing.
- [Project Structure](#project-structure) – Directory layout and architecture overview.

## Setup

> back to [outline](#table-of-contents)

writing...

## Develop

> back to [outline](#table-of-contents)

writing...

## Project Structure

> back to [outline](#table-of-contents)

- [Top-Level](#top-level) – Root files like `main.go`, `go.mod`, configs, scripts.
- [Directories](#directories) – Organized project code.
  - [app](#app) – Application bootstrap and core wiring.
  - [build](#build) – Dockerfiles and build-related artifacts.
  - [cmd](#cmd) – Entrypoint commands (main binaries).
  - [config](#config) – Configuration structs, readers, and tests.
  - [internal](#internal) – Private application code (domain, services, adapters).
    - [adapter/repository](#adapterrepository) – Database repository implementations.
    - [api](#api) – Transport + generated protobufs.
      - [handlers/grpc](#handlersgrpc) – gRPC service endpoints.
      - [pb](#pb) – Generated protobuf code.
    - [common](#common) – Shared domain primitives/utilities.
    - [plan](#plan) – Plan domain, ports, and service logic.
    - [user](#user) – User domain, ports, and service logic.
  - [pkg](#pkg) – Reusable helper libraries (e.g. logger).
- [Conventions](#conventions) – Architectural and code organization principles.
- [Extending](#extending) – How to add new domains, transports, or infra.

### Top-Level

> back to [outline](#project-structure)

- **main.go** – Application entrypoint.
- **example.env** – Example env vars for local/dev setup.
- **go.mod / go.sum** – Go module definitions.
- **README.md** – Project overview & usage.
- **genproto.sh** – Script to regenerate gRPC/protobuf files.

### Directories

> back to [outline](#project-structure)

### app

Application bootstrap (constructs core application, wires dependencies).

### build

Artifacts for building and deployment.

- Dockerfile – Production image.
- Dockerfile.test – Test/development image.

### cmd

Entrypoint commands.

- Typically one file per executable (e.g., CLI or server binaries).

### config

Configuration logic.

- config.go – Configuration struct.
- read.go – Load config (env/files).
- read_test.go – Unit tests for config.

### **internal**:

---

Private application code (not imported outside this module).

### adapter/repository

Persistence layer (DB repositories).

- plan_repo.go, user_repo.go, etc.

### api

- **handlers/grpc** – gRPC transport layer.
- **pb** – Generated protobuf files.

### common

Shared domain primitives.

### plan

- **domain** – Plan entities.
- **port** – Interfaces (ports).
- **service.go** – Business logic.

### user

- **domain** – User entities.
- **port** – Interfaces (ports).
- **service.go** – Business logic.

---

### pkg

Public helper libraries (safe to import elsewhere).

- **logger** – Zap-based logging wrapper.

---

### Conventions

> back to [outline](#project-structure)

- **Hexagonal architecture (Ports & Adapters)**
  - domain – Pure business entities.
  - port – Interfaces for boundaries.
  - service.go – Use-case logic.
  - adapter/... – Infrastructure implementations.
- **Generated Code** → internal/api/pb.

---

### Extending

> back to [outline](#project-structure)

- Add a new domain: create internal/\<domain\>/.
- Add a new transport: extend internal/api/handlers/.
- Add new infra: extend internal/adapter/repository/.
