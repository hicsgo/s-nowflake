package snowflake

import (
	"time"
	"sync"
	"errors"
	"math/rand"
)

const (
	nodeBits  uint8 = 10
	stepBits  uint8 = 12
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	stepMax   int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

// 起始时间戳 (毫秒数显示)
var Epoch int64 = 1288834974657 // timestamp 2006-03-21:20:50:14 GMT

// ID 结构
type ID int64

// 存储基础信息的 Node 结构
type Node struct {
	mu        sync.Mutex // 保证并发安全
	timestamp int64
	node      int64
	step      int64
}

// 生成、返回唯一 snowflake ID
func (n *Node) Generate() ID {

	//随机设置序列号位,预防跨毫秒位每次都取0,Hash取模入库数据分布不均匀
	rand.Seed(time.Now().UnixNano())
	randStep := rand.Int63n(9)

	n.mu.Lock()         // 保证并发安全, 加锁
	defer n.mu.Unlock() // 解锁

	// 获取当前时间的时间戳 (毫秒数显示)
	now := time.Now().UnixNano() / 1e6

	if n.timestamp == now {
		// step 步进 1
		n.step ++

		// 当前 step 用完
		if n.step > stepMax {
			// 等待本毫秒结束
			for now <= n.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}

	} else {
		// 本毫秒内 step 用完
		n.step = randStep
	}

	n.timestamp = now

	result := ID((now-Epoch)<<timeShift | (n.node << nodeShift) | (n.step))

	return result
}

func NewNode(node int64) (*Node, error) {
	// 如果超出节点的最大范围，产生一个 error
	if node < 0 || node > nodeMax {
		return nil, errors.New("Node number must be between 0 and 1023")
	}
	// 生成并返回节点实例的指针
	return &Node{
		timestamp: 0,
		node:      node,
		step:      0,
	}, nil
}
