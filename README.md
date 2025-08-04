# Go State Machine with Gorm Persistence

A flexible **state machine** implementation in Go where each state runs a callback function that returns a status (`success`, `failed`, or `unknown`), and transitions to the next state are decided based on that status. The machine supports:

- **Retrying** the current state before transitioning.
- **Persisting state and payload** in PostgreSQL using [GORM](https://gorm.io/).
- Payloads as **arbitrary Go structs** serialized as JSON.
- **Thread-safe execution** using Postgres row-level locking.
- Explicit **end state support** to mark job completion.

---

## Features

- Register states with their callback functions and status-based transitions.
- Callbacks receive strongly-typed payload structs and can modify them.
- Retry logic with configurable max retries.
- State and payload saved in PostgreSQL JSONB column.
- Concurrent-safe execution via DB transactions with `SELECT FOR UPDATE`.
- Supports terminating the state machine flow gracefully.

---

## Installation

```bash
go get github.com/inamvar/state-go
```

