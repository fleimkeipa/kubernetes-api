package model

// MiniNamespace is a Namespace with only the information needed for the UI.
type MiniNamespace struct {
	MiniObjectMeta `json:"metadata,omitempty"`
}

// MiniNamespaceList is a list of Namespaces with only the information needed for the UI.
type MiniNamespaceList struct {
	ListMeta `json:"metadata,omitempty"`
	Items    []MiniNamespace `json:"items"`
}

// ConvertMini converts a NamespaceList object into a MiniNamespaceList object.
func (rc *NamespaceList) ConvertMini() MiniNamespaceList {
	return MiniNamespaceList{
		ListMeta: ListMeta(rc.ListMeta),
		Items:    rc.convertNamespacesToMini(),
	}
}

// convertNamespacesToMini converts a slice of Namespace objects into a slice of MiniNamespace objects.
func (rc *NamespaceList) convertNamespacesToMini() []MiniNamespace {
	namespaces := make([]MiniNamespace, len(rc.Items))
	for i, namespace := range rc.Items {
		namespaces[i] = MiniNamespace{
			MiniObjectMeta: MiniObjectMeta{
				UID:               namespace.UID,
				CreationTimestamp: namespace.CreationTimestamp,
				Name:              namespace.Name,
				GenerateName:      namespace.GenerateName,
				Namespace:         namespace.Namespace,
			},
		}
	}

	return namespaces
}
