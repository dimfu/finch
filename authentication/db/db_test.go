package db

import "testing"

func TestDBNoEnvironments(t *testing.T) {
	t.Setenv("POSTGRES_USER", "")
	t.Setenv("POSTGRES_PASSWORD", "")
	t.Setenv("POSTGRES_DBNAME", "")
	_, err := New()
	if err == nil {
		t.Fatal(err)
	}
}
