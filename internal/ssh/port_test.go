package ssh

import (
	"testing"
)

func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
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
	tests := []struct {
		input    string
		expected bool
	}{
		{"0A", true},
		{"01", false},
	}
	for _, test := range tests {
		got := canListen(test.input)
		if got != test.expected {
			t.Errorf("Error: expect %t with %s, but %t", test.expected, test.input, got)
		}
	}
}

func TestEqualsUid(t *testing.T) {
	uid := "0"
	tests := []struct {
		input    string
		expected bool
	}{
		{"2", false},
		{"0", true},
	}
	for _, test := range tests {
		got := equalsUid(uid, test.input)
		if got != test.expected {
			t.Errorf("Error: expect %t with uid %s, but %t", test.expected, test.input, got)
		}
	}
}

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
	for i, test := range tests {
		got := canPortForward(test.input, "0")
		if got != test.want {
			t.Errorf("Error: got %t, but want %t in %d times", got, test.want, i)
		}
	}
}

func TestParsePort(t *testing.T) {
	tests := []struct {
		input string
		want string
	} {
		{"270F", "9999"}, 
		{"AD15", "44309"}, 
	}
	for _, test := range tests {
		got := parsePort(test.input)
		if got != test.want {
			t.Errorf("Error got: %s, want: %s", got, test.want)
		}
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		line string
		port string
		errMsg string
	} {
		{"0: 00000000:270F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37861 1 0000000000000000 100 0 0 10 0", "9999", ""},
		{"1: 0B000000:AD15 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 37784 1 0000000000000000 100 0 0 10 0", "", "This port is not forwardable"},
	}
	for _, test := range tests {
		port, err := parseLine(test.line, "0")
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

