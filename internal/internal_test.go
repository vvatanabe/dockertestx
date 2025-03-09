package internal

import "testing"

func TestGetEnvValue(t *testing.T) {
	// Test case 1: Key exists in the environment slice.
	env := []string{
		"MYSQL_ROOT_PASSWORD=secret",
		"MYSQL_DATABASE=test",
		"FOO=bar",
	}

	got := GetEnvValue(env, "MYSQL_ROOT_PASSWORD")
	want := "secret"
	if got != want {
		t.Errorf("GetEnvValue(env, %q) = %q; want %q", "MYSQL_ROOT_PASSWORD", got, want)
	}

	// Test case 2: Another existing key.
	got = GetEnvValue(env, "MYSQL_DATABASE")
	want = "test"
	if got != want {
		t.Errorf("GetEnvValue(env, %q) = %q; want %q", "MYSQL_DATABASE", got, want)
	}

	// Test case 3: Key that does not exist.
	got = GetEnvValue(env, "NOT_EXISTENT")
	want = ""
	if got != want {
		t.Errorf("GetEnvValue(env, %q) = %q; want empty string", "NOT_EXISTENT", got)
	}

	// Test case 4: Empty environment slice.
	got = GetEnvValue([]string{}, "ANY_VAR")
	if got != "" {
		t.Errorf("GetEnvValue(empty slice, %q) = %q; want empty string", "ANY_VAR", got)
	}

	// Test case 5: Keys with similar prefixes.
	env2 := []string{
		"MY_VAR=123",
		"MY_VARIABLE=456",
	}
	got = GetEnvValue(env2, "MY_VAR")
	want = "123"
	if got != want {
		t.Errorf("GetEnvValue(env2, %q) = %q; want %q", "MY_VAR", got, want)
	}

	got = GetEnvValue(env2, "MY_VARIABLE")
	want = "456"
	if got != want {
		t.Errorf("GetEnvValue(env2, %q) = %q; want %q", "MY_VARIABLE", got, want)
	}
}
