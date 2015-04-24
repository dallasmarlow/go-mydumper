package mydumper

import (
	"path"
	"strconv"
	"strings"
	"time"
)

type MyLoader struct {
	CommonOpts
	ImportPath string
}

func NewLoader(execPath, importPath, username, passwd, db string, threads int) MyLoader {
	return MyLoader{
		CommonOpts: CommonOpts{
			ExecPath:   execPath,
			Username:   username,
			Password:   passwd,
			Database:   db,
			NumThreads: threads,
		},
		ImportPath: importPath,
	}
}

func (l MyLoader) Exec(timeout time.Duration) (string, string, error) {
	stdout, stderr, err := bash.ExecWithTimeout(l.String(), timeout)
	return stdout.String(), stderr.String(), err
}

func (l MyLoader) String() string {
	return strings.Join([]string{
		path.Join(l.ExecPath, `myloader`),
		`--user=` + l.Username,
		`--password=` + l.Password,
		`--database=` + l.Database,
		`--threads=` + strconv.Itoa(l.NumThreads),
		`--directory=` + l.ImportPath}, ` `)
}
