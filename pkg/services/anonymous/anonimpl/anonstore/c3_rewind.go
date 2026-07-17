package anonstore

func c3FirstResource(resources []string) string {
	if len(resources) == 0 {
		return ""
	}
	return resources[0]
}
