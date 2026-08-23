# Bug Reproduction

## Bug

An event queue push can return success after its context is canceled during retry backoff, so callers cannot tell that the operation was aborted.

## Trigger

Run the targeted queue cancellation regression test against the buggy baseline:

```bash
go test ./internal/repository -run '^TestQueueAbortDuringBackoff$' -count=1
```

## Observed Error

```text
--- FAIL: TestQueueAbortDuringBackoff (0.00s)
    cancellation_score_test.go:23: push error = <nil>
FAIL
FAIL github.com/example/api-quota-service/internal/repository
```
