package app

import (
	stderrors "errors"
	"fmt"
	"os"
	"strings"

	"unifai-cli/internal/command"
	clierrors "unifai-cli/internal/errors"
)

func Run(args []string) int {
	root := command.NewRootCommand()
	root.SetArgs(args)

	if err := root.Execute(); err != nil {
		var usageErr *clierrors.UsageError
		if stderrors.As(err, &usageErr) || isUsageError(err) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return clierrors.ExitUsage
		}

		fmt.Fprintln(os.Stderr, "Error:", err)
		return clierrors.ExitError
	}

	return clierrors.ExitOK
}

func isUsageError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	usagePrefixes := []string{
		"unknown command",
		"unknown flag",
		"unknown shorthand flag",
		"required flag",
		"invalid argument",
		"flag needs an argument",
	}
	for _, p := range usagePrefixes {
		if strings.HasPrefix(msg, p) {
			return true
		}
	}
	return false
}
