package gorec

// This is useful for a simple conversion of url.Values (i.e., map[string][]string) into a format that is assignable on this system.
func MakeFormAssignable(origin map[string][]string) map[string]interface{} {
	newmap := map[string]interface{}{}
	for k, v := range origin {
		if len(v) == 0 {
			newmap[k] = nil
		} else {
			newmap[k] = v[0]
		}
	}
	return newmap
}
