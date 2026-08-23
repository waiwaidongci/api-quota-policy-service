# Bug Reproduction

## Bug

Concurrent token-bucket decisions race while updating the shared counter map and token state. The race detector reports conflicting accesses in `MemoryCounter.Allow` and `RefillTokens`.

## Trigger

Run the targeted race-enabled token refill regression test against the buggy baseline:

```bash
go test -race ./internal/infrastructure -run '^TestCounterConcurrentTokenRefill$' -count=1
```

## Observed Error

```text
WARNING: DATA RACE
Read at 0x00c000092f60 by goroutine 9:
  runtime.mapdelete_fast64()
  github.com/example/api-quota-service/internal/infrastructure.(*MemoryCounter).Allow()
      internal/infrastructure/memory.go:129

Previous write at 0x00c000092f60 by goroutine 14:
  runtime.mapaccess2_faststr()
  github.com/example/api-quota-service/internal/infrastructure.(*MemoryCounter).Allow()
      internal/infrastructure/memory.go:132
```
