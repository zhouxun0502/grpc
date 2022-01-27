# Day02 流式RPC

### 为什么不用 Simple RPC

1、流式为什么要存在呢，是 Simple RPC 有什么问题吗？通过模拟业务场景，可得知在使用 Simple RPC 时，有如下问题：

2、数据包过大造成的瞬时压力

3、接收数据包时，需要所有数据包都接受成功且正确后，才能够回调响应，进行业务处理（无法客户端边发送，服务端边处理）

### 为什么用 Streaming RPC

1、大规模数据包

2、实时场景

### 模拟场景 

1、每天早上 6 点，都有一批百万级别的数据集要同从 A 同步到 B，在同步的时候，会做一系列操作（归档、数据分析、画像、日志等）。这一次性涉及的数据量确实大

2、在同步完成后，也有人马上会去查阅数据，为了新的一天筹备。也符合实时性。

**两者相较下，这个场景下更适合使用 Streaming RPC**

## cd day02/proto && protoc --go_out=. --go-grpc_out=. *.proto

## cd day02/server && go run server.go

## cd day02/client && go run client.go
