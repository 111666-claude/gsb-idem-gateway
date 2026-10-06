# idem-gateway

幂等网关：按幂等键去重，并按保留窗口缓存结果。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/idemgateway --sample=inflight
go run ./cmd/idemgateway --sample=expiry
go run ./cmd/idemgateway --sample=empty
go run ./cmd/idemgateway --sample=overwrite
go run ./cmd/idemgateway --sample=work
```

## 口径

- **首次**：没有记录、或记录已过期时，申请返回 `new` 并建一条待处理记录，过期时刻为当前时刻加窗口。
- **处理中**：记录存在、未过期且还没完成时，再次申请返回 `inflight`。
- **重放**：记录存在、未过期且已完成时，返回 `replay`，并沿用首次完成时缓存的结果。
- **完成**：首次完成落结果；重复完成不改动已缓存的结果。

## 不变量

- 同一个幂等键在窗口内至多有一条记录，其状态只从待处理变为已完成一次。
- 结果一旦落定就不再被后续调用改写。
- `scanned` 不随请求数平方级放大：一千个键共三千次调用，`scanned` 不超过 6000。

## 输出契约

场景打印一行 JSON，键随场景而定。`--sample=work` 打印 `{"keys": N, "scanned": N}`。
