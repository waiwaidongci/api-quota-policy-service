# Bug Reproduction

## Bug

Creating a service through the HTTP endpoint with a missing request body can dereference a nil body. The recovery middleware turns the panic into an HTTP 500 response.

## Trigger

Run the targeted HTTP regression test against the buggy baseline:

```bash
go test ./internal/http -run '^TestServiceDefaultsDoNotPanic$' -count=1
```

## Observed Error

```text
ERROR panic value="runtime error: invalid memory address or nil pointer dereference"
--- FAIL: TestServiceDefaultsDoNotPanic (0.00s)
    request_score_test.go:23: status = 500, body={"error":"Internal Server Error","message":"internal server error"}
FAIL
```
