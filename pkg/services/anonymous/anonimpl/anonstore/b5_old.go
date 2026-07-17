package anonstore

const b5Marker = "seed"

func b5FirstResource(resources []string) string {
	if len(resources) == 0 {
		return ""
	}
	return resources[0]
}
