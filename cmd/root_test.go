package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/springload/ssm-parent/ssm/transformations"
)

// loadConfig writes the given yaml to a temp file and runs initSettings against it
func loadConfig(t *testing.T, yaml string) []transformations.Transformation {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0600); err != nil {
		t.Fatalf("can't write the config: %s", err)
	}

	oldConfig, oldList := config, transformationsList
	t.Cleanup(func() {
		config, transformationsList = oldConfig, oldList
		viper.Reset()
	})

	config, transformationsList = path, nil
	initSettings()

	return transformationsList
}

func TestInitSettingsTransformations(t *testing.T) {
	got := loadConfig(t, `
recursive: true
paths: ["/$PROJECT/common/"]

transformations:
    - action: template
      rule:
          SS_DATABASE_SERVER: "{{ url_host .DATABASE_URL }}"
    - action: rename
      rule:
          AWS_BUCKET: AWS_S3_BUCKET
    - action: delete
      rule:
          - DATABASE_URL
    - action: trim_name_prefix
      rule:
          trim: "_"
          starts_with: "_PHP"
`)

	want := []transformations.Transformation{
		&transformations.TemplateTransformation{Action: "template", Rule: map[string]string{"SS_DATABASE_SERVER": "{{ url_host .DATABASE_URL }}"}},
		&transformations.RenameTransformation{Action: "rename", Rule: map[string]string{"AWS_BUCKET": "AWS_S3_BUCKET"}},
		&transformations.DeleteTransformation{Action: "delete", Rule: []string{"DATABASE_URL"}},
		&transformations.TrimTransformation{Action: "trim_name_prefix", Rule: map[string]string{"trim": "_", "starts_with": "_PHP"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}

	if !viper.GetBool("recursive") {
		t.Error("expected 'recursive' to be read from the config file")
	}
	if paths := viper.GetStringSlice("paths"); !reflect.DeepEqual(paths, []string{"/$PROJECT/common/"}) {
		t.Errorf("expected 'paths' to be read from the config file, got %v", paths)
	}
}

func TestInitSettingsUnknownAction(t *testing.T) {
	got := loadConfig(t, `
transformations:
    - action: nonexistent
      rule:
          KEY: value
    - action: delete
      rule:
          - DATABASE_URL
`)

	want := []transformations.Transformation{
		&transformations.DeleteTransformation{Action: "delete", Rule: []string{"DATABASE_URL"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unknown actions should be skipped: got %#v, want %#v", got, want)
	}
}

func TestInitSettingsNoTransformations(t *testing.T) {
	if got := loadConfig(t, "recursive: true\n"); len(got) != 0 {
		t.Errorf("expected no transformations, got %#v", got)
	}
}
