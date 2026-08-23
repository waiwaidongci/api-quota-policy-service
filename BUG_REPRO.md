# Bug Reproduction

## Bug

A stale lock release can remove the current lease. During rolling shutdown, lifecycle workers can also outlive `Stop`, leaving a policy lock held for the next process.

## Trigger

Run the targeted race-enabled lifecycle regression test against the buggy baseline:

```bash
go test -race ./internal/config -run '^TestLockLifecycleStopWaitsForWorkers$' -count=1
```

## Observed Error

```text
--- FAIL: TestLockLifecycleStopWaitsForWorkers (0.02s)
    lock_score_test.go:20: stale release removed the current lease
FAIL
FAIL github.com/example/api-quota-service/internal/config
```
