# Bug Reproduction

## Bug

A missing quota policy loses its not-found classification while the error crosses package boundaries. The adapter reports an internal-server status instead of the expected not-found status.

## Trigger

Run the targeted error-chain regression test against the buggy baseline:

```bash
go test ./internal/adapter -run '^TestErrorChainPreservesNotFound$' -count=1
```

## Observed Error

```text
--- FAIL: TestErrorChainPreservesNotFound (0.00s)
    error_score_test.go:13: not found status = 500
FAIL
FAIL github.com/example/api-quota-service/internal/adapter
```
