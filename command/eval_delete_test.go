package command

import (
	"errors"
	"testing"

	"github.com/hashicorp/nomad/ci"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/require"
)

func TestEvalDeleteCommand_Run(t *testing.T) {
	ci.Parallel(t)

	testCases := []struct {
		testFn func()
		name   string
	}{
		{
			testFn: func() {

				testServer, client, url := testServer(t, false, nil)
				defer testServer.Shutdown()

				// Create the UI and command.
				ui := cli.NewMockUi()
				cmd := &EvalDeleteCommand{
					Meta: Meta{
						Ui:          ui,
						flagAddress: url,
					},
				}

				// Test basic command input validation.
				require.Equal(t, 1, cmd.Run([]string{"-address=" + url}))
				require.Contains(t, ui.ErrorWriter.String(), "Error validating command args and flags")
				ui.ErrorWriter.Reset()
				ui.OutputWriter.Reset()

				// Try deleting an eval by its ID that doesn't exist.
				require.Equal(t, 1, cmd.Run([]string{"-address=" + url, "fa3a8c37-eac3-00c7-3410-5ba3f7318fd8"}))
				require.Contains(t, ui.ErrorWriter.String(), "404 (eval not found)")
				ui.ErrorWriter.Reset()
				ui.OutputWriter.Reset()

				// Try deleting evals using a filter, which doesn't match any
				// evals found in Nomad.
				require.Equal(t, 1, cmd.Run([]string{"-address=" + url, "-scheduler=service"}))
				require.Contains(t, ui.ErrorWriter.String(), "failed to find any evals that matched filter criteria")
				ui.ErrorWriter.Reset()
				ui.OutputWriter.Reset()

				// Ensure the scheduler config broker is un-paused.
				schedulerConfig, _, err := client.Operator().SchedulerGetConfiguration(nil)
				require.NoError(t, err)
				require.False(t, schedulerConfig.SchedulerConfig.PauseEvalBroker)
			},
			name: "failures",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFn()
		})
	}
}

func TestEvalDeleteCommand_verifyArgsAndFlags(t *testing.T) {
	ci.Parallel(t)

	testCases := []struct {
		inputEvalDeleteCommand *EvalDeleteCommand
		inputArgs              []string
		expectedError          error
		name                   string
	}{
		{
			inputEvalDeleteCommand: &EvalDeleteCommand{
				scheduler:  []string{"batch", "service"},
				status:     []string{"pending", "complete"},
				job:        []string{"example", "another"},
				deployment: []string{"fake-dep-1", "fake-dep-2"},
				node:       []string{"fake-node-1", "fake-node-2"},
			},
			inputArgs:     []string{"fa3a8c37-eac3-00c7-3410-5ba3f7318fd8"},
			expectedError: errors.New("evaluation ID or filter flags required"),
			name:          "arg and flags",
		},
		{
			inputEvalDeleteCommand: &EvalDeleteCommand{
				scheduler:  []string{},
				status:     []string{},
				job:        []string{},
				deployment: []string{},
				node:       []string{},
			},
			inputArgs:     []string{},
			expectedError: errors.New("evaluation ID or filter flags required"),
			name:          "no arg or flags",
		},
		{
			inputEvalDeleteCommand: &EvalDeleteCommand{
				scheduler:  []string{},
				status:     []string{},
				job:        []string{},
				deployment: []string{},
				node:       []string{},
			},
			inputArgs:     []string{"fa3a8c37-eac3-00c7-3410-5ba3f7318fd8", "fa3a8c37-eac3-00c7-3410-5ba3f7318fd9"},
			expectedError: errors.New("expected 1 argument, got 2"),
			name:          "multiple args",
		},
		{
			inputEvalDeleteCommand: &EvalDeleteCommand{
				scheduler:  []string{},
				status:     []string{},
				job:        []string{"example"},
				deployment: []string{},
				node:       []string{},
				filterOp:   "everything",
			},
			inputArgs:     []string{},
			expectedError: errors.New(`got filter-operator "everything", supports "and" "or"`),
			name:          "incorrect filter-operator",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualError := tc.inputEvalDeleteCommand.verifyArgsAndFlags(tc.inputArgs)
			require.Equal(t, tc.expectedError, actualError)
		})
	}
}

func TestEvalDeleteCommand_buildFilter(t *testing.T) {
	ci.Parallel(t)

	testCases := []struct {
		inputEvalDeleteCommand *EvalDeleteCommand
		expectedOutput         string
		name                   string
	}{
		{
			inputEvalDeleteCommand: &EvalDeleteCommand{
				scheduler:  []string{"batch", "service"},
				status:     []string{"pending", "complete"},
				job:        []string{"example", "another"},
				deployment: []string{"fake-dep-1", "fake-dep-2"},
				node:       []string{"fake-node-1", "fake-node-2"},
				filterOp:   " or ",
			},
			expectedOutput: `Type == "batch" or Type == "service" or Status == "pending" or Status == "complete" or JobID == "example" or JobID == "another" or DeploymentID == "fake-dep-1" or DeploymentID == "fake-dep-2" or NodeID == "fake-node-1" or NodeID == "fake-node-2"`,
			name:           "all filter flags with multiple entries",
		},
		{
			inputEvalDeleteCommand: &EvalDeleteCommand{
				scheduler:  []string{"batch", "service"},
				status:     []string{},
				job:        []string{"example", "another"},
				deployment: []string{},
				node:       []string{"fake-node-1", "fake-node-2"},
				filterOp:   " or ",
			},
			expectedOutput: `Type == "batch" or Type == "service" or JobID == "example" or JobID == "another" or NodeID == "fake-node-1" or NodeID == "fake-node-2"`,
			name:           "some filter flags with multiple entries",
		},
		{
			inputEvalDeleteCommand: &EvalDeleteCommand{
				scheduler:  []string{},
				status:     []string{},
				job:        []string{},
				deployment: []string{},
				node:       []string{},
			},
			expectedOutput: "",
			name:           "no filters",
		},
		{
			inputEvalDeleteCommand: &EvalDeleteCommand{
				scheduler:  []string{},
				status:     []string{"pending"},
				job:        []string{"example"},
				deployment: []string{},
				node:       []string{},
				filterOp:   " and ",
			},
			expectedOutput: `Status == "pending" and JobID == "example"`,
			name:           "using and filter operation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := tc.inputEvalDeleteCommand.buildFilter()
			require.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}
