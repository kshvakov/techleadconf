package play

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kshvakov/techleadconf/kinectl/config"
)

type Options struct {
	Env      string
	Diff     bool
	Debug    bool
	Limit    string
	OpsDir   string
	DryRun   bool
	SpecFile string
}

func New(o *Options) (*Command, error) {
	conf, err := config.ParseConf(filepath.Join(o.OpsDir, ".kinectl.yml"))
	if err != nil {
		return nil, err
	}
	spec, err := config.ParseSpec(o.SpecFile)
	if err != nil {
		return nil, err
	}
	cmd := Command{
		opt:   o,
		spec:  spec,
		conf:  conf,
		debug: func(string, ...any) {},
	}
	if o.Debug {
		cmd.debug = log.New(os.Stdout, "[kinectl debug] ", log.Lshortfile).Printf
	}
	return &cmd, nil
}

type Command struct {
	opt   *Options
	spec  *config.Spec
	conf  *config.Config
	debug func(string, ...any)
}

func (cmd *Command) Run() error {
	return cmd.play()
}
