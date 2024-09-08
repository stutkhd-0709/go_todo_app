package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := 3333
	t.Setenv("PORT", fmt.Sprint(wantPort))

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	if got.Port != wantPort {
		t.Error("want port", wantPort, "got", got.Port)
	}
	wantEnv := "dev"
	if got.Env != wantEnv {
		t.Error("want env", wantEnv, "got", got.Env)
	}

}
