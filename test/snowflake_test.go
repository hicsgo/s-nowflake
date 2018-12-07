package test

import (
	"testing"
	"snowflake"
	"fmt"
	"sync"
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 注意node全局实例化一次就行,
 * 而且传入的node值不能重复(1-1023)否则产生的Id会有重复
 * 下面两组中测试节点（0，1，2）并发在（0,1，2）节点生成Id判断是否有重复(成功无重复)
 * 缺点：节点值需要配置，分布式环境中这个值需要不同，可以以读取环境变量的方式在不同服务器设置不同的值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func TestGenerateId(t *testing.T) {
	ch := make(chan snowflake.ID)
	var m sync.Map
	b := 0
	for i := 0; i < 3; i++ {

		go func(j int64) {

			count := 1000000

			n, _ := snowflake.NewNode(j)

			// 并发 count 个 goroutine 进行 snowflake ID 生成
			for i := 0; i < count; i++ {
				go func() {

					id := n.Generate()
					ch <- id

				}()

			}
		}(int64(i))

	}

	for i := 0; i < 3000000; i++ {

		id := <-ch

		b++
		// 如果 map 中存在为 id 的 key, 说明生成的 snowflake ID 有重复
		_, ok := m.Load(id)
		if ok {
			fmt.Printf("ID is not unique! %d", id)
			return
		}
		// 将 id 作为 key 存入 map
		m.Store(id, "exist")
	}
	// 成功生成 snowflake ID

	fmt.Println("success", b)
}
