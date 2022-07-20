package objecthelper

import "testing"

func TestIsStrictlyAlphNumeric(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "isStrictlyAlphNumeric 1",
			args: args{"abc123"},
			want: true,
		},
		{
			name: "isStrictlyAlphNumeric 2",
			args: args{"abc@#123"},
			want: false,
		},
		{
			name: "isStrictlyAlphNumeric 3",
			args: args{"abc#%!"},
			want: false,
		},
		{
			name: "isStrictlyAlphNumeric 4",
			args: args{")&|123"},
			want: false,
		},
		{
			name: "isStrictlyAlphNumeric 5",
			args: args{"abc"},
			want: false,
		},
		{
			name: "isStrictlyAlphNumeric 6",
			args: args{"123"},
			want: false,
		},

		{
			name: "isStrictlyAlphNumeric 7",
			args: args{"!@#$^&"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStrictlyAlphNumeric(tt.args.word); got != tt.want {
				t.Errorf("IsStrictlyAlphNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsAlphNumeric(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "containsAlphNumeric 1",
			args: args{"abc123"},
			want: true,
		},
		{
			name: "containsAlphNumeric 2",
			args: args{"abc@#$12#3"},
			want: true,
		},
		{
			name: "containsAlphNumeric 3",
			args: args{"abc@#$"},
			want: false,
		},
		{
			name: "containsAlphNumeric 4",
			args: args{"@#$123"},
			want: false,
		},
		{
			name: "containsAlphNumeric 5",
			args: args{"abc"},
			want: false,
		},
		{
			name: "containsAlphNumeric 6",
			args: args{"123"},
			want: false,
		},
		{
			name: "containsAlphNumeric 7",
			args: args{"%^&*"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAlphNumeric(tt.args.word); got != tt.want {
				t.Errorf("ContainsAlphNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStrictlyNumeric(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStrictlyNumeric(tt.args.num); got != tt.want {
				t.Errorf("IsStrictlyNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsNumeric(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsNumeric(tt.args.num); got != tt.want {
				t.Errorf("ContainsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStrictlyAlpha(t *testing.T) {
	type args struct {
		alpha string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStrictlyAlpha(tt.args.alpha); got != tt.want {
				t.Errorf("IsStrictlyAlpha() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsAlpha(t *testing.T) {
	type args struct {
		alpha string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAlpha(tt.args.alpha); got != tt.want {
				t.Errorf("ContainsAlpha() = %v, want %v", got, tt.want)
			}
		})
	}
}
