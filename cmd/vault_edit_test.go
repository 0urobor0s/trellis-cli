package cmd

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"trellis-cli/trellis"
)

func TestVaultEditRunValidations(t *testing.T) {
	ui := cli.NewMockUi()

	cases := []struct {
		name            string
		projectDetected bool
		args            []string
		out             string
		code            int
	}{
		{
			"no_project",
			false,
			nil,
			"No Trellis project detected",
			1,
		},
		{
			"too_many_args",
			true,
			[]string{"group_vars/all/vault.yml", "group_vars/production/vault.yml"},
			"Error: too many arguments",
			1,
		},
	}

	for _, tc := range cases {
		mockProject := &MockProject{tc.projectDetected}
		trellis := trellis.NewTrellis(mockProject)
		vaultEditCommand := VaultEditCommand{ui, trellis}

		code := vaultEditCommand.Run(tc.args)

		if code != tc.code {
			t.Errorf("expected code %d to be %d", code, tc.code)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

		if !strings.Contains(combined, tc.out) {
			t.Errorf("expected output %q to contain %q", combined, tc.out)
		}
	}
}

func TestVaultEditRun(t *testing.T) {
	defer trellis.LoadFixtureProject(t)()
	defer MockExec(t)()

	ui := cli.NewMockUi()
	project := &trellis.Project{}
	trellis := trellis.NewTrellis(project)
	vaultEditCommand := VaultEditCommand{ui, trellis}

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"default",
			[]string{"group_vars/development/vault.yml"},
			"ansible-vault edit group_vars/development/vault.yml",
			0,
		},
	}

	for _, tc := range cases {
		code := vaultEditCommand.Run(tc.args)

		if code != tc.code {
			t.Errorf("expected code %d to be %d", code, tc.code)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

		if !strings.Contains(combined, tc.out) {
			t.Errorf("expected output %q to contain %q", combined, tc.out)
		}
	}
}
