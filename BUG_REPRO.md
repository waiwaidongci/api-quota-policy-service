# Bug Reproduction

## Bug

A newly stored policy-cache entry is treated as expired. The affected cache and filtering paths also fail to preserve an isolated policy snapshot.

## Trigger

Run the targeted cache snapshot regression test against the buggy baseline:

```bash
go test ./internal/domain -run '^TestPolicyCacheSnapshotIsolation$' -count=1
```

## Observed Error

```text
--- FAIL: TestPolicyCacheSnapshotIsolation (0.00s)
    cache_score_test.go:16: fresh cache entry was treated as expired
FAIL
FAIL github.com/example/api-quota-service/internal/domain
```
