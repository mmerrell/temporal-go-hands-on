# SailPoint Dev Day — Hands-On Labs (Go SDK)

Two labs using a warehouse fulfillment domain. Each lab builds on the previous.

| Lab | Topic | Time |
|---|---|---|
| [Lab 1](./lab1-child-workflows/) | Child workflows | ~10 min |
| [Lab 2](./lab2-local-and-parallel/) | Parallel activities + local activities | ~15 min |

## Prerequisites

- Go 1.22+
- Temporal CLI: `brew install temporal` (or download from temporal.io)
- A running local Temporal server: `temporal server start-dev`

## Structure

Each lab has a `practice/` directory (with TODO comments) and a `solution/`
directory. Work in `practice/`. Check `solution/` only if you're stuck.

```
lab1-child-workflows/
├── practice/         ← work here
├── solution/         ← reference
└── README.md         ← instructions

lab2-local-and-parallel/
├── practice/         ← work here
├── solution/         ← reference
└── README.md         ← instructions
```

## The domain

An order fulfillment pipeline:

1. **Validate + fraud check** — fast in-process checks (local activities in Lab 2)
2. **Reserve inventory** — check warehouses for stock (child workflow; parallel in Lab 2)
3. **Process payment** — charge the customer
4. **Dispatch** — hand off to the fulfillment center

This mirrors patterns you'll find in SailPoint's identity provisioning workflows:
parallel checks across systems, isolated sub-workflows with their own retry
boundaries, and cheap in-process validation before expensive external calls.
