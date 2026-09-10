package transformations

import (
	"reflect"
	"testing"
)

func TestDeleteTransformation(t *testing.T) {
	tr := &DeleteTransformation{Action: "delete", Rule: []string{"DATABASE_URL", "NOT_THERE"}}

	got, err := tr.Transform(map[string]string{"DATABASE_URL": "postgres://localhost", "KEEP": "me"})
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}

	want := map[string]string{"KEEP": "me"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestRenameTransformation(t *testing.T) {
	tr := &RenameTransformation{Action: "rename", Rule: map[string]string{
		"AWS_BUCKET": "AWS_S3_BUCKET",
		"MISSING":    "NEW_NAME",
	}}

	got, err := tr.Transform(map[string]string{"AWS_BUCKET": "bucket", "OTHER": "value"})
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}

	want := map[string]string{"AWS_S3_BUCKET": "bucket", "OTHER": "value"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTemplateTransformation(t *testing.T) {
	tr := &TemplateTransformation{Action: "template", Rule: map[string]string{
		"SS_DATABASE_SERVER":   "{{ url_host .DATABASE_URL }}",
		"SS_DATABASE_PORT":     "{{ url_port .DATABASE_URL }}",
		"SS_DATABASE_USERNAME": "{{ url_user .DATABASE_URL }}",
		"SS_DATABASE_PASSWORD": "{{ url_password .DATABASE_URL }}",
		"SS_DATABASE_NAME":     `{{ with $x := url_path .DATABASE_URL }}{{ trim_prefix $x "/" }}{{ end }}`,
	}}

	source := map[string]string{"DATABASE_URL": "postgres://user:secret@db.example.com:5432/mydb"}
	got, err := tr.Transform(source)
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}

	want := map[string]string{
		"DATABASE_URL":         "postgres://user:secret@db.example.com:5432/mydb",
		"SS_DATABASE_SERVER":   "db.example.com",
		"SS_DATABASE_PORT":     "5432",
		"SS_DATABASE_USERNAME": "user",
		"SS_DATABASE_PASSWORD": "secret",
		"SS_DATABASE_NAME":     "mydb",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTemplateTransformationBadTemplate(t *testing.T) {
	tr := &TemplateTransformation{Action: "template", Rule: map[string]string{"BROKEN": "{{ url_host }"}}

	if _, err := tr.Transform(map[string]string{}); err == nil {
		t.Error("expected an error for an unparsable template, got nil")
	}
}

func TestTemplateTransformationExecutionError(t *testing.T) {
	tr := &TemplateTransformation{Action: "template", Rule: map[string]string{"BROKEN": "{{ env \"SSM_PARENT_DEFINITELY_UNSET\" }}"}}

	if _, err := tr.Transform(map[string]string{}); err == nil {
		t.Error("expected an error when the template function fails, got nil")
	}
}

func TestTrimTransformation(t *testing.T) {
	tr := &TrimTransformation{Action: "trim_name_prefix", Rule: map[string]string{
		"trim":        "_",
		"starts_with": "_PHP",
	}}

	got, err := tr.Transform(map[string]string{"_PHP_MEMORY_LIMIT": "512M", "_OTHER": "kept"})
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}

	want := map[string]string{"PHP_MEMORY_LIMIT": "512M", "_OTHER": "kept"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTrimTransformationMissingRules(t *testing.T) {
	for name, rule := range map[string]map[string]string{
		"no trim":        {"starts_with": "_PHP"},
		"no starts_with": {"trim": "_"},
		"empty":          {},
	} {
		t.Run(name, func(t *testing.T) {
			tr := &TrimTransformation{Action: "trim_name_prefix", Rule: rule}
			if _, err := tr.Transform(map[string]string{}); err == nil {
				t.Error("expected an error for an incomplete rule, got nil")
			}
		})
	}
}
