package ssh

import (
	"testing"
)

func TestCanPortForward(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"0: 00000000:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0", true},
		{"0: 00000000:270F 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0", false},
		{"0: 00000000:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     1        0 37861 1 0000000000000000 100 0 0 10 0", false},
		{"0: 20000001:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0", false},
	}
	uid := Uid("0")
	for i, test := range tests {
		got := canPortForward(test.input, uid)
		if got != test.want {
			t.Errorf("Error: got %t, but want %t in %d times", got, test.want, i)
		}
	}
}

func TestParsePort(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"270F", "9999"},
		{"AD15", "44309"},
	}
	for _, test := range tests {
		got, err := parsePort(test.input)
		if err != nil {
			t.Errorf("parsePort error: %v", err)
			continue
		}
		if got != test.want {
			t.Errorf("Error got: %s, want: %s", got, test.want)
		}
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		line   string
		port   string
		errMsg string
	}{
		{"0: 00000000:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0", "9999", ""},
		{"1: 0B000000:AD15 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37784 1 0000000000000000 100 0 0 10 0", "", "This port is not forwardable"},
	}
	uid := Uid("0")
	for _, test := range tests {
		port, err := parseLine(test.line, uid)
		if port != test.port {
			t.Errorf("Error: got %s, expected: %s", port, test.port)
		}
		if err != nil {
			errMsg := err.Error()
			if errMsg != test.errMsg {
				t.Errorf("Error: got Error Massage: %s, expected: %s", errMsg, test.errMsg)
			}
		}
	}
}

func TestFindForwardablePorts(t *testing.T) {
	lines :=
		`  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 56148 1 0000000000000000 100 0 0 10 0
   1: 00000000:0378 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 58927 1 0000000000000000 100 0 0 10 0
   2: 00000000:00EA 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 57816 1 0000000000000000 100 0 0 10 0
   3: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 57777 1 0000000000000000 100 0 0 10 0
   4: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     2        0 57777 1 0000000000000000 100 0 0 10 0
   5: 00000000:1F40 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 56126 1 0000000000000000 100 0 0 10 0
   6: 0B00007F:81A1 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 59682 1 0000000000000000 100 0 0 10 0`

	expected := []string{"9999", "888", "234", "22"}
	expectedLen := len(expected)
	uid := Uid("0")
	ports := FindForwardablePorts(lines, uid)
	portsLen := len(ports)

	if portsLen != expectedLen {
		t.Errorf("Error: port length: %d, expected: %d", portsLen, expectedLen)
		return
	}

	for i := 0; i < expectedLen; i++ {
		if ports[i] != expected[i] {
			t.Errorf("Error: got %s, expected: %s", ports[i], expected[i])
		}
	}
}
