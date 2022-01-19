package repoconfig

import (
	pirov1 "github.com/bhojpur/piro/pkg/api/v1"
	"github.com/bhojpur/piro/pkg/filterexpr"
	corev1 "k8s.io/api/core/v1"
)

// C is the struct we expect to find in the repo root which configures how we build things
type C struct {
	DefaultJob string          `yaml:"defaultJob"`
	Rules      []*JobStartRule `yaml:"rules"`
}

// JobStartRule determines if a Kubernetes Job will be started
type JobStartRule struct {
	Path string                     `yaml:"path"`
	Expr []*pirov1.FilterExpression `yaml:"matchesAll"`
}

// UnmarshalYAML unmarshals the filter expressions
func (r *JobStartRule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawJobStartRule struct {
		Path string           `yaml:"path"`
		Expr []JobStartRuleOr `yaml:"matchesAll"`
	}
	err := unmarshal(&rawJobStartRule)
	if err != nil {
		return err
	}

	r.Path = rawJobStartRule.Path
	for _, expr := range rawJobStartRule.Expr {
		terms, err := filterexpr.Parse(expr.Or)
		if err != nil {
			return err
		}
		r.Expr = append(r.Expr, &pirov1.FilterExpression{Terms: terms})
	}

	return nil
}

// JobStartRuleOr contains an "OR'ed" list of conditions which have to match for a Job to run
type JobStartRuleOr struct {
	Or []string `yaml:"or"`
}

// TemplatePath returns the path to the Job template in the repo
func (rc *C) TemplatePath(md *pirov1.JobMetadata) string {
	js := &pirov1.JobStatus{Metadata: md}
	for _, rule := range rc.Rules {
		if filterexpr.MatchesFilter(js, rule.Expr) {
			return rule.Path
		}
	}

	return rc.DefaultJob
}

// ShouldRun determines based on the repo config if the Job should run
func (rc *C) ShouldRun(md *pirov1.JobMetadata) bool {
	return rc.TemplatePath(md) != ""
}

// JobSpec is the format of the files we expect to find when starting jobs
type JobSpec struct {
	// Desc describes the purpose of this Kubernetes Job spec.
	Desc string `yaml:"description,omitempty"`

	// Pod is the actual Job spec to start. Prior to deploying this to
	// Kubernetes, we'll run this as a Go template.
	Pod *corev1.PodSpec `yaml:"pod"`

	// Mutex makes Job execution exclusive, with new ones canceling the
	// currently running one. For example: job A is running at the moment,
	// and job B is about to start. If A and B share the same mutex, B will
	// cancel A.
	Mutex string `yaml:"mutex,omitempty"`

	// Args describe annotations which this Job expects. This list is only
	// used on the UI when manually starting the Kubernetes Job.
	// This list is neither exhaustive (i.e. Jobs can use annotations not
	// listed here), nor binding (i.e. Jobs can run even when annotations
	// listed here are not present). What matters for a Job to run is only
	// if Kubernetes accepts the produced podspec.
	Args []ArgSpec `yaml:"args,omitempty"`

	Sidecars []string `yaml:"sidecars,omitempty"`
}

// ArgSpec specifies an argument/annotation for a job.
type ArgSpec struct {
	Name string `yaml:"name"`
	Req  bool   `yaml:"required"`
	Desc string `yaml:"description"`
}
