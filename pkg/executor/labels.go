package executor

import "strings"

const (
	// DefaultLabelPrefix is used when no explicit label prefix is set
	defaultLabelPrefix = "piro.bhojpur.net"
)

type labelSet struct {
	// LabelPiroMarker is the label applied to all jobs and configmaps. This label can be used
	// to search for Bhojpur Piro job objects in Kubernetes.
	LabelPiroMarker string

	// LabelJobName adds the ID of the job to the k8s object
	LabelJobName string

	// LabelMutex makes jobs findable via their mutex
	LabelMutex string

	// UserDataAnnotationPrefix is prepended together with the label prefix to all user annotations added to jobs
	UserDataAnnotationPrefix string

	// AnnotationFailureLimit is the annotation denoting the max times a job may fail
	AnnotationFailureLimit string

	// AnnotationMetadata stores the JSON encoded metadata available at creation
	AnnotationMetadata string

	// AnnotationFailed explicitelly fails the job
	AnnotationFailed string

	// AnnotationResults stores JSON encoded list of a job results
	AnnotationResults string

	// AnnotationCanReplay stores if this job can be replayed
	AnnotationCanReplay string

	// AnnotationWaitUntil stores the start time of waiting job
	AnnotationWaitUntil string

	// AnnotationSidecars lists all container whose lifecycle depends on that of the others
	AnnotationSidecars string
}

// newLabelSetet returns a new label set initialized with a particular prefix
func newLabelSetet(prefix string) labelSet {
	if prefix == "" {
		prefix = defaultLabelPrefix
	}
	prefix = strings.TrimSuffix(prefix, "/") + "/"

	return labelSet{
		LabelPiroMarker:          prefix + "job",
		LabelJobName:             prefix + "jobName",
		LabelMutex:               prefix + "mutex",
		UserDataAnnotationPrefix: "userdata." + prefix,
		AnnotationFailureLimit:   prefix + "failureLimit",
		AnnotationMetadata:       prefix + "metadata",
		AnnotationFailed:         prefix + "failed",
		AnnotationResults:        prefix + "results",
		AnnotationCanReplay:      prefix + "canReplay",
		AnnotationWaitUntil:      prefix + "waitUntil",
		AnnotationSidecars:       prefix + "sidecars",
	}
}
