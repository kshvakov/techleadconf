package deb

import (
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/goreleaser/nfpm/v2"
	"github.com/goreleaser/nfpm/v2/files"
	"github.com/kshvakov/techleadconf/kinectl/config"
)

type Options struct {
	Target   string
	Version  string
	SpecFile string
}

func New(o *Options) (*Command, error) {
	if _, err := semver.NewVersion(o.Version); err != nil {
		return nil, err
	}

	spec, err := config.ParseSpec(o.SpecFile)
	if err != nil {
		return nil, err
	}
	return &Command{
		spec:    spec,
		target:  o.Target,
		version: o.Version,
	}, nil
}

type Command struct {
	spec    *config.Spec
	target  string
	version string
}

func (cmd *Command) Run() error {
	pkg, err := nfpm.Get("deb")
	if err != nil {
		return err
	}
	var fileInfo *files.ContentFileInfo

	fileInfo = &files.ContentFileInfo{
		Owner: cmd.spec.Security.Owner,
		Group: cmd.spec.Security.Group,
	}

	info := nfpm.WithDefaults(&nfpm.Info{
		Name:        cmd.spec.Name,
		Vendor:      "kinescope",
		Version:     cmd.version,
		Section:     "default",
		Platform:    "linux",
		Priority:    "extra",
		Homepage:    "https://kinescope.io",
		Maintainer:  "corp@kinescope.io",
		Description: cmd.spec.Description,
		Overridables: nfpm.Overridables{
			Contents: files.Contents{
				{
					Source:      filepath.Join(cmd.target, cmd.spec.Name),
					Destination: filepath.Join("/usr/bin", cmd.spec.Name),
					FileInfo:    fileInfo,
				},
			},
		},
	})

	info.Target = filepath.Join(cmd.target, pkg.ConventionalFileName(info))

	f, err := os.Create(info.Target)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := pkg.Package(info, f); err != nil {
		os.Remove(info.Target)
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return f.Close()
}
