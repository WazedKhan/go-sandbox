package aggregate

import (
	"maps"
	"testing"
)

func TestAggregate(t *testing.T) {
	tests := []struct {
		name string
		in   []map[string]int
		want map[string]int
	}{
		{
			name: "interview example",
			in:   []map[string]int{{"m1": 20}, {"m2": 30}, {"m1": 20}},
			want: map[string]int{"m1": 40, "m2": 30},
		},
		{
			name: "several keys per map",
			in:   []map[string]int{{"a": 1, "b": 2}, {"b": 3, "c": 4}},
			want: map[string]int{"a": 1, "b": 5, "c": 4},
		},
		{
			name: "negative and zero values",
			in:   []map[string]int{{"x": 5}, {"x": -5}, {"y": 0}},
			want: map[string]int{"x": 0, "y": 0},
		},
		{
			name: "nil map inside the list",
			in:   []map[string]int{{"a": 1}, nil, {"a": 2}},
			want: map[string]int{"a": 3},
		},
		{
			name: "empty list returns an empty map, not nil",
			in:   []map[string]int{},
			want: map[string]int{},
		},
		{
			name: "nil list returns an empty map, not nil",
			in:   nil,
			want: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Aggregate(tt.in)
			if got == nil {
				t.Fatalf("got nil map; callers should be able to write to the result")
			}
			if !maps.Equal(got, tt.want) {
				t.Errorf("Aggregate(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestAggregateDoesNotMutateInput(t *testing.T) {
	first := map[string]int{"m1": 20}
	in := []map[string]int{first, {"m1": 20}}

	Aggregate(in)

	if first["m1"] != 20 {
		t.Errorf("input map was modified: m1 = %d, want 20", first["m1"])
	}
}

func TestResultIsWritable(t *testing.T) {
	got := Aggregate(nil)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("writing to the result panicked: %v", r)
		}
	}()
	got["new"] = 1
}
