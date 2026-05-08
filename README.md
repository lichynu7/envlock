# envlock

> Snapshot and diff environment variables across shell sessions and CI pipelines.

---

## Installation

```bash
go install github.com/yourusername/envlock@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envlock.git && cd envlock && go build -o envlock .
```

---

## Usage

**Take a snapshot of the current environment:**

```bash
envlock snapshot --output env.lock
```

**Diff two snapshots:**

```bash
envlock diff env.lock env-new.lock
```

**Compare current environment against a saved snapshot:**

```bash
envlock diff env.lock --current
```

Example output:

```
+ NEW_VAR=hello
- REMOVED_VAR=old_value
~ PATH=/usr/bin → /usr/local/bin
```

envlock is useful for detecting unexpected environment changes between CI pipeline stages, debugging session drift, or auditing environment configuration over time.

---

## License

MIT © [yourusername](https://github.com/yourusername)