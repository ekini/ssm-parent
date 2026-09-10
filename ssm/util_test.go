package ssm

import (
	"reflect"
	"testing"
)

func TestStringSliceDifference(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{"all missing", []string{"a", "b"}, []string{}, []string{"a", "b"}},
		{"none missing", []string{"a", "b"}, []string{"a", "b"}, []string{}},
		{"some missing", []string{"a", "b", "c"}, []string{"b"}, []string{"a", "c"}},
		{"extra in b", []string{"a"}, []string{"a", "z"}, []string{}},
		{"empty a", nil, []string{"a"}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringSliceDifference(tt.a, tt.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringSliceDifference(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestEnvGet(t *testing.T) {
	t.Setenv("SSM_PARENT_TEST_VAR", "value")

	var e Env

	if val, ok := e.Get("SSM_PARENT_TEST_VAR"); !ok || val != "value" {
		t.Errorf("expected ('value', true), got ('%s', %t)", val, ok)
	}

	if val, ok := e.Get("SSM_PARENT_TEST_VAR_UNSET"); ok || val != "" {
		t.Errorf("expected ('', false), got ('%s', %t)", val, ok)
	}
}
