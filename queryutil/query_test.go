package queryutil_test

import (
	"testing"

	"github.com/shoraid/stx-go-utils/queryutil"

	"github.com/stretchr/testify/assert"
)

func TestQueryUtil_CalculatePagination(t *testing.T) {
	tests := []struct {
		name           string
		page           string
		perPage        string
		defaultPerPage int
		wantPage       int
		wantPerPage    int
		wantOffset     int
	}{
		{
			name:           "Valid inputs",
			page:           "2",
			perPage:        "10",
			defaultPerPage: 15,
			wantPage:       2,
			wantPerPage:    10,
			wantOffset:     10,
		},
		{
			name:           "Default perPage when zero",
			page:           "1",
			perPage:        "0",
			defaultPerPage: 15,
			wantPage:       1,
			wantPerPage:    15,
			wantOffset:     0,
		},
		{
			name:           "Default page when zero",
			page:           "0",
			perPage:        "20",
			defaultPerPage: 15,
			wantPage:       1,
			wantPerPage:    20,
			wantOffset:     0,
		},
		{
			name:           "Invalid page input",
			page:           "abc",
			perPage:        "10",
			defaultPerPage: 15,
			wantPage:       1,
			wantPerPage:    10,
			wantOffset:     0,
		},
		{
			name:           "Invalid perPage input",
			page:           "3",
			perPage:        "xyz",
			defaultPerPage: 15,
			wantPage:       3,
			wantPerPage:    15,
			wantOffset:     30,
		},
		{
			name:           "Negative page input",
			page:           "-1",
			perPage:        "5",
			defaultPerPage: 15,
			wantPage:       1,
			wantPerPage:    5,
			wantOffset:     0,
		},
		{
			name:           "Negative perPage input",
			page:           "2",
			perPage:        "-5",
			defaultPerPage: 15,
			wantPage:       2,
			wantPerPage:    15,
			wantOffset:     15,
		},
		{
			name:           "Default page and perPage when both invalid",
			page:           "0",
			perPage:        "0",
			defaultPerPage: 15,
			wantPage:       1,
			wantPerPage:    15,
			wantOffset:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPerPage, gotOffset := queryutil.CalculatePagination(tt.page, tt.perPage, tt.defaultPerPage)

			assert.Equal(t, tt.wantPage, gotPage, "Page mismatch")
			assert.Equal(t, tt.wantPerPage, gotPerPage, "PerPage mismatch")
			assert.Equal(t, tt.wantOffset, gotOffset, "Offset mismatch")
		})
	}
}

func TestQueryUtil_CalculateTotalPage(t *testing.T) {
	tests := []struct {
		name       string
		totalData  int
		perPage    int
		wantResult int
	}{
		{
			name:       "Perfect division",
			totalData:  100,
			perPage:    10,
			wantResult: 10,
		},
		{
			name:       "Rounding up with remainder",
			totalData:  101,
			perPage:    10,
			wantResult: 11,
		},
		{
			name:       "Single page",
			totalData:  5,
			perPage:    10,
			wantResult: 1,
		},
		{
			name:       "No data",
			totalData:  0,
			perPage:    10,
			wantResult: 0,
		},
		{
			name:       "Zero perPage",
			totalData:  100,
			perPage:    0,
			wantResult: 0, // Handle invalid perPage scenario
		},
		{
			name:       "Negative perPage",
			totalData:  100,
			perPage:    -10,
			wantResult: 0, // Handle invalid perPage scenario
		},
		{
			name:       "Negative totalData",
			totalData:  -100,
			perPage:    10,
			wantResult: 0, // Handle invalid totalData scenario
		},
		{
			name:       "Large numbers",
			totalData:  1000000,
			perPage:    1000,
			wantResult: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := queryutil.CalculateTotalPage(tt.totalData, tt.perPage)

			assert.Equal(t, tt.wantResult, gotResult, "Mismatch in total pages calculation")
		})
	}
}

func TestQueryUtil_ResolveAllowedFields(t *testing.T) {
	type args struct {
		input   string
		allowed map[string]any
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "allow exact fields (bool)",
			args: args{
				input: "name,email",
				allowed: map[string]any{
					"name":  true,
					"email": true,
				},
			},
			want: []string{"name", "email"},
		},
		{
			name: "alias field (string value in map)",
			args: args{
				input: "email",
				allowed: map[string]any{
					"email": "user_email",
				},
			},
			want: []string{"user_email"},
		},
		{
			name: "mixed allowed and disallowed fields",
			args: args{
				input: "name,age,address",
				allowed: map[string]any{
					"name": true,
					"age":  "user_age",
				},
			},
			want: []string{"name", "user_age"},
		},
		{
			name: "empty input",
			args: args{
				input: "",
				allowed: map[string]any{
					"name": true,
				},
			},
			want: []string{},
		},
		{
			name: "no allowed fields match",
			args: args{
				input: "unknown,field",
				allowed: map[string]any{
					"name": true,
				},
			},
			want: []string{},
		},
		{
			name: "field with whitespace",
			args: args{
				input: "  name ,  email ",
				allowed: map[string]any{
					"name":  true,
					"email": "user_email",
				},
			},
			want: []string{"name", "user_email"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := queryutil.ResolveAllowedFields(tt.args.input, tt.args.allowed)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestQueryUtil_ResolveSingleField(t *testing.T) {
	type testCase struct {
		name         string
		input        string
		allowed      map[string]any
		defaultField string
		expected     string
	}

	cases := []testCase{
		{
			name:         "Valid field with boolean true",
			input:        "username",
			allowed:      map[string]any{"username": true},
			defaultField: "default",
			expected:     "username",
		},
		{
			name:         "Valid field with alias",
			input:        "email",
			allowed:      map[string]any{"email": "user_email"},
			defaultField: "default",
			expected:     "user_email",
		},
		{
			name:         "Invalid field should return default",
			input:        "invalid",
			allowed:      map[string]any{"email": true},
			defaultField: "default",
			expected:     "default",
		},
		{
			name:         "Empty input should return default",
			input:        "",
			allowed:      map[string]any{"email": true},
			defaultField: "default",
			expected:     "default",
		},
		{
			name:         "Field with whitespace",
			input:        "  email  ",
			allowed:      map[string]any{"email": "user_email"},
			defaultField: "default",
			expected:     "user_email",
		},
		{
			name:         "Field with false boolean should return default",
			input:        "username",
			allowed:      map[string]any{"username": false},
			defaultField: "default",
			expected:     "default",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := queryutil.ResolveSingleField(tc.input, tc.defaultField, tc.allowed)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func BenchmarkQueryUtil_CalculatePagination(b *testing.B) {
	page := "5"
	perPage := "10"
	defaultPerPage := 15

	for i := 0; i < b.N; i++ {
		queryutil.CalculatePagination(page, perPage, defaultPerPage)
	}
}

func BenchmarkQueryUtil_CalculateTotalPage(b *testing.B) {
	totalData := 1000
	perPage := 20

	for i := 0; i < b.N; i++ {
		queryutil.CalculateTotalPage(totalData, perPage)
	}
}

func BenchmarkQueryUtil_ResolveAllowedFields(b *testing.B) {
	type benchCase struct {
		name    string
		input   string
		allowed map[string]any
	}

	cases := []benchCase{
		{
			name: "10-bool",
			input: "field0,field1,field2,field3,field4," +
				"field5,field6,field7,field8,field9",
			allowed: map[string]any{
				"field0": true,
				"field1": true,
				"field2": true,
				"field3": true,
				"field4": true,
				"field5": true,
				"field6": true,
				"field7": true,
				"field8": true,
				"field9": true,
			},
		},
		{
			name: "10-alias",
			input: "field0,field1,field2,field3,field4," +
				"field5,field6,field7,field8,field9",
			allowed: map[string]any{
				"field0": "alias0",
				"field1": "alias1",
				"field2": "alias2",
				"field3": "alias3",
				"field4": "alias4",
				"field5": "alias5",
				"field6": "alias6",
				"field7": "alias7",
				"field8": "alias8",
				"field9": "alias9",
			},
		},
		{
			name: "10-mixed",
			input: "field0,field1,field2,field3,field4," +
				"field5,field6,field7,field8,field9",
			allowed: map[string]any{
				"field0": true,
				"field1": true,
				"field2": true,
				"field3": true,
				"field4": true,
				"field5": "alias5",
				"field6": "alias6",
				"field7": "alias7",
				"field8": "alias8",
				"field9": "alias9",
			},
		},
		{
			name:  "invalid-fields",
			input: "xxx,yyy,zzz",
			allowed: map[string]any{
				"field1": true,
			},
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = queryutil.ResolveAllowedFields(tc.input, tc.allowed)
			}
		})
	}
}

func BenchmarkQueryUtil_ResolveSingleField(b *testing.B) {
	type benchCase struct {
		name         string
		input        string
		allowed      map[string]any
		defaultField string
	}

	cases := []benchCase{
		{
			name:         "ValidBool",
			input:        "username",
			allowed:      map[string]any{"username": true, "email": true},
			defaultField: "default",
		},
		{
			name:         "ValidAlias",
			input:        "email",
			allowed:      map[string]any{"email": "user_email"},
			defaultField: "default",
		},
		{
			name:         "InvalidField",
			input:        "invalid",
			allowed:      map[string]any{"username": true},
			defaultField: "default",
		},
		{
			name:         "EmptyInput",
			input:        "",
			allowed:      map[string]any{"username": true},
			defaultField: "default",
		},
		{
			name:         "FalseBool",
			input:        "disabled",
			allowed:      map[string]any{"disabled": false},
			defaultField: "default",
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = queryutil.ResolveSingleField(tc.input, tc.defaultField, tc.allowed)
			}
		})
	}
}
