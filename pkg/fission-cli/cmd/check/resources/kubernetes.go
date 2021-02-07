package resources

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	//KubernetesSupportedVersion is Kubernetes version supported by fission
	KubernetesSupportedVersion = "v1.9"
	//KubernetesVersionLabel is label for KubernetesVersion check
	KubernetesVersionLabel = "kubernetes version compatibility"
	//KubernetesPodStatusLabel is label for KubernetesPodStatus check
	KubernetesPodStatusLabel = "determine pods are running"
)

// KubernetesVersion is used to check Kubernetes version compatibility with fission
type KubernetesVersion struct {
	client *kubernetes.Clientset
}

// NewKubernetesVersion is used to create new KubernetesVersion instance
func NewKubernetesVersion(clientset *kubernetes.Clientset) Resource {
	return KubernetesVersion{client: clientset}
}

// Check performs the check and returns the result
func (res KubernetesVersion) Check() Results {
	serverVer, err := res.client.ServerVersion()
	if err != nil {
		return getResults(err.Error(), false)
	}

	if semver.Compare(serverVer.GitVersion, KubernetesSupportedVersion) == -1 {
		return getResults("kubernetes version is incompatible", false)
	}

	return getResults("kubernetes version is compatible", true)
}

// GetLabel returns the label for check
func (res KubernetesVersion) GetLabel() string {
	return KubernetesVersionLabel
}

// KubernetesPodStatus is used to check selected pods are running
type KubernetesPodStatus struct {
	client *kubernetes.Clientset
}

// NewKubernetesPodStatus is used to create new KubernetesPodStatus
func NewKubernetesPodStatus(clientset *kubernetes.Clientset) Resource {
	return KubernetesPodStatus{client: clientset}
}

// Check performs the check and returns the result
func (res KubernetesPodStatus) Check() Results {
	podsToVerify := []string{
		"buildermgr",
		"controller",
		"executor",
		"influxdb",
		"kubewatcher",
		"logger",
		"mqtrigger",
		"nats-streaming",
		"router",
		"storagesvc",
		"timer",
	}
	podsToVerifyMap := make(map[string]bool, len(podsToVerify))
	for _, idx := range podsToVerify {
		podsToVerifyMap[idx] = false
	}

	objs, err := res.client.CoreV1().Pods(metav1.NamespaceAll).List(
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("svc in (%s)", strings.Join(podsToVerify, ", ")),
		},
	)
	if err != nil {
		return getResults(err.Error(), false)
	}

	var results Results
	for _, item := range objs.Items {
		if item.Status.Phase != corev1.PodRunning {
			results = appendResult(results, fmt.Sprintf("pod %s is not running", item.Name), false)
		} else {
			results = appendResult(results, fmt.Sprintf("pod %s is running", item.Name), true)
		}
		val, ok := item.Labels["svc"]
		if ok {
			podsToVerifyMap[val] = true
		}
	}

	for key, val := range podsToVerifyMap {
		if !val {
			results = appendResult(results, fmt.Sprintf("not able to find pod with label svc=%s", key), false)
		}
	}

	return results
}

// GetLabel returns the label for check
func (res KubernetesPodStatus) GetLabel() string {
	return KubernetesPodStatusLabel
}
