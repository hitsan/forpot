package ssh

import (
	"testing"
)

func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		input			string
		expected  bool
	} {
		{"00000000", true},
		{"020012AC", false},
	}
	for _, test := range tests {
		got := isLocalhost(test.input)
		if got != test.expected {
			t.Errorf("Error: expect %v, but %v in %s", test.expected, got, test.input)
		}
	}
}

func TestCanListen(t *testing.T) {
	statuses := []string{"0A", "01"}
	wants := []bool{true, false}
	for i, status := range statuses {
		got := canListen(status)
		want := wants[i]
		if got != want {
			t.Errorf("Error: get %t, but want %t", got, want)
		}
	}
}

func TestEqualsUid(t *testing.T) {
	testUid := "0"
	uids := []string{"2", "0"}
	wants := []bool{false, true}
	for i, uid := range uids {
		got := equalsUid(testUid, uid)
		want := wants[i]
		if got != want {
			t.Errorf("Error: test UID %s, but uid %s", testUid, uid)
		}
	}
}

func TestCanPortForward(t *testing.T) {
	tests := []struct {
		input    string
		want bool
	}{
		{"0: 00000000:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0", true},
		{"0: 00000000:270F 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0", false},
		{"0: 00000000:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     1        0 37861 1 0000000000000000 100 0 0 10 0", false},
		{"0: 20000001:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0", false},
	}
	for i, test := range tests {
		got := canPortForward(test.input)
		if got != test.want {
			t.Errorf("Error: got %t, but want %t in %d times", got, test.want, i)
		}
	}
}

	//func TestParsePort(t *testing.T) {
	//		net := `
	//   sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
	//   0: 00000000:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0
	//   1: 0B00007F:AD15 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37784 1 0000000000000000 100 0 0 10 0
	//   2: 00000000:1F40 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37835 1 0000000000000000 100 0 0 10 0
	//   5: 020012AC:EA6C 9D4E7822:01BB 01 00000000:00000000 00:00000000 00000000     0        0 35787 1 0000000000000000 20 4 1 10 -1
	//		`
	//		want := []int{9000, 8000}
	//		got := ParsePort(net)
	//		for i, v := range want {
	//				if got[i] != v {
	//						t.Errorf("Difference got: %d, want: %d", got[i], v)
	//				}
	//		}
	//}
