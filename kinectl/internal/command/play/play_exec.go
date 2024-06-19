package play

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type (
	PlayBook []Target
	Target   struct {
		Hosts        string            `yaml:"hosts"`
		Name         string            `yaml:"name,omitempty"`
		User         string            `yaml:"user,omitempty"`
		Serial       int               `yaml:"serial,omitempty"`
		Become       bool              `yaml:"become,omitempty"`
		BecomeMethod string            `yaml:"become_method,omitempty"`
		GatherFacts  bool              `yaml:"gather_facts"`
		Vars         map[string]string `yaml:"vars,omitempty"`
		Tasks        []map[string]any  `yaml:"tasks,omitempty"`
		Handlers     []map[string]any  `yaml:"handlers,omitempty"`
	}
)

func (cmd *Command) play() error {
	var (
		tasks    []map[string]any
		vPath    = nameToPath(cmd.spec.Name) + ".version"
		pbPath   = filepath.Join(os.TempDir(), "kinectl", "plb_"+uuid.New().String())
		envPath  = filepath.Join(os.TempDir(), "kinectl", "env_"+uuid.New().String())
		srvPath  = filepath.Join(os.TempDir(), "kinectl", "srt_"+uuid.New().String())
		fDSPath  = filepath.Join(os.TempDir(), "kinectl", "fds_"+uuid.New().String())
		varsPath = filepath.Join(os.TempDir(), "kinectl", "var_"+uuid.New().String()+".yml")
	)
	if err := os.MkdirAll(filepath.Join(os.TempDir(), "kinectl"), os.ModePerm); err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(filepath.Join(os.TempDir(), "kinectl"))
	}()
	unit, err := cmd.systemdTPL()
	if err != nil {
		return err
	}

	if err := os.WriteFile(srvPath, []byte(unit), 0644); err != nil {
		return err
	}

	cmd.debug("unit file \n %s", unit)

	if err := os.WriteFile(envPath, []byte(cmd.envTPL()), 0644); err != nil {
		return err
	}

	cmd.debug("environment \n %s", cmd.envTPL())

	if err := os.WriteFile(fDSPath, []byte(cmd.fileDiscoveryTPL()), 0644); err != nil {
		return err
	}

	cmd.debug("monitoring FD \n %s", cmd.fileDiscoveryTPL())

	var (
		sources  []string
		playbook PlayBook
	)

	list, err := os.ReadDir(filepath.Join(cmd.opt.OpsDir, "inventories", "vars"))
	if err != nil {
		return err
	}

	for _, e := range list {
		sources = append(sources, filepath.Join("inventories", "vars", e.Name()))
		sources = append(sources, filepath.Join("inventories", cmd.opt.Env, "override_vars", e.Name()))
	}

	playbook = append(playbook, Target{
		Name:  "include vars",
		Hosts: "localhost",
		Tasks: []map[string]any{
			{
				"merge_yaml": map[string]any{
					"sources": sources,
					"dest":    varsPath,
				},
				"check_mode": false,
			},
		},
	})

	tasks = append(tasks, map[string]any{
		"ansible.builtin.include_vars": map[string]any{
			"file": varsPath,
		},
	})

	tasks = append(tasks, map[string]any{
		"name": "set version",
		"set_fact": map[string]string{
			"app_version": "{% if " + vPath + " is defined and " + vPath + "|length %}{{ " + vPath + " }}{% else %}{% endif %}",
		},
	})

	tasks = append(tasks, map[string]any{
		"name": "install " + cmd.spec.Name,
		"ansible.builtin.apt": map[string]any{
			"name":         cmd.spec.Name + "{% if app_version|length %}={{ app_version }}{% endif %}",
			"state":        "{% if app_version|length %}present{% else %}latest{% endif %}",
			"force":        true,
			"update_cache": true,
			"dpkg_options": "force-downgrade",
		}},
	)

	tasks = append(tasks, map[string]any{
		"ansible.builtin.file": map[string]string{
			"path":  filepath.Join("/etc", cmd.spec.Name),
			"mode":  "0755",
			"state": "directory",
			"owner": cmd.spec.Security.Owner,
			"group": cmd.spec.Security.Group,
		},
	})

	tasks = append(tasks, map[string]any{
		"ansible.builtin.template": map[string]string{
			"src":  envPath,
			"dest": filepath.Join("/etc", cmd.spec.Name, "environment"),
			"mode": "0644",
		}},
	)

	if !cmd.spec.Cli {
		tasks = append(tasks, map[string]any{
			"ansible.builtin.template": map[string]string{
				"src":  srvPath,
				"dest": filepath.Join("/etc/systemd/system/", cmd.spec.Name+".service"),
				"mode": "0644",
			},
		})
		tasks = append(tasks, map[string]any{
			"ansible.builtin.service": map[string]any{
				"name":          cmd.spec.Name,
				"state":         "started",
				"enabled":       true,
				"daemon_reload": true,
			},
		})
	}

	for _, role := range cmd.spec.Deploy.Roles {
		tasks = append(tasks, map[string]any{
			"ansible.builtin.include_role": map[string]string{
				"name": role,
			},
		})
	}

	if !cmd.spec.Cli {
		tasks = append(tasks, map[string]any{
			"ansible.builtin.service": map[string]any{
				"name":          cmd.spec.Name,
				"state":         "restarted",
				"enabled":       true,
				"daemon_reload": true,
			},
		})
	}

	if probe := cmd.spec.Deploy.Probe; !cmd.spec.Cli && probe != nil {
		switch probe.Type {
		case "http", "https":
			url := probe.Type + "://" + probe.Host + ":" + probe.Port
			if len(probe.Path) != 0 {
				url += probe.Path
			}

			tasks = append(tasks, map[string]any{
				"name": "http probe",
				"ansible.builtin.uri": map[string]any{
					"url":              url,
					"follow_redirects": true,
					"method":           "GET",
				},
				//"delegate_to": "127.0.0.1",
				"register": "_result",
				"until":    "_result.status == 200",
				"retries":  12,
				"delay":    5,
			})
		case "tcp":
			tasks = append(tasks, map[string]any{
				"name": "tcp probe",
				"ansible.builtin.wait_for": map[string]any{
					"host":    probe.Host,
					"port":    probe.Port,
					"state":   "started",
					"delay":   0,
					"timeout": 5,
				},
				//"delegate_to": "127.0.0.1",
			})
		default:
			return fmt.Errorf("invalid probe type: %q", probe.Type)
		}
	}

	playbook = append(playbook, Target{
		Name:         "deploy " + cmd.spec.Name,
		User:         cmd.conf.User,
		Vars:         cmd.conf.Vars,
		Hosts:        cmd.spec.Group,
		Serial:       cmd.spec.Deploy.Serial,
		Become:       cmd.conf.Become,
		BecomeMethod: cmd.conf.BecomeMethod,
		GatherFacts:  cmd.conf.GatherFacts,
		Tasks:        tasks,
	})

	if !cmd.spec.Cli {
		playbook = append(playbook, Target{
			Name:         "monitoring SD file",
			User:         cmd.conf.User,
			Vars:         cmd.conf.Vars,
			Hosts:        cmd.conf.Monitoring.Hosts,
			Become:       cmd.conf.Become,
			BecomeMethod: cmd.conf.BecomeMethod,
			GatherFacts:  cmd.conf.GatherFacts,
			Tasks: []map[string]any{
				{
					"ansible.builtin.include_vars": map[string]any{
						"file": varsPath,
					},
				},
				{
					"ansible.builtin.template": map[string]string{
						"src":   fDSPath,
						"dest":  filepath.Join(cmd.conf.Monitoring.FileSDDir, cmd.spec.Name+".json"),
						"mode":  "0644",
						"owner": "prometheus",
						"group": "prometheus",
					},
				},
			},
		})
	}

	pb, err := yaml.Marshal(playbook)
	if err != nil {
		return err
	}

	cmd.debug("playbook \n %s", pb)

	if err := os.WriteFile(pbPath, pb, 0644); err != nil {
		return err
	}
	args := []string{
		"-i", "inventories/" + cmd.opt.Env + "/hosts.yml",
		pbPath,
		"--extra-vars", "service_name=" + cmd.spec.Name,
	}
	if cmd.opt.DryRun {
		args = append(args, "--check")
	}
	if cmd.opt.Diff {
		args = append(args, "--diff")
	}
	if len(cmd.opt.Limit) != 0 {
		args = append(args, "--limit", cmd.opt.Limit)
	}
	if cmd.opt.Debug {
		args = append(args, "--syntax-check")
	}
	c := exec.Command("ansible-playbook", args...)
	{
		c.Dir = cmd.opt.OpsDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
	}
	fmt.Println(c.String())
	return c.Run()
}
