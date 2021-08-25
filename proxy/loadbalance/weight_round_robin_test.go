package loadbalance

import (
	"fmt"
	"os"
	"testing"
)

//TestMain会在下面所有测试方法执行开始前先执行，一般用于初始化资源和执行完后释放资源
func TestMain(m *testing.M) {
	fmt.Println("初始化资源")
	result := m.Run()
	fmt.Println("释放资源")
	os.Exit(result)
}

//单元测试
func TestRandWeight(t *testing.T) {
	rw := &WeightRoundRobinBalance{}
	rw.Add("127.0.0.1:2001|10", "127.0.0.1:2002|20", "127.0.0.1:2003|30", "127.0.0.1:2004|40")

	c := make(map[string]int)
	for i := 0; i < 20; i++ {
		addr := rw.Next()
		c[addr]++
	}
	t.Logf("%#v\n", c)
}

func BenchmarkRandWeight(b *testing.B) {
	rw := &WeightRoundRobinBalance{}
	rw.Add("127.0.0.1:2001|10", "127.0.0.1:2002|20", "127.0.0.1:2003|30", "127.0.0.1:2004|40")
	c := make(map[string]int)
	for i := 0; i < 200; i++ {
		addr := rw.Next()
		c[addr]++
	}
}
