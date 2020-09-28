package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	if os.Getenv("TEST_SUCCESS") == "true" {
		os.Args = []string{"gelato", "-version"}
		main()
	}

	if os.Getenv("TEST_FAIL") == "true" {
		os.Args = []string{"gelato"}
		main()
	}

	name := os.Args[0]
	args := []string{"-test.run=TestMain"}

	t.Run("Success", func(t *testing.T) {
		cmd := exec.Command(name, args...)
		cmd.Env = append(os.Environ(), "TEST_SUCCESS=true")
		err := cmd.Run()
		assert.NoError(t, err)
	})

	t.Run("Fail", func(t *testing.T) {
		cmd := exec.Command(name, args...)
		cmd.Env = append(os.Environ(), "TEST_FAIL=true")
		err := cmd.Run()
		e, ok := err.(*exec.ExitError)
		assert.True(t, ok)
		assert.False(t, e.Success())
	})
}
