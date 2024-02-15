package duplicate

import (
	"fmt"

	duputil "dup/pkg/util"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/resource"
)

type DupPodSpec struct {
	Name    string
	Spec    *corev1.PodSpec
	Labels  map[string]string
	Options *DupPodOptions
}

type DupPodOptions struct {
	DuplicatePod  bool
	DisableProbes bool
	Image         string
	Command       string
}

func Clone(opts *DupPodOptions, objects []*resource.Info) ([]*runtime.Object, error) {
	var ret []*runtime.Object
	for i := range objects {
		obj := objects[i]
		objKind := obj.Object.GetObjectKind().GroupVersionKind().Kind
		if opts.DuplicatePod && hasPodSpec(objKind) {
			dupPodSpec, err := generateDupPodSpec(obj)
			if err != nil {
				return nil, err
			}
			dupPodSpec.Options = opts
			dupPod := applyDupPodSpec(dupPodSpec)
			ret = append(ret, &dupPod)
			continue
		}
		clonedResource, err := cloneResource(obj.Object)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &clonedResource)
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

func generateDupPodSpec(obj *resource.Info) (DupPodSpec, error) {
	objType := obj.Object.GetObjectKind().GroupVersionKind().Kind
	switch objType {
	case "StatefulSet":
		sts := &appsv1.StatefulSet{}
		unstructuredToType(obj.Object, sts)
		return DupPodSpec{
			Name:   duputil.GenerateResourceName(sts.Name),
			Spec:   &sts.Spec.Template.Spec,
			Labels: DeepCopyLabels(sts.Spec.Template.GetLabels()),
		}, nil
	case "Deployment":
		deploy := &appsv1.Deployment{}
		unstructuredToType(obj.Object, deploy)
		return DupPodSpec{
			Name:   duputil.GenerateResourceName(deploy.Name),
			Spec:   &deploy.Spec.Template.Spec,
			Labels: DeepCopyLabels(deploy.Spec.Template.GetLabels()),
		}, nil
	case "CronJob":
		cron := &batchv1.CronJob{}
		unstructuredToType(obj.Object, cron)
		return DupPodSpec{
			Name:   duputil.GenerateResourceName(cron.Name),
			Spec:   &cron.Spec.JobTemplate.Spec.Template.Spec,
			Labels: DeepCopyLabels(cron.Spec.JobTemplate.Spec.Template.GetLabels()),
		}, nil
	case "Job":
		job := &batchv1.Job{}
		unstructuredToType(obj.Object, job)
		return DupPodSpec{
			Name:   duputil.GenerateResourceName(job.Name),
			Spec:   &job.Spec.Template.Spec,
			Labels: DeepCopyLabels(job.Spec.Template.GetLabels()),
		}, nil
	case "Pod":
		pod := &corev1.Pod{}
		unstructuredToType(obj.Object, pod)
		return DupPodSpec{
			Name:   duputil.GenerateResourceName(pod.Name),
			Spec:   &pod.Spec,
			Labels: DeepCopyLabels(pod.ObjectMeta.GetLabels()),
		}, nil
	default:
		return DupPodSpec{}, fmt.Errorf("object type %s does not have PodSpec or is not supported", objType)
	}
}

func DeepCopyLabels(labels map[string]string) map[string]string {
	labelMap := make(map[string]string)

	for key, value := range labels {
		labelMap[key] = value
	}
	return labelMap
}

func cloneResource(obj runtime.Object) (runtime.Object, error) {
	objCopy := obj.DeepCopyObject()
	accessor := meta.NewAccessor()
	name, err := accessor.Name(objCopy)
	if err != nil {
		return nil, err
	}
	suffixedName := duputil.GenerateResourceName(name)
	accessor.SetName(objCopy, suffixedName)

	return objCopy, nil
}
func applyDupPodSpec(pod DupPodSpec) runtime.Object {
	var clonedPod corev1.Pod

	pod.Spec.DeepCopyInto(&clonedPod.Spec)
	clonedPod.ObjectMeta.Labels = pod.Labels
	clonedPod.Name = pod.Name

	if pod.Options != nil && pod.Options.DisableProbes {
		disableProbes(&clonedPod)
	}

	return &clonedPod
}

func disableProbes(pod *corev1.Pod) {
	containers := pod.Spec.Containers
	for i := range containers {
		containers[i].ReadinessProbe = nil
		containers[i].LivenessProbe = nil
	}
}
