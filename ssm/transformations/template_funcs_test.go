package transformations

import "testing"

func TestGetEnv(t *testing.T) {
	t.Setenv("SSM_PARENT_TEST_VAR", "value")

	got, err := GetEnv("SSM_PARENT_TEST_VAR")
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}
	if got != "value" {
		t.Errorf("expected 'value', got '%s'", got)
	}

	if _, err := GetEnv("SSM_PARENT_DEFINITELY_UNSET"); err == nil {
		t.Error("expected an error for an unset variable, got nil")
	}
}

func TestURLFuncs(t *testing.T) {
	tests := []struct {
		name  string
		fn    func(string) (string, error)
		input string
		want  string
	}{
		{"user", URLUser, "postgres://user:secret@db.example.com:5432/mydb", "user"},
		{"user unset", URLUser, "postgres://db.example.com/mydb", ""},
		{"password", URLPassword, "postgres://user:secret@db.example.com:5432/mydb", "secret"},
		{"password unset", URLPassword, "postgres://user@db.example.com/mydb", ""},
		{"scheme", URLScheme, "postgres://user:secret@db.example.com:5432/mydb", "postgres"},
		{"host", URLHost, "postgres://user:secret@db.example.com:5432/mydb", "db.example.com"},
		{"host without port", URLHost, "https://example.com/path", "example.com"},
		{"port", URLPort, "postgres://user:secret@db.example.com:5432/mydb", "5432"},
		{"port unset", URLPort, "https://example.com/path", ""},
		{"path", URLPath, "postgres://user:secret@db.example.com:5432/mydb", "/mydb"},
		{"path unset", URLPath, "https://example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fn(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got: %s", err)
			}
			if got != tt.want {
				t.Errorf("got '%s', want '%s'", got, tt.want)
			}
		})
	}
}

func TestURLFuncsInvalidURL(t *testing.T) {
	const invalid = "://not a url"

	for name, fn := range map[string]func(string) (string, error){
		"user":     URLUser,
		"password": URLPassword,
		"scheme":   URLScheme,
		"host":     URLHost,
		"port":     URLPort,
		"path":     URLPath,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := fn(invalid); err == nil {
				t.Error("expected an error for an unparsable URL, got nil")
			}
		})
	}
}
