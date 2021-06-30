package ring_buffer

import (
	"reflect"
	"testing"
)

func Test_ringBuffer(t *testing.T) {
	type item struct {
		op  string
		val interface{}
	}
	tests := []struct {
		name  string
		size  int
		tYpe  string
		block bool
		items []item
	}{
		// TODO: Add test cases.
		//array
		{
			name: "0",
			size: 1,
			tYpe: TYPE_ARRAY,
			items: []item{
				item{"get", nil},
			},
		},
		{
			name: "2",
			size: 1,
			tYpe: TYPE_ARRAY,
			items: []item{
				item{"append", 1},
				item{"get", 1},
				item{"append", 2},
				item{"append", 3},
				item{"get", 3},
			},
		},
		{
			name: "4",
			size: 3,
			tYpe: TYPE_ARRAY,
			items: []item{
				item{"append", 1},
				item{"append", 2},
				item{"append", 3},
				item{"get", 1},
				item{"get", 2},
				item{"get", 3},
				item{"append", 4},
				item{"append", 5},
				item{"append", 6},
				item{"append", 7},
				item{"get", 5},
				item{"get", 6},
				item{"get", 7},
			},
		},

		//list
		{
			name: "1",
			size: 1,
			tYpe: TYPE_LIST,
			items: []item{
				item{"get", nil},
			},
		},
		{
			name: "3",
			size: 1,
			tYpe: TYPE_LIST,
			items: []item{
				item{"append", 1},
				item{"get", 1},
				item{"append", 2},
				item{"append", 3},
				item{"get", 3},
			},
		},
		{
			name: "5",
			size: 3,
			tYpe: TYPE_LIST,
			items: []item{
				item{"append", 1},
				item{"append", 2},
				item{"append", 3},
				item{"get", 1},
				item{"get", 2},
				item{"get", 3},
				item{"append", 4},
				item{"append", 5},
				item{"append", 6},
				item{"append", 7},
				item{"get", 5},
				item{"get", 6},
				item{"get", 7},
			},
		},
		{
			name:  "6",
			size:  3,
			tYpe:  TYPE_LIST,
			block: true,
			items: []item{
				item{"append", 1},
				item{"append", 2},
				item{"append", 3},
				item{"get", 1},
				item{"get", 2},
				item{"get", 3},
				item{"append", 1},
				item{"append", 2},
				item{"append", 3},
				item{"append", 4},
				item{"get", 2},
				item{"get", 3},
				item{"get", 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewRingBuffer(tt.size).EvictType(tt.tYpe)
			if tt.block {
				cb = cb.Block()
			}
			cache := cb.Build()
			for index, v := range tt.items {
				switch v.op {
				case "get":
					if got := cache.Get(); !reflect.DeepEqual(got, v.val) {
						t.Errorf("test ringBuffer failed. index :%d, got: %v, want: %v", index, got, v.val)
					}
				case "append":
					cache.Append(v.val)
				default:
				}
			}
		})
	}
}
