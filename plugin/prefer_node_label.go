package plugin

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	schedconfig "k8s.io/kube-scheduler/config/v1"
	framework "k8s.io/kube-scheduler/framework"
)

const (
	Name = "PreferNodeLabel"
)

// Args are configured via KubeSchedulerConfiguration (pluginConfig.args)
type Args struct {
	// LabelKey and LabelValue to prefer on nodes, e.g. workload=true
	LabelKey   string `json:"labelKey,omitempty"`
	LabelValue string `json:"labelValue,omitempty"`
}

// PreferNodeLabelPlugin implements Score and (optionally) Filter.
type PreferNodeLabelPlugin struct {
	handle framework.Handle
	args   *Args
}

// Verify the interfaces we implement.
var (
	_ framework.ScorePlugin     = &PreferNodeLabelPlugin{}
	_ framework.ScoreExtensions = &PreferNodeLabelPlugin{}
	_ framework.FilterPlugin    = &PreferNodeLabelPlugin{}
)

func (p *PreferNodeLabelPlugin) Name() string { return Name }

// New is called by the framework to create the plugin.
func New(_ context.Context, obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	args := &Args{}
	if obj != nil {
		if err := framework.DecodeInto(obj, args); err != nil {
			return nil, fmt.Errorf("decoding args: %w", err)
		}
	}
	if args.LabelKey == "" {
		args.LabelKey = "workload"
	}
	if args.LabelValue == "" {
		args.LabelValue = "true"
	}
	return &PreferNodeLabelPlugin{handle: handle, args: args}, nil
}

// Filter: (optional) reject nodes that are unschedulable for other reasons.
// Here we keep it permissive: everyone passes. You could enforce the label here instead.
func (p *PreferNodeLabelPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	// Example (strict): require the label
	// if nodeInfo.Node().Labels[p.args.LabelKey] != p.args.LabelValue {
	//   return framework.NewStatus(framework.Unschedulable, "missing preferred label")
	// }
	return framework.NewStatus(framework.Success, "")
}

// Score gives higher score to nodes with the preferred label.
func (p *PreferNodeLabelPlugin) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	ni, err := p.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.AsStatus(err)
	}
	node := ni.Node()
	if node == nil {
		return 0, framework.NewStatus(framework.Error, "node not found")
	}
	if node.Labels[p.args.LabelKey] == p.args.LabelValue {
		// Raw score before normalization
		return framework.MaxNodeScore, framework.NewStatus(framework.Success, "")
	}
	return 0, framework.NewStatus(framework.Success, "")
}

// ScoreExtensions allows us to normalize scores (optional here).
func (p *PreferNodeLabelPlugin) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	// Already using 0 / MaxNodeScore; nothing to normalize.
	return framework.NewStatus(framework.Success, "")
}

// Helper for v1 config registration
func NewConfig() runtime.Object { return &schedconfig.PluginConfig{} }
