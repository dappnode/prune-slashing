# Prune slashing protection on demmand

Prune a given slashing database by taking the latest 10 `signed_blocks` and `signed_attestations` for each `pubkey`

### Instructions of use

**Compile**: should be compiled with golang version 1.17

```
go build -o slashing-prune slashing-prune.go
```

The output will be an executable named `slashing-prune`

**Run**:

```
./slashing-prune --source-path <source path of the slashing protection file to be prunned> --target-path <target path to create the prunned slashing protection file>
```
