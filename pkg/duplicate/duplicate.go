package duplicate

import (
	duputil "dup/pkg/util"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/resource"
)

const LOOP_COMMAND = "tail -f /dev/null"

type PodOptions struct {
	DuplicateInnerPod bool
	DisableProbes     bool
	LoopCommand       bool
	Image             string
}

func Clone(opts *PodOptions, objects []*resource.Info) ([]*runtime.Object, error) {
	var ret []*runtime.Object
	for i := range objects {
		obj := objects[i]
		objKind := obj.Object.GetObjectKind().GroupVersionKind().Kind
		if hasPodSpec(objKind) {
			dResource, err := cloneResourceWithPod(obj.Object, opts)
			if err != nil {
				return nil, err
			}
			ret = append(ret, dResource)
		} else {
			dResource, err := cloneGenericResource(obj.Object)
			if err != nil {
				return nil, err
			}
			ret = append(ret, dResource)
		}
	}
	return ret, nil
}

func hasPodSpec(kind string) bool {
	switch kind {
	case "Deployment", "CronJob", "StatefulSet", "Job", "Pod":
		return true
	}
	return false
}

func unstructuredToType(in runtime.Object, out runtime.Object) error {
	unstructured, ok := in.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("Error converting runtime.Object to unstructured")
	}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.Object, out)
	if err != nil {
		return fmt.Errorf("Error converting unstructured to Deploy: %v\n", err)
	}
	return nil
}

func cloneResourceWithPod(obj runtime.Object, opts *PodOptions) (*runtime.Object, error) {
	var dupObject runtime.Object
	var metadata *metav1.ObjectMeta
	var spec *corev1.PodSpec
	objType := obj.GetObjectKind().GroupVersionKind().Kind
	switch objType {
	case "StatefulSet":
		dupObject = &appsv1.StatefulSet{}
		unstructuredToType(obj, dupObject)
		metadata, spec = extractPod[StatefulSetAdapter](StatefulSetAdapter{dupObject.(*appsv1.StatefulSet)})
	case "Deployment":
		dupObject = &appsv1.Deployment{}
		unstructuredToType(obj, dupObject)
		metadata, spec = extractPod[DeploymentAdapter](DeploymentAdapter{dupObject.(*appsv1.Deployment)})
	case "CronJob":
		dupObject = &batchv1.CronJob{}
		unstructuredToType(obj, dupObject)
		metadata, spec = extractPod[CronJobAdapter](CronJobAdapter{dupObject.(*batchv1.CronJob)})
	case "Job":
		dupObject = &batchv1.Job{}
		unstructuredToType(obj, dupObject)
		metadata, spec = extractPod[JobAdapter](JobAdapter{dupObject.(*batchv1.Job)})
	case "Pod":
		dupObject = &corev1.Pod{}
		unstructuredToType(obj, dupObject)
		metadata, spec = extractPod[PodAdapter](PodAdapter{dupObject.(*corev1.Pod)})
	default:
		return nil, fmt.Errorf("object type %s does not have PodSpec or is not supported", objType)
	}

	applyOptions(objType, spec, metadata, opts)
	err := setSuffixedName(&dupObject)
	if err != nil {
		return nil, err
	}
	return &dupObject, nil
}

func setSuffixedName(obj *runtime.Object) error {
	accessor := meta.NewAccessor()
	name, err := accessor.Name(*obj)
	if err != nil {
		return err
	}
	suffixedName := duputil.GenerateResourceName(name)
	accessor.SetName(*obj, suffixedName)
	return nil
}

func cloneGenericResource(obj runtime.Object) (*runtime.Object, error) {
	objCopy := obj.DeepCopyObject()
	err := setSuffixedName(&objCopy)
	if err != nil {
		return nil, err
	}

	return &objCopy, nil
}
func applyOptions(kind string, spec *corev1.PodSpec, meta *metav1.ObjectMeta, opts *PodOptions) {
	if opts != nil {
		if opts.DisableProbes {
			disableProbes(spec)
		}
		if opts.LoopCommand {
			setCommand(spec)
		}
		if kind == "Pod" {
			removeOwnership(meta)
		}
	}
}

func disableProbes(podSpec *corev1.PodSpec) {
	containers := podSpec.Containers
	for i := range containers {
		containers[i].ReadinessProbe = nil
		containers[i].LivenessProbe = nil
	}
}

func setCommand(podSpec *corev1.PodSpec) {
	command := strings.Split(LOOP_COMMAND, " ")
	containers := podSpec.Containers
	for i := range containers {
		containers[i].Command = command
	}
}

func removeOwnership(metadata *metav1.ObjectMeta) {
	delete(metadata.Labels, "app.kubernetes.io/instance")
	delete(metadata.Labels, "app.kubernetes.io/name")
	metadata.OwnerReferences = nil
	metadata.UID = ""
	metadata.ResourceVersion = ""
}
