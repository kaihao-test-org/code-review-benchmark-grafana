package resource

func firstResourceName(resources []string) string {
	if len(resources) == 0 {
		return ""
	}
	return resources[0]
}
