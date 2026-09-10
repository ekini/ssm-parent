package ssm

import (
	"reflect"
	"testing"
)

func TestExpandArgs(t *testing.T) {
	t.Setenv("PROJECT", "myproject")
	t.Setenv("ENVIRONMENT", "production")

	args := []string{"/$PROJECT/$ENVIRONMENT/backend/", "/${PROJECT}/common/", "/literal/path/"}
	want := []string{"/myproject/production/backend/", "/myproject/common/", "/literal/path/"}

	if got := ExpandArgs(args); !reflect.DeepEqual(got, want) {
		t.Errorf("ExpandArgs(%v) = %v, want %v", args, got, want)
	}
}

func TestExpandArgsUnsetVar(t *testing.T) {
	args := []string{"/$SSM_PARENT_DEFINITELY_UNSET/path/"}
	want := []string{"//path/"}

	if got := ExpandArgs(args); !reflect.DeepEqual(got, want) {
		t.Errorf("ExpandArgs(%v) = %v, want %v", args, got, want)
	}
}

func TestExpandValueTrimsSpace(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")

	if got := expandValue("  $ENVIRONMENT  "); got != "production" {
		t.Errorf("expected 'production', got '%s'", got)
	}
}

func TestExpandParametersNoop(t *testing.T) {
	parameters := map[string]string{"SOME_SECRET": "abc$abc"}

	if err := expandParameters(parameters, false, nil); err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}

	if parameters["SOME_SECRET"] != "abc$abc" {
		t.Errorf("expected value to be untouched, got '%s'", parameters["SOME_SECRET"])
	}
}
