package main

import (
	"context"
	"os"

	"k8s.io/component-base/logs"
	schedulerapp "k8s.io/kubernetes/cmd/kube-scheduler/app"
	scheduleroptions "k8s.io/kubernetes/cmd/kube-scheduler/app/options"

	"example.com/custom-scheduler/plugin"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	ctx := context.Background()
	opts := scheduleroptions.NewOptions()

	// Allow passing --config=/config/scheduler-config.yaml, etc.
	command := schedulerapp.NewSchedulerCommand(
		ctx,
		schedulerapp.WithPlugin(plugin.Name, plugin.New),
		schedulerapp.WithPluginConfig(plugin.Name, plugin.NewConfig()),
	)

	// Add standard flags (so --config works)
	opts.AddFlags(command.Flags())

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
