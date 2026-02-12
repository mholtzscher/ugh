package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestScripts(t *testing.T) {
	binDir := t.TempDir()
	binPath := filepath.Join(binDir, "ugh")

	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %v\n%s", err, output)

	params := testscript.Params{
		Dir: filepath.Join("testdata", "script"),
		Setup: func(env *testscript.Env) error {
			path := binDir + string(os.PathListSeparator) + os.Getenv("PATH")
			env.Setenv("PATH", path)
			env.Setenv("HOME", filepath.Join(env.WorkDir, "home"))
			env.Setenv("XDG_CONFIG_HOME", filepath.Join(env.WorkDir, "home", ".config"))
			env.Setenv("XDG_DATA_HOME", filepath.Join(env.WorkDir, "home", ".local", "share"))
			env.Setenv("TZ", "UTC")
			return nil
		},
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"dbquery": cmdDBQuery,
			"dbexec":  cmdDBExec,
		},
	}

	testscript.Run(t, params)
}

func cmdDBQuery(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("dbquery does not support !")
	}
	if len(args) < 2 {
		ts.Fatalf("usage: dbquery <db> <sql...>")
	}

	db := openScriptDB(ts, args[0])
	defer func() { _ = db.Close() }()
	query := strings.Join(args[1:], " ")

	rows, err := db.QueryContext(context.Background(), query)
	ts.Check(err)
	defer rows.Close()

	cols, err := rows.Columns()
	ts.Check(err)

	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	for rows.Next() {
		err = rows.Scan(ptrs...)
		ts.Check(err)

		fields := make([]string, len(vals))
		for i, v := range vals {
			fields[i] = sqlValueString(v)
		}
		_, err = fmt.Fprintln(ts.Stdout(), strings.Join(fields, "\t"))
		ts.Check(err)
	}

	ts.Check(rows.Err())
}

func cmdDBExec(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("dbexec does not support !")
	}
	if len(args) < 2 {
		ts.Fatalf("usage: dbexec <db> <sql...>")
	}

	db := openScriptDB(ts, args[0])
	defer func() { _ = db.Close() }()
	statement := strings.Join(args[1:], " ")

	result, err := db.ExecContext(context.Background(), statement)
	ts.Check(err)

	affected, err := result.RowsAffected()
	ts.Check(err)

	_, err = fmt.Fprintln(ts.Stdout(), affected)
	ts.Check(err)
}

func openScriptDB(ts *testscript.TestScript, path string) *sql.DB {
	dbPath := ts.MkAbs(path)

	db, err := sql.Open("turso", dbPath)
	ts.Check(err)

	_, err = db.ExecContext(context.Background(), "PRAGMA foreign_keys=ON;")
	ts.Check(err)

	return db
}

func sqlValueString(v any) string {
	switch value := v.(type) {
	case nil:
		return ""
	case []byte:
		return string(value)
	case string:
		return value
	default:
		return fmt.Sprint(value)
	}
}
