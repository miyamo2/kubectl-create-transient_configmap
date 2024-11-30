/*
Copyright (c) 2024 miyamo2

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"os/exec"
	"strings"
)

const version = "0.1.0"

var (
	// transientConfigMapUse is the one-line usage message for the transient_configmap command
	transientConfigMapUse = "kubectl create transient_configmap CONFIGMAP_NAME [--from-file=[key=]source] [--from-literal=key1=value1] [--from-env-file=[key=]source] --job-name=my-job [--from=cronjob/name]"

	// transientConfigMapShort is the short description of the transient_configmap command
	transientConfigMapShort = "Create a ConfigMap and a Job. And after the job is complete, delete them."

	// transientConfigMapLong is the long description of the transient_configmap command
	transientConfigMapLong = `Create a ConfigMap and a Job. And after the job is complete, delete them.`

	// transientConfigMapExample is the example of the transient_configmap command
	transientConfigMapExample = `
		# Create a new transient configmap named my-config based on folder bar
		kubectl create transient_configmap my-config --from-file=path/to/bar --job-name=test-job --job-from=cronjob/a-cronjob

		# Create a new transient configmap named my-config with specified keys instead of file basenames on disk
		kubectl create transient_configmap my-config --from-file=key1=/path/to/bar/file1.txt --from-file=key2=/path/to/bar/file2.txt --job-name=test-job --job-from=cronjob/a-cronjob

		# Create a new transient configmap named my-config with key1=config1 and key2=config2
		kubectl create transient_configmap my-config --from-literal=key1=config1 --from-literal=key2=config2 --job-name=test-job --job-from=cronjob/a-cronjob

		# Create a new transient configmap named my-config from the key=value pairs in the file
		kubectl create transient_configmap my-config --from-file=path/to/bar --job-name=test-job --job-from=cronjob/a-cronjob

		# Create a new transient configmap named my-config from an env file
		kubectl create transient_configmap my-config --from-env-file=path/to/foo.env --from-env-file=path/to/bar.env --job-name=test-job --job-from=cronjob/a-cronjob`
)

var (
	errConfigMapNameIsRequired = errors.New("configmap name is required")
)

// TransientConfigMapOptions provides information required to create a transient configmap
type TransientConfigMapOptions struct {
	command        []string
	name           string
	version        bool
	jobFlags       JobFlags
	configMapFlags ConfigMapFlags
}

// JobFlags contains the flags to create a job
type JobFlags struct {
	Name  string
	Image string
	From  string
}

// ConfigMapFlags contains the flags to create a configmap
type ConfigMapFlags struct {
	FileSources    []string
	LiteralSources []string
	EnvFileSources []string
	AppendHash     bool
}

// NewTransientConfigMapOptions provides an instance of TransientConfigMapOptions with default values
func NewTransientConfigMapOptions() *TransientConfigMapOptions {
	return &TransientConfigMapOptions{}
}

// NewCmdTransientConfigMap provides a cobra command wrapping TransientConfigMapOptions
func NewCmdTransientConfigMap(ioStreams genericiooptions.IOStreams) *cobra.Command {
	o := NewTransientConfigMapOptions()
	cmd := &cobra.Command{
		Use:                   transientConfigMapUse,
		DisableFlagsInUseLine: true,
		Short:                 transientConfigMapShort,
		Long:                  transientConfigMapLong,
		Example:               transientConfigMapExample,
		SilenceUsage:          true,
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "kubectl create transient_configmap",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run(cmd)
		},
	}
	cmd.SetOut(ioStreams.Out)
	cmd.SetErr(ioStreams.ErrOut)
	cmd.SetIn(ioStreams.In)

	cmd.Flags().BoolVar(&o.version, "version", false, "Print the version of this plugin")

	// job flags
	cmd.Flags().StringVar(&o.jobFlags.Name, "job-name", o.jobFlags.Name, "The name of the job to create.")
	cmd.Flags().StringVar(&o.jobFlags.Image, "job-image", o.jobFlags.Image, "Image name to run.")
	cmd.Flags().StringVar(&o.jobFlags.From, "job-from", o.jobFlags.From, "The name of the resource to create a Job from (only cronjob is supported).")

	// onceconfigmap flags
	cmd.Flags().StringSliceVar(&o.configMapFlags.FileSources, "from-file", o.configMapFlags.FileSources, "Key file can be specified using its file path, in which case file basename will be used as configmap key, or optionally with a key and file path, in which case the given key will be used.  Specifying a directory will iterate each named file in the directory whose basename is a valid configmap key.")
	cmd.Flags().StringArrayVar(&o.configMapFlags.LiteralSources, "from-literal", o.configMapFlags.LiteralSources, "Specify a key and literal value to insert in configmap (i.e. mykey=somevalue)")
	cmd.Flags().StringSliceVar(&o.configMapFlags.EnvFileSources, "from-env-file", o.configMapFlags.EnvFileSources, "Specify the path to a file to read lines of key=val pairs to create a configmap.")

	return cmd
}

// Complete completes all the required options
func (o *TransientConfigMapOptions) Complete(_ *cobra.Command, args []string) error {
	if o.version {
		return nil
	}
	if len(args) == 0 {
		return errConfigMapNameIsRequired
	}
	o.name = args[0]
	if len(args) > 1 {
		o.command = args[1:]
	}
	return nil
}

// Validate validates the provided options
func (o *TransientConfigMapOptions) Validate() error {
	// delegate validate processing to 'create configmap' and 'create job'
	return nil
}

// kubectl args
var (
	// createConfigMapArgs is the base args to create a configmap
	createConfigMapArgs = []string{"create", "configmap"}

	// createJobArgs is the base args to create a job
	createJobArgs = []string{"create", "job"}

	// waitJobArgs is the base args to wait for a job
	waitJobArgs = []string{"wait", "job"}

	// deleteJobArgs is the base args to delete a job
	deleteJobArgs = []string{"delete", "job"}

	// deleteConfigMapArgs is the base args to delete a configmap
	deleteConfigMapArgs = []string{"delete", "configmap"}
)

// Run executes the command
func (o *TransientConfigMapOptions) Run(cmd *cobra.Command) (err error) {
	if o.version {
		cmd.Printf("kubectl-create-transient_configmap %s\n", version)
		return
	}
	// Create the configmap
	cmd.Printf("creating configmap %s...\n", o.name)
	createConfigMapArgs = append(createConfigMapArgs, o.name)
	if len(o.configMapFlags.FileSources) != 0 {
		for _, source := range o.configMapFlags.FileSources {
			createConfigMapArgs = append(createConfigMapArgs, fmt.Sprintf("--from-file=%s", source))
		}
	}
	if len(o.configMapFlags.LiteralSources) != 0 {
		for _, source := range o.configMapFlags.LiteralSources {
			createConfigMapArgs = append(createConfigMapArgs, fmt.Sprintf("--from-literal=%s", source))
		}
	}
	if len(o.configMapFlags.EnvFileSources) != 0 {
		for _, source := range o.configMapFlags.EnvFileSources {
			createConfigMapArgs = append(createConfigMapArgs, fmt.Sprintf("--from-env-file=%s", source))
		}
	}
	out, err := executeKubectlCommand(cmd.Context(), createConfigMapArgs...)
	if err != nil {
		return
	}
	if len(out) > 0 {
		cmd.Printf("%s", out)
	}

	// Cleanup the configmap
	defer func() {
		cmd.Printf("deleting configmap %s...\n", o.name)
		deleteConfigMapArgs = append(deleteConfigMapArgs, o.name)
		out, dErr := executeKubectlCommand(cmd.Context(), deleteConfigMapArgs...)
		if dErr != nil {
			if err != nil {
				return
			}
			err = dErr
			return
		}
		if len(out) > 0 {
			cmd.Printf("%s", out)
		}
	}()

	// Create the job
	cmd.Printf("creating job %s...\n", o.jobFlags.Name)
	jobName := o.jobFlags.Name
	createJobArgs = append(createJobArgs, jobName)
	if len(o.jobFlags.From) != 0 {
		createJobArgs = append(createJobArgs, fmt.Sprintf("--from=%s", o.jobFlags.From))
	}
	if len(o.jobFlags.Image) != 0 {
		createJobArgs = append(createJobArgs, fmt.Sprintf("--image=%s", o.jobFlags.Image))
	}
	out, err = executeKubectlCommand(cmd.Context(), createJobArgs...)
	if err != nil {
		return
	}
	if len(out) > 0 {
		cmd.Printf("%s", out)
	}

	// Cleanup the job
	defer func() {
		cmd.Printf("deleting job %s...\n", jobName)
		deleteJobArgs = append(deleteJobArgs, jobName)
		out, dErr := executeKubectlCommand(cmd.Context(), deleteJobArgs...)
		if dErr != nil {
			if err != nil {
				return
			}
			err = dErr
			return
		}
		if len(out) > 0 {
			cmd.Printf("%s", out)
		}
	}()

	// Wait for the job to complete
	waitJobContext, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	complete := make(chan bool)
	go func() {
		waitJobCompleteArgs := append(waitJobArgs, jobName, "--for=condition=complete")
		_, wErr := executeKubectlCommand(waitJobContext, waitJobCompleteArgs...)
		select {
		case <-waitJobContext.Done():
			return
		default:
			if wErr != nil {
				if err == nil {
					err = wErr
				}
				complete <- false
				return
			}
		}
		complete <- true
	}()
	go func() {
		waitJobFailedArgs := append(waitJobArgs, jobName, "--for=condition=failed")
		_, wErr := executeKubectlCommand(waitJobContext, waitJobFailedArgs...)
		select {
		case <-waitJobContext.Done():
			return
		default:
			if wErr != nil && err == nil {
				err = wErr
			}
		}
		complete <- false
	}()

	success := <-complete
	if !success {
		err = fmt.Errorf(`job "%s" failed`, jobName)
		return
	}
	cmd.Printf("job \"%s\" completed\n", jobName)
	return nil
}

// executeKubectlCommand executes a kubectl command
func executeKubectlCommand(ctx context.Context, args ...string) ([]byte, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	kubectlCmd := exec.CommandContext(ctx, "kubectl", args...)
	kubectlCmd.Stdout = &stdout
	kubectlCmd.Stderr = &stderr
	err := kubectlCmd.Run()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		return nil, errors.New(strings.TrimSpace(strings.TrimPrefix(stderr.String(), "error:")))
	}
	return stdout.Bytes(), nil
}
