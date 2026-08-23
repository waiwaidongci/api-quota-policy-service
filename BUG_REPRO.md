# Bug Reproduction

## Bug

Canceling an audit export is reported as a successful export with no error, so callers cannot distinguish an aborted operation from a valid empty result.

## Trigger

Run the targeted audit export regression test against the buggy baseline:

```bash
go test ./internal/logging -run '^TestAuditExportAbortError$' -count=1
```

## Observed Error

```text
--- FAIL: TestAuditExportAbortError (0.00s)
    export_score_test.go:20: audit export error = <nil>
FAIL
FAIL github.com/example/api-quota-service/internal/logging
```
