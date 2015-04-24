package mydumper

import (
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type CommonOpts struct {
	ExecPath, ExportRoot, Username, Password, Database string
	NumThreads                                         int
}

type MyDumper struct {
	CommonOpts
	Table, OutputDir       string
	Compress, NoSchemas    bool
	RowRangeSize, StmtSize int
}

func New(execPath, expRoot, username, passwd, db, table string, threads int) MyDumper {
	return MyDumper{
		CommonOpts: CommonOpts{
			ExecPath:   execPath,
			ExportRoot: expRoot,
			Username:   username,
			Password:   passwd,
			Database:   db,
			NumThreads: threads,
		},
		Table: table,
	}
}

func (d MyDumper) ExportPath() string {
	return path.Join(d.ExportRoot, d.Database, d.Table)
}

func (d MyDumper) OutputDirOrExportPath() string {
	if d.OutputDir == "" {
		return d.ExportPath()
	}

	return d.OutputDir
}

func (d MyDumper) CreateExportDir() error {
	return os.MkdirAll(d.OutputDirOrExportPath(), 0755)
}

func (d MyDumper) Exec(timeout time.Duration) (string, string, error) {
	if err := d.CreateExportDir(); err != nil {
		return ``, ``, err
	}

	stdout, stderr, err := ExecWithTimeout(d.String(), timeout)
	return stdout.String(), stderr.String(), err
}

func (d MyDumper) String() string {
	args := []string{
		`--user=` + d.Username,
		`--password=` + d.Password,
		`--database=` + d.Database,
		`--tables-list=` + d.Table,
		`--threads=` + strconv.Itoa(d.NumThreads),
		`--outputdir=` + d.OutputDirOrExportPath(),
	}

	if d.Compress {
		args = append(args, `--compress`)
	}
	if d.NoSchemas {
		args = append(args, `--no-schemas`)
	}
	if d.RowRangeSize != 0 {
		args = append(args, `--rows=`+strconv.Itoa(d.RowRangeSize))
	}
	if d.StmtSize != 0 {
		args = append(args, `--statement-size=`+strconv.Itoa(d.StmtSize))
	}

	return strings.Join(append([]string{path.Join(d.ExecPath, `mydumper`)}, args...), ` `)
}
