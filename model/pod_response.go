package model

// MiniPod is a Pod with only the information needed for the UI.
type MiniPod struct {
	MiniObjectMeta `json:"metadata,omitempty"`
}

// MiniPodList is a list of Pods with only the information needed for the UI.
type MiniPodList struct {
	ListMeta `json:"metadata,omitempty"`
	Items    []MiniPod `json:"items"`
}

// ConvertMini converts a PodList object into a MiniPodList object.
func (rc *PodList) ConvertMini() MiniPodList {
	return MiniPodList{
		ListMeta: ListMeta(rc.ListMeta),
		Items:    rc.convertPodsToMiniPods(),
	}
}

// convertPodsToMiniPods converts a slice of Pod objects into a slice of MiniPod objects.
func (rc *PodList) convertPodsToMiniPods() []MiniPod {
	pods := make([]MiniPod, len(rc.Items))
	for i, pod := range rc.Items {
		pods[i] = MiniPod{
			MiniObjectMeta: MiniObjectMeta{
				UID:               pod.UID,
				CreationTimestamp: pod.CreationTimestamp,
				Name:              pod.Name,
				GenerateName:      pod.GenerateName,
				Namespace:         pod.Namespace,
			},
		}
	}

	return pods
}
