package play

import (
	"strings"

	"github.com/flosch/pongo2/v6"
)

func nameToPath(name string) string {
	return strings.Join(strings.Split(name, "-"), ".")
}

func (cmd *Command) envTPL() string {
	base := `
{{ ansible_managed }}
MONITORING_HTTP_ADDR={{ ansible_host }}:{{ ` + nameToPath(cmd.spec.Name) + `.monitoring_port }}
`
	if len(cmd.spec.Resources) != 0 {
		base += "\n## resources\n"
		for _, r := range cmd.spec.Resources {
			if v, ok := cmd.conf.Resources[r]; ok {
				base += v + "\n"
			}
		}
	}

	return base + "\n## from spec file\n" + cmd.spec.Environments
}

func (cmd *Command) fileDiscoveryTPL() string {
	return `
[
{% for server in ` + cmd.spec.Group + `_hostvars %}
  {
    "targets" : ["{{ server.ansible_host }}:{{ ` + nameToPath(cmd.spec.Name) + `.monitoring_port }}"],
    "labels"  : {
      "job"      : "kinescope-app",
      "dc"       : "{{ server.dc }}",
      "rack"     : "{{ server.rack }}",
      "group"    : "` + cmd.spec.Group + `",
      "service"  : "` + cmd.spec.Name + `",
      "hostname" : "{{ server.inventory_hostname }}"
    }
  } {% if not loop.last %},{% endif %}
{% endfor %}
]
`
}

func (cmd *Command) systemdTPL() (string, error) {
	return systemdPongoTPL.Execute(pongo2.Context{
		"name": cmd.spec.Name,
		"exec": map[string]string{
			"stop":   cmd.spec.Exec.Stop,
			"start":  cmd.spec.Exec.Start,
			"reload": cmd.spec.Exec.Reload,
		},
		"limits": map[string]any{
			"mem":     cmd.spec.Limits.Mem,
			"no_file": cmd.spec.Limits.NoFile,
		},
		"security": map[string]string{
			"owner": cmd.spec.Security.Owner,
			"group": cmd.spec.Security.Group,
		},
		"description":  cmd.spec.Description,
		"capabilities": cmd.spec.Capabilities,
	})
}

var systemdPongoTPL = pongo2.Must(pongo2.FromString(`{% verbatim %}## {{ ansible_managed }}{% endverbatim %}
[Unit]
Description={{ description }}.
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User={{ security.owner }}
Group={{ security.group }}
EnvironmentFile=/etc/{{ name }}/environment
ExecStop={{ exec.stop }}
ExecStart={{ exec.start }}
{% if exec.reload != "" %}
ExecReload={{ exec.reload }}
{% endif %}

Restart=always
RestartSec=5

{% if limits.mem != "unlimited" %}
MemoryLimit={{ limits.mem }}
MemoryAccounting=true
{% endif %}

LimitNOFILE={{ limits.no_file }}

RuntimeDirectory={{ name }}
{% for cap in capabilities %}
AmbientCapabilities={{ cap }}
{% endfor %}

[Install]
WantedBy=multi-user.target
`))
