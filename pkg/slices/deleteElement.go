package slices

func DeleteElement(a []string, i string) []string {
	var j int
	for _, v := range a {
		if v != i {
			a[j] = v
			j++
		}
	}
	return a[:j]
}
