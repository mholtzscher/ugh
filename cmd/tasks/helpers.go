package tasks

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"
)

func parseIDs(args []string) ([]int64, error) {
	ids := make([]int64, 0, len(args))
	for _, arg := range args {
		val, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid id %q", arg)
		}
		ids = append(ids, val)
	}
	if len(ids) == 0 {
		return nil, errors.New("at least one id is required")
	}
	return ids, nil
}

func commandArgs(cmd *cli.Command) []string {
	if cmd == nil {
		return nil
	}
	args := make([]string, cmd.Args().Len())
	for i := 0; i < cmd.Args().Len(); i++ {
		args[i] = cmd.Args().Get(i)
	}
	return args
}
