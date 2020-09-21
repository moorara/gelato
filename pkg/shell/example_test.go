package shell_test

import (
	"context"
	"fmt"

	"github.com/moorara/gelato/pkg/shell"
)

func ExampleRun() {
	_, out, _ := shell.Run(context.Background(), "echo", "foo", "bar")
	fmt.Println(out)
}

func ExampleRunWith() {
	opts := shell.RunOptions{
		Environment: map[string]string{
			"PLACEHOLDER": "foo bar",
		},
	}

	_, out, _ := shell.RunWith(context.Background(), opts, "printenv", "PLACEHOLDER")
	fmt.Println(out)
}

func ExampleRunner() {
	echo := shell.Runner(context.Background(), "echo", "foo", "bar")
	_, out, _ := echo("baz")
	fmt.Println(out)
}

func ExampleRunnerWith() {
	opts := shell.RunOptions{
		Environment: map[string]string{
			"PLACEHOLDER": "foo bar baz",
		},
	}

	printenv := shell.RunnerWith(context.Background(), opts, "printenv")
	_, out, _ := printenv("PLACEHOLDER")
	fmt.Println(out)
}

func ExampleRunnerFunc_WithArgs() {
	echo := shell.Runner(context.Background(), "echo", "foo")
	echo = echo.WithArgs("bar")
	_, out, _ := echo("baz")
	fmt.Println(out)
}
