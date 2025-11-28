# Slim Archive Node Implementation Notes

## Goal

Build a **trace-capable pruned node** that supports `debug_traceTransaction` on **any historical block from genesis**, with minimal extra storage (~50-100 GB for traces only).

## Key Insight

**What gets pruned vs what doesn't:**

| Data | Pruned? | Where Stored |
|------|---------|--------------|
| Block headers | NO | Chain DB |
| Block bodies (txs) | NO | Chain DB |
| Receipts | NO | Chain DB |
| State trie | YES | Only recent in Chain DB |
| **Traces** | N/A - requires state | **Slim Archive** |

**Only traces need separate storage** - everything else survives pruning!

## Solution

Full sync from genesis with a hook that captures traces before state is pruned:

```
Block Processing (insertBlock)
         │
         ▼
   Process() + ValidateState()
         │
         ▼
   ┌─────────────────────────────────────────────┐
   │  BlockHook.OnBlockProcessed()               │
   │  - Run callTracer on each tx                │
   │  - Run prestateTracer (diffMode) on each tx │
   │  - Store both in separate PebbleDB          │
   └─────────────────────────────────────────────┘
         │
         ▼
   writeBlockAndSetHead() + Commit()
   (state may now be pruned)
```

## Files

**Created:**
- `plugin/evm/slimarchive/store.go` - PebbleDB wrapper (1MB blocks, ZSTD L1), traces only
- `plugin/evm/slimarchive/indexer.go` - Implements `core.BlockHook`, runs callTracer + prestateTracer
- `plugin/evm/slimarchive/rpc.go` - RPC API (`slimarchive_TraceTransaction`)

**Modified:**
- `plugin/evm/config/config.go` - Added `SlimArchiveConfig` with validation
- `core/blockchain.go` - Added `BlockHook` interface, call site in `insertBlock`
- `plugin/evm/vm.go` - Wires up slim archive, consistency checks
- `eth/tracers/api.go` - When enabled, uses slim archive directly (no fallback)

## Config

```json
{
  "pruning-enabled": true,
  "state-sync-enabled": false,
  "slim-archive": {
    "enabled": true,
    "path": "/data/slim-archive"
  }
}
```

## Guardrails

1. **Requires pruning=true** - archive mode stores full trie which defeats the purpose
2. **Requires state-sync=false** - state sync skips blocks, would miss history  
3. **Cannot disable once enabled** - marker in DB prevents this
4. **Cannot re-enable after disable** - would have incomplete history

## Storage Schema (PebbleDB)

| Key Pattern | Value |
|-------------|-------|
| `ct:{txHash}` | callTracer result JSON |
| `sd:{txHash}` | prestateTracer diffMode result JSON |
| `meta:head` | Last indexed block |
| `meta:enabled` | Enabled marker |

## Supported Tracers

Both tracers are run and stored for every transaction:

```bash
# Call trace (default)
curl -X POST ... -d '{"method": "debug_traceTransaction", "params": ["0x...", {"tracer": "callTracer"}]}'

# State diff
curl -X POST ... -d '{"method": "debug_traceTransaction", "params": ["0x...", {"tracer": "prestateTracer", "tracerConfig": {"diffMode": true}}]}'
```

## Testing Strategy

Run slim archive node alongside real archival node, compare `debug_traceTransaction` output byte-by-byte.

## TODO / Future Work

- Consider async indexing for better sync performance
- Metrics for indexing lag, storage size
- Support for other tracers (4byteTracer, etc.) if needed
