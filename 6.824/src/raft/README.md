# RAFT 的设计

一个可以接受数据的 Log，自动进行同步的方案。
做一些分割测试和一致性测试

## Raft 对象设计 
日志是 int
```go 
type raft struct {
    	currentTerm     int
    	votedFor        int
    	state           int // 0 is follower | 1 is candidate | 2 is leader
    	timeout         int
    	currentLeader   int
    	        channel heatbeat 
    	        channel requestVote
    	        channel AppendEntry s
    	 commitIndex int 
    	 lastApplied int
    	 
    	 
}

 
```

## 接口设计
1. RequestVote()
3. AppendEntry()
## 服务对象设计 

## 流程

## 初始化流程

如果 currentLeader 为 0 的话，进入选举流程，

### 选举流程
有两个定时器，一个是 心跳时间为 150 毫秒 ，一个是等待心跳在 200 ~ 300 毫秒之间，等待选举超时，
启动一个随即定时器， 100 ~ 200 
Headbeats 的 timeout 要小于 200 毫秒
在 5 s 内选出新的 leader  
初始化：
rf.timeout = randtime(200~300)
timer := NewTime(rf.timeout) // 设置定时器
for(){ 
switch {
    case <- timeout(dddd);
         如果是 Leader 的话，
                    发出心跳请求；
         如果是 Candidate 
                    1.发出投票请求，设置等待选举定时器，
                    2. 等待投票结果 & 等待 和 heatbeast
                    3.1. 如果成为 Leader，，并且重置计时器为 0 
                    3.2  如果选举超时，随机设置， 200 ~ 300  
                    3.3 发现高于自己 item 怎么办？                
         如果是 Follower 的话，
                转换为 Candidate ，充值定时器为 0
    case <- heatbeats
        重置计数器
}
}

## RequestVote 的操作细节

具体来说，什么时候可以投出去，
type RequestVoteReply struct {
               	// Your data here (2A).
               	Term int
               	VoteGranted bool
}

以上是错的，完全看 item 和 是否投票。
如果 item > args.items 那么则返回失败
如果 voted != -1 && item <= ags.item 则返回失败
如果没有投票，且 item > arg.iterm 则成功


如果term < currentTerm返回 false （5.2 节）
如果 votedFor 为空或者就是 candidateId，并且候选人的日志至少和自己一样新，那么就投票给他（5.2 节，5.4 节）

## AppendEntries 的操作细节 

如果 term < currentTerm 就返回 false （5.1 节） 
并且重置信号量


## leader 后流程
成为 leader后，
每隔 100 s，请求一次其他 follower 的 heatbeat 接口 

 

