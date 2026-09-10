package ssm

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func TestCollectJsonParameters(t *testing.T) {
	input := []*ssm.Parameter{
		{Name: aws.String("/project/env/one"), Value: aws.String(`{"USERNAME": "myuser", "DATABASE": "production"}`)},
		{Name: aws.String("/project/env/two"), Value: aws.String(`{"DATABASE": "test"}`)},
	}

	parameters, errs := collectJsonParameters(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	want := []map[string]string{
		{"USERNAME": "myuser", "DATABASE": "production"},
		{"DATABASE": "test"},
	}
	if !reflect.DeepEqual(parameters, want) {
		t.Errorf("got %v, want %v", parameters, want)
	}
}

func TestCollectJsonParametersInvalidJson(t *testing.T) {
	input := []*ssm.Parameter{
		{Name: aws.String("/project/env/broken"), Value: aws.String("not json")},
		{Name: aws.String("/project/env/ok"), Value: aws.String(`{"KEY": "value"}`)},
	}

	parameters, errs := collectJsonParameters(input)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}

	want := []map[string]string{{"KEY": "value"}}
	if !reflect.DeepEqual(parameters, want) {
		t.Errorf("got %v, want %v", parameters, want)
	}
}

func TestCollectPlainParameters(t *testing.T) {
	input := []*ssm.Parameter{
		{Name: aws.String("/project/environment/myParameter"), Value: aws.String("supervalue")},
		{Name: aws.String("bare"), Value: aws.String("barevalue")},
	}

	parameters, errs := collectPlainParameters(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}

	want := []map[string]string{
		{"myParameter": "supervalue"},
		{"bare": "barevalue"},
	}
	if !reflect.DeepEqual(parameters, want) {
		t.Errorf("got %v, want %v", parameters, want)
	}
}
