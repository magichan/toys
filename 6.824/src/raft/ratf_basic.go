package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"golang.org/x/text/feature/plural"
	"sync"
	"time"
)
import "labrpc"

// import "bytes"
// import "encoding/gob"



//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     struct{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}
const(
	FOLLOWER = iota
	CANDIDATE
	LEADER
)
//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	currentTerm     int // default 0
	votedFor        int // -1 is noting
	state           int // 0 is follower | 1 is candidate | 2 is leader
	timeout         int
	currentLeader   int // -1 is noting

	heartbeat chan struct{}
	cancelSelection chan struct{}
}

func (rf * Raft) ToFollower(term int, candidate int){
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.currentTerm = term
	rf.votedFor = candidate
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).

	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := gob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := gob.NewDecoder(r)
	// d.Decode(&rf.xxx)
	// d.Decode(&rf.yyy)
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
}




//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term int
	CandidateId int
	LastLogIndex int
	LastLogTerm int

}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term int
	VoteGranted bool
}
type AppendEntriesArgs struct {
	Term int
	LeaderId int

}
type AppendEntriesReply struct {
	term int
	success bool
}
//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	// 什么时候可以投，什么时候不可以投，
	// 投票后是否重置计时器。
	// 什么时候可以投
	// 		如果 votedFor 为空或者就是 candidateId，并且候选人的日志至少和自己一样新，那么就投票给他（5.2 节，5.4 节）
	//      并且发送 heartbeat
	rf.mu.Lock()
	if rf.currentTerm <= args.Term && ( rf.votedFor == -1 || rf.votedFor == args.CandidateId )  {

		reply.Term = args.Term
		reply.VoteGranted = true
		rf.heartbeat <- struct {}{}
		return
	}else if rf.currentTerm > args.Term {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}
	rf.mu.Unlock()

	reply.VoteGranted = false
	return
}

func (rf *Raft) AppendEntries(args * AppendEntriesArgs, reply * AppendEntriesReply){
	// 2A 没有日志功能
	// 如果为 HeatBeats 的功能
	//如果 term < currentTerm 就返回 false （5.1 节）
	//获取这个心跳信号，意味这投票已经结束了，所有的接受该信号的人都已经变为 Leader
	// 	将这个 Raft 转换为 Leader 的 Follower，重置投票记录，如果是 Candidate 转换的话，需要取消这个 raft 端的
	rf.mu.Lock()
	defer  rf.mu.Unlock()
	if rf.currentTerm > args.Term {
		reply.term = rf.currentTerm
		reply.success = false
	}else{
		rf.heartbeat <- struct {}{}
		rf.currentTerm = args.Term
		rf.currentLeader = args.LeaderId
		if rf.state == CANDIDATE {
			rf.cancelSelection <- struct{}{}
		}
		rf.state = FOLLOWER
		rf.votedFor = -1 // reset voted
		reply.success = true
	}

}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}


//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).


	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}
func (rf * Raft) startUp(){
	timer := time.NewTimer()


	/*
	有两个定时器，一个是 心跳时间大于 150 毫秒 ，一个是等待心跳在 200 ~ 300 毫秒之间，等待选举超时，
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
			如果是 Follower 的话，
			转换为 Candidate ，充值定时器为 0
		case <- heatbeats
			重置计数器
		}
	}
	*/
}
//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	rf.state = FOLLOWER
	rf.votedFor = -1
	rf.timeout = randInt(200,300)
	rf.heartbeat = make(chan interface{})


	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	go rf.startUp()


	return rf
}
