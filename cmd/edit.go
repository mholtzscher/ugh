package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mholtzscher/ugh/internal/editor"
	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/urfave/cli/v3"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var editCmd = &cli.Command{
	Name:     "edit",
	Aliases:  []string{"e"},
	Usage:    "Edit a task",
	Category: "Tasks",
	Description: `Edit a task by ID.

Opens the task in your editor ($VISUAL or $EDITOR) by default.
Use flags for quick single-field changes without opening an editor.

		Examples:
		  ugh edit 1                          # Open in editor (default)
		  ugh edit 1 --state now              # Move to Now
		  ugh edit 1 --due 2026-02-10         # Set due date
		  ugh edit 1 --title "New title"      # Change title
		  ugh edit 1 -p urgent                # Add project 'urgent'
		  ugh edit 1 --remove-project old     # Remove project 'old'
		  ugh edit 1 -c work -m key:val       # Add context and metadata`,
	ArgsUsage: "<id>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    flags.FlagTitle,
			Aliases: []string{flags.FlagDescription},
			Usage:   "update title",
		},
		&cli.StringFlag{
			Name:  flags.FlagNotes,
			Usage: "update notes",
		},
		&cli.StringFlag{
			Name:  flags.FlagState,
			Usage: "set state (" + flags.TaskStatesUsage + ")",
			Action: flags.StringAction(
				flags.OneOfCaseInsensitiveRule(flags.FieldState, flags.TaskStates()...),
			),
		},
		&cli.StringFlag{
			Name:  flags.FlagDueOn,
			Usage: "set due date (" + flags.DateTextYYYYMMDD + ")",
			Action: flags.StringAction(
				flags.DateLayoutRule(flags.FieldDate, flags.DateLayoutYYYYMMDD, flags.DateTextYYYYMMDD),
			),
		},
		&cli.BoolFlag{
			Name:  flags.FlagNoDue,
			Usage: "clear due date",
		},
		&cli.StringFlag{
			Name:  flags.FlagWaitingFor,
			Usage: "set waiting-for value",
		},
		&cli.BoolFlag{
			Name:  flags.FlagNoWaitingFor,
			Usage: "clear waiting-for value",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagProject,
			Aliases: []string{"p"},
			Usage:   "add project (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagContext,
			Aliases: []string{"c"},
			Usage:   "add context (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagMeta,
			Aliases: []string{"m"},
			Usage:   "set metadata " + flags.MetaTextKeyValue + " (repeatable)",
			Action: flags.StringSliceAction(
				flags.EachContainsSeparatorRule(flags.FieldMeta, flags.MetaSeparatorColon, flags.MetaTextKeyValue),
			),
		},
		&cli.BoolFlag{
			Name:    flags.FlagDone,
			Aliases: []string{"x"},
			Usage:   "mark as done (state=done)",
			Action: flags.BoolAction(
				flags.MutuallyExclusiveBoolFlagsRule(flags.FlagDone, flags.FlagUndone),
				flags.BoolRequiresStringOneOfCaseInsensitiveRule(flags.FlagDone, flags.FlagState, flags.TaskStateDone),
			),
		},
		&cli.BoolFlag{
			Name:  flags.FlagUndone,
			Usage: "reopen (undo done)",
			Action: flags.BoolAction(
				flags.MutuallyExclusiveBoolFlagsRule(flags.FlagDone, flags.FlagUndone),
			),
		},
		&cli.StringSliceFlag{
			Name:  flags.FlagRemoveProject,
			Usage: "remove project (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:  flags.FlagRemoveContext,
			Usage: "remove context (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:  flags.FlagRemoveMeta,
			Usage: "remove metadata key (repeatable)",
		},
		&cli.BoolFlag{
			Name:    flags.FlagEditor,
			Aliases: []string{"e"},
			Usage:   "open in $VISUAL/$EDITOR",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != 1 {
			return errors.New("edit requires a task id")
		}
		ids, err := parseIDs(commandArgs(cmd))
		if err != nil {
			return err
		}
		id := ids[0]

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		err = maybeSyncBeforeWrite(ctx, svc)
		if err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		var saved *store.Task
		changed := false
		hasFields := hasFieldFlags(cmd)
		if cmd.Bool(flags.FlagEditor) && hasFields {
			return errors.New("cannot combine field flags with --editor")
		}

		if cmd.Bool(flags.FlagEditor) || !hasFields {
			saved, changed, err = runEditorMode(ctx, svc, id)
			if err != nil {
				return err
			}
		} else {
			saved, err = runFlagsMode(ctx, cmd, svc, id)
			if err != nil {
				return err
			}
			changed = true
		}

		if !changed {
			writer := outputWriter()
			_, err = fmt.Fprintln(writer.Out, "No changes made")
			return err
		}

		err = maybeSyncAfterWrite(ctx, svc)
		if err != nil {
			return fmt.Errorf("sync push: %w", err)
		}

		writer := outputWriter()
		return writer.WriteTask(saved)
	},
}

func hasFieldFlags(cmd *cli.Command) bool {
	return cmd.String(flags.FlagTitle) != "" ||
		cmd.String(flags.FlagNotes) != "" ||
		cmd.String(flags.FlagState) != "" ||
		cmd.String(flags.FlagDueOn) != "" ||
		cmd.Bool(flags.FlagNoDue) ||
		cmd.String(flags.FlagWaitingFor) != "" ||
		cmd.Bool(flags.FlagNoWaitingFor) ||
		len(cmd.StringSlice(flags.FlagProject)) > 0 ||
		len(cmd.StringSlice(flags.FlagContext)) > 0 ||
		len(cmd.StringSlice(flags.FlagMeta)) > 0 ||
		cmd.Bool(flags.FlagDone) ||
		cmd.Bool(flags.FlagUndone) ||
		len(cmd.StringSlice(flags.FlagRemoveProject)) > 0 ||
		len(cmd.StringSlice(flags.FlagRemoveContext)) > 0 ||
		len(cmd.StringSlice(flags.FlagRemoveMeta)) > 0
}

func runEditorMode(ctx context.Context, svc service.Service, id int64) (*store.Task, bool, error) {
	task, err := svc.GetTask(ctx, id)
	if err != nil {
		return nil, false, err
	}

	edited, changed, err := editor.Edit(task)
	if err != nil {
		return nil, false, fmt.Errorf("editor: %w", err)
	}

	if !changed || edited == nil {
		return nil, false, nil
	}

	updatedTask, updateErr := svc.FullUpdateTask(ctx, service.FullUpdateTaskRequest{
		ID:         id,
		Title:      edited.Title,
		Notes:      edited.Notes,
		State:      edited.State,
		DueOn:      edited.DueOn,
		WaitingFor: edited.WaitingFor,
		Projects:   edited.Projects,
		Contexts:   edited.Contexts,
		Meta:       edited.Meta,
	})
	if updateErr != nil {
		return nil, false, updateErr
	}

	return updatedTask, true, nil
}

func runFlagsMode(ctx context.Context, cmd *cli.Command, svc service.Service, id int64) (*store.Task, error) {
	meta, err := parseMetaFlags(cmd.StringSlice(flags.FlagMeta))
	if err != nil {
		return nil, fmt.Errorf("parse meta: %w", err)
	}

	req := service.UpdateTaskRequest{
		ID:              id,
		AddProjects:     cmd.StringSlice(flags.FlagProject),
		AddContexts:     cmd.StringSlice(flags.FlagContext),
		SetMeta:         meta,
		RemoveProjects:  cmd.StringSlice(flags.FlagRemoveProject),
		RemoveContexts:  cmd.StringSlice(flags.FlagRemoveContext),
		RemoveMetaKeys:  cmd.StringSlice(flags.FlagRemoveMeta),
		ClearDueOn:      cmd.Bool(flags.FlagNoDue),
		ClearWaitingFor: cmd.Bool(flags.FlagNoWaitingFor),
	}

	if title := cmd.String(flags.FlagTitle); title != "" {
		req.Title = &title
	}
	if notes := cmd.String(flags.FlagNotes); notes != "" {
		req.Notes = &notes
	}
	if state := cmd.String(flags.FlagState); state != "" {
		req.State = &state
	}

	if due := cmd.String(flags.FlagDueOn); due != "" {
		req.DueOn = &due
	}
	if waitingFor := cmd.String(flags.FlagWaitingFor); waitingFor != "" {
		req.WaitingFor = &waitingFor
	}

	// Apply field updates first.
	updated, err := svc.UpdateTask(ctx, req)
	if err != nil {
		return nil, err
	}
	if cmd.Bool(flags.FlagDone) {
		_, err = svc.SetDone(ctx, []int64{id}, true)
		if err != nil {
			return nil, err
		}
		return svc.GetTask(ctx, id)
	}
	if cmd.Bool(flags.FlagUndone) {
		_, err = svc.SetDone(ctx, []int64{id}, false)
		if err != nil {
			return nil, err
		}
		return svc.GetTask(ctx, id)
	}
	return updated, nil
}

func parseMetaFlags(meta []string) (map[string]string, error) {
	if len(meta) == 0 {
		return map[string]string{}, nil
	}
	result := make(map[string]string, len(meta))
	for _, m := range meta {
		k, v, ok := strings.Cut(m, flags.MetaSeparatorColon)
		if !ok {
			return nil, fmt.Errorf("invalid meta format: %s (expected %s)", m, flags.MetaTextKeyValue)
		}
		result[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return result, nil
}
