package ssh

import (
		"fmt"
		"strconv"
		"testing"
)

func TestParsePort(t *testing.T) {
		net := `
  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0
   1: 0B00007F:AD15 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37784 1 0000000000000000 100 0 0 10 0
   2: 00000000:1F40 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37835 1 0000000000000000 100 0 0 10 0
   5: 020012AC:EA6C 9D4E7822:01BB 01 00000000:00000000 00:00000000 00000000     0        0 35787 1 0000000000000000 20 4 1 10 -1
		`
		want = []int{9000, 
		get = ParsePort(net)
		if got != want {
				t.Errorf("err")
		}
}
 
