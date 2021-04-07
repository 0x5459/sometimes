package vm

import (
	"reflect"
	"sometimes/ir/value"
	"testing"
)

func TestStack(t *testing.T) {

	tests := []struct {
		pushData []value.Value
		want     []value.Value
	}{
		{
			pushData: []value.Value{
				&value.Number{Val: 3},
				&value.Number{Val: 1},
				&value.Number{Val: 4},
				&value.Number{Val: 1},
				&value.Number{Val: 5},
				&value.Number{Val: 9},
			},
			want: []value.Value{
				&value.Number{Val: 9},
				&value.Number{Val: 5},
				&value.Number{Val: 1},
				&value.Number{Val: 4},
				&value.Number{Val: 1},
				&value.Number{Val: 3},
			},
		},
	}

	for _, testcase := range tests {
		s := NewOperandStack(512)
		for _, v := range testcase.pushData {
			s.Push(v)
		}
		got := make([]value.Value, 0, len(testcase.pushData))
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
	s := NewOperandStack(512)
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
	s := NewOperandStack(int(capacity))
	for ; i < capacity+1; i++ {
		s.Push(&value.Number{Val: int(i)})
	}
}
