# logstep

Structured log tailer with filtering, field extraction, and output formatting for JSON log streams.

---

## Installation

```bash
go install github.com/yourusername/logstep@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logstep.git && cd logstep && go build -o logstep .
```

---

## Usage

Tail a JSON log file and filter by field value:

```bash
logstep tail app.log --filter 'level=error'
```

Extract specific fields and format output:

```bash
logstep tail app.log --fields time,level,msg --format pretty
```

Pipe from stdin:

```bash
kubectl logs -f my-pod | logstep --filter 'status>=500' --fields time,method,path,status
```

### Flags

| Flag | Description |
|------|-------------|
| `--filter` | Filter expression (e.g. `level=error`, `status>=500`) |
| `--fields` | Comma-separated list of fields to display |
| `--format` | Output format: `pretty`, `json`, or `compact` (default: `pretty`) |
| `--follow` | Continue watching file for new lines |

---

## Example Output

```
2024-01-15T10:23:01Z  ERROR  failed to connect to database  host=db.internal retries=3
2024-01-15T10:23:05Z  ERROR  request timeout                path=/api/users duration=30s
```

---

## License

MIT © [yourusername](https://github.com/yourusername)