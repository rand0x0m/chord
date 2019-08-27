package chord

import (
	"testing"
)

import (
	"net"
)

func Test_hashData(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{s: "TestValue"}, want: "1412516341862220528149878343959718123327499787032"},
		{args: args{s: "TestValue87"}, want: "202751359931511688565338992679126862383343664194"},
		{args: args{s: "TestValue90"}, want: "606622844616454782750912665649096917170760627300"},
		{args: args{s: "TestValue20"}, want: "702343121414328359148025766316339326501878348599"},
		{args: args{s: "TestValue6"}, want: "598303051810575104682413810545506446060435977830"},
		{args: args{s: "TestValue65"}, want: "396778959088901340448484647779371769610863113696"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashData(tt.args.s); got != tt.want {
				t.Errorf("hashData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashAddr(t *testing.T) {
	type args struct {
		addr net.TCPAddr
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{addr: net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5480}}, want: "134852913187402694000634224619367376407810526135"},
		{args: args{addr: net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5487}}, want: "528390268536630496616574877981867029929644276646"},
		{args: args{addr: net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5488}}, want: "1144168734263290924232640306221451784150775350879"},
		{args: args{addr: net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5483}}, want: "739126638082231409668787470808025286961754925349"},
		{args: args{addr: net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5484}}, want: "61284563622561060941371771439854551201444769279"},
		{args: args{addr: net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5485}}, want: "50967384523505140176522743257572828238404511884"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashAddr(tt.args.addr); got != tt.want {
				t.Errorf("hashAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_exclusiveBetween(t *testing.T) {
	type args struct {
		left  string
		i     string
		right string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{"a", "a", "a"}, want: false},
		{args: args{"aaa", "aaa", "aaa"}, want: false},
		{args: args{"a", "b", "c"}, want: true},
		{args: args{"aaa", "aab", "aac"}, want: true},
		{args: args{"b", "a", "c"}, want: false},
		{args: args{"aab", "aaa", "aac"}, want: false},
		{args: args{"c", "b", "a"}, want: false},
		{args: args{"aac", "aab", "aaa"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exclusiveBetween(tt.args.left, tt.args.i, tt.args.right); got != tt.want {
				t.Errorf("exclusiveBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inclusiveBetween(t *testing.T) {
	type args struct {
		left  string
		i     string
		right string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{"a", "a", "a"}, want: true},
		{args: args{"aaa", "aaa", "aaa"}, want: true},
		{args: args{"a", "b", "c"}, want: true},
		{args: args{"aaa", "aab", "aac"}, want: true},
		{args: args{"b", "a", "c"}, want: false},
		{args: args{"aab", "aaa", "aac"}, want: false},
		{args: args{"c", "b", "a"}, want: false},
		{args: args{"aac", "aab", "aaa"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inclusiveBetween(tt.args.left, tt.args.i, tt.args.right); got != tt.want {
				t.Errorf("inclusiveBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inclusiveRightBetween(t *testing.T) {
	type args struct {
		left  string
		i     string
		right string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{"a", "a", "a"}, want: false},
		{args: args{"aaa", "aaa", "aaa"}, want: false},
		{args: args{"a", "b", "c"}, want: true},
		{args: args{"aaa", "aab", "aac"}, want: true},
		{args: args{"b", "a", "c"}, want: false},
		{args: args{"aab", "aaa", "aac"}, want: false},
		{args: args{"c", "b", "a"}, want: false},
		{args: args{"aac", "aab", "aaa"}, want: false},
		{args: args{"a", "b", "b"}, want: true},
		{args: args{"aaa", "aab", "aab"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inclusiveRightBetween(tt.args.left, tt.args.i, tt.args.right); got != tt.want {
				t.Errorf("inclusiveRightBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}
