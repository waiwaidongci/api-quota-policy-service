# Bug Reproduction

## Bug

The middleware chain executes in reverse order, which places the request-context boundary on the wrong side of the handler and breaks propagation through the chain.

## Trigger

Run the targeted middleware context regression test against the buggy baseline:

```bash
go test ./internal/transport -run '^TestMiddlewareContextBoundary$' -count=1
```

## Observed Error

```text
--- FAIL: TestMiddlewareContextBoundary (0.00s)
    middleware_score_test.go:41: middleware order = []string{"second", "first", "base"}
FAIL
FAIL github.com/example/api-quota-service/internal/transport
```
