# Phase 1
Commit hash: https://github.com/tuananh131001/go-distributed-system-6.5840-lab/commit/198f60ac29080c49a15821425d0095ceee0c5a8c
1. Define RPC message between coordinator and workerDefine Coordinator to service files to Workers
2. Define Coordinator to server files to workers
3. Define worker to receive files with its information

### QA Strategy
1. Run this command to compline the code
`go build -buildmode=plugin ../mrapps/wc.go`
2. Run these commands at the same time
` go run mrcoordinator.go pg-*.txt`
`go run mrworker.go wc.so`
3. Observe result Coordinator transfer the files to its workers

### Test Evidence:
![image](https://github.com/user-attachments/assets/ecbcbc59-fd99-48a5-a0dc-c9fe1a6c690d)
