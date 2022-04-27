package executor

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
	"golang.org/x/xerrors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
)

// extracts the phase from the Kubernetes Job object
func getStatus(obj *corev1.Pod, labels labelSet) (status *v1.JobStatus, err error) {
	defer func() {
		if status != nil && status.Phase == v1.JobPhase_PHASE_DONE {
			status.Metadata.Finished = timestamppb.Now()
		}
	}()

	name, hasName := getJobName(obj, labels)
	if !hasName {
		return nil, xerrors.Errorf("job has no name: %v", obj.Name)
	}

	rawmd, ok := obj.Annotations[labels.AnnotationMetadata]
	if !ok {
		return nil, xerrors.Errorf("job has no metadata")
	}
	var md v1.JobMetadata
	unmarshaler := &protojson.UnmarshalOptions{}
	err = unmarshaler.Unmarshal([]byte(rawmd), &md)
	if err != nil {
		return nil, xerrors.Errorf("cannot unmarshal metadata %v :%w", rawmd, err)
	}

	var results []*v1.JobResult
	if c, ok := obj.Annotations[labels.AnnotationResults]; ok {
		err = json.Unmarshal([]byte(c), results)
		if err != nil {
			return nil, xerrors.Errorf("cannot unmarshal results: %w", err)
		}
	}

	annotationCanReplay := labels.AnnotationCanReplay
	_, canReplay := obj.Annotations[annotationCanReplay]

	annotationWaitUntil := labels.AnnotationWaitUntil
	var waitUntil *timestamppb.Timestamp
	if wt, ok := obj.Annotations[annotationWaitUntil]; ok {
		ts, err := time.Parse(time.RFC3339, wt)
		if err != nil {
			return nil, xerrors.Errorf("cannot parse %s annotation: %w", annotationWaitUntil, err)
		}
		waitUntil = timestamppb.New(ts)
	}

	status = &v1.JobStatus{
		Name:     name,
		Metadata: &md,
		Phase:    v1.JobPhase_PHASE_UNKNOWN,
		Conditions: &v1.JobConditions{
			Success:   true,
			CanReplay: canReplay,
			WaitUntil: waitUntil,
		},
		Results: results,
	}

	var (
		statuses      = append(obj.Status.InitContainerStatuses, obj.Status.ContainerStatuses...)
		anyFailed     bool
		maxRestart    int32
		allTerminated = len(statuses) != 0
	)
	for _, cs := range statuses {
		if w := cs.State.Waiting; w != nil && w.Reason == "ErrImagePull" {
			status.Phase = v1.JobPhase_PHASE_DONE
			status.Conditions.Success = false
			status.Details = w.Message
			return
		}

		isSidecarContainer := strings.Contains(obj.Annotations[labels.AnnotationSidecars], cs.Name)
		if cs.State.Terminated != nil {
			if cs.State.Terminated.ExitCode != 0 {
				anyFailed = true
			}
		} else if !isSidecarContainer {
			allTerminated = false
		}

		if cs.RestartCount >= maxRestart {
			maxRestart = cs.RestartCount
		}
	}
	status.Conditions.FailureCount = maxRestart
	status.Conditions.Success = !(anyFailed || maxRestart > getFailureLimit(obj, labels))
	status.Conditions.DidExecute = obj.Status.Phase != "" || len(statuses) > 0

	if msg, failed := obj.Annotations[labels.AnnotationFailed]; failed {
		status.Phase = v1.JobPhase_PHASE_DONE
		if obj.DeletionTimestamp != nil {
			status.Phase = v1.JobPhase_PHASE_CLEANUP
		}
		status.Conditions.Success = false
		status.Details = msg

		return
	}
	if obj.DeletionTimestamp != nil {
		status.Phase = v1.JobPhase_PHASE_CLEANUP
		return
	}
	if maxRestart > getFailureLimit(obj, labels) {
		status.Phase = v1.JobPhase_PHASE_DONE
		return
	}
	if allTerminated {
		status.Phase = v1.JobPhase_PHASE_DONE
		return
	}

	switch obj.Status.Phase {
	case corev1.PodPending:
		status.Phase = v1.JobPhase_PHASE_PREPARING
		return
	case corev1.PodRunning:
		status.Phase = v1.JobPhase_PHASE_RUNNING
	}

	return
}

func getFailureLimit(obj *corev1.Pod, labels labelSet) int32 {
	val := obj.Annotations[labels.AnnotationFailureLimit]
	if val == "" {
		val = "0"
	}

	res, _ := strconv.ParseInt(val, 10, 32)
	return int32(res)
}

func getJobName(obj *corev1.Pod, labels labelSet) (id string, ok bool) {
	id, ok = obj.Labels[labels.LabelJobName]
	return
}
