package duplicate

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MetadataSpecExtractor[T any] interface {
	GetPodMetadata() *metav1.ObjectMeta
	GetPodSpec() *corev1.PodSpec
}

func extractPod[T MetadataSpecExtractor[T]](obj T) (*metav1.ObjectMeta, *corev1.PodSpec) {
	return obj.GetPodMetadata(), obj.GetPodSpec()
}

type StatefulSetAdapter struct {
	*appsv1.StatefulSet
}

func (s StatefulSetAdapter) GetPodMetadata() *metav1.ObjectMeta {
	return &s.StatefulSet.Spec.Template.ObjectMeta
}

func (s StatefulSetAdapter) GetPodSpec() *corev1.PodSpec {
	return &s.StatefulSet.Spec.Template.Spec
}

type DeploymentAdapter struct {
	*appsv1.Deployment
}

func (s DeploymentAdapter) GetPodMetadata() *metav1.ObjectMeta {
	return &s.Deployment.Spec.Template.ObjectMeta
}

func (s DeploymentAdapter) GetPodSpec() *corev1.PodSpec {
	return &s.Deployment.Spec.Template.Spec
}

type CronJobAdapter struct {
	*batchv1.CronJob
}

func (s CronJobAdapter) GetPodMetadata() *metav1.ObjectMeta {
	return &s.CronJob.Spec.JobTemplate.Spec.Template.ObjectMeta
}

func (s CronJobAdapter) GetPodSpec() *corev1.PodSpec {
	return &s.CronJob.Spec.JobTemplate.Spec.Template.Spec
}

type JobAdapter struct {
	*batchv1.Job
}

func (s JobAdapter) GetPodMetadata() *metav1.ObjectMeta {
	return &s.Job.Spec.Template.ObjectMeta
}

func (s JobAdapter) GetPodSpec() *corev1.PodSpec {
	return &s.Job.Spec.Template.Spec
}

type PodAdapter struct {
	*corev1.Pod
}

func (s PodAdapter) GetPodMetadata() *metav1.ObjectMeta {
	return &s.Pod.ObjectMeta
}

func (s PodAdapter) GetPodSpec() *corev1.PodSpec {
	return &s.Pod.Spec
}
