package vm

import (
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {

	tests := []struct {
		pushData []uint
		want     []uint
	}{
		{
			pushData: []uint{3, 1, 4, 1, 5, 9},
			want:     []uint{9, 5, 1, 4, 1, 3},
		},
	}

	for _, testcase := range tests {
		s := NewVMStack(512)
		for _, v := range testcase.pushData {
			s.Push(v)
		}
		got := make([]uint, 0, len(testcase.pushData))
		for i := 0; i < len(testcase.pushData); i++ {
			got = append(got, s.Pop())
		}
		if !reflect.DeepEqual(got, testcase.want) {
			t.Errorf("pop stack want: %v; got: %v", testcase.want, got)
		}
	}
}

func TestStack_Pop_stackoverflow(t *testing.T) {
	defer func() {
		want := StackOverflow
		got := recover()
		if want != got {
			t.Errorf("When popping an empty stack it should panic. want `%s`; got %v", want, got)
		}
	}()
	s := NewVMStack(512)
	s.Pop()
}


func TestStack_Push_stackoverflow(t *testing.T) {
	defer func() {
		want := StackOverflow
		got := recover()
		if want != got {
			t.Errorf("push elements to a full stack. want `%s`; got %v", want, got)
		}
	}()
	var i, capacity uint = 0, 32
	s := NewVMStack(int(capacity))
	for ; i < capacity+1; i++ {
		s.Push(i)
	}
}

