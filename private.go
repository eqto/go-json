package json

func getFromMap(obj map[string]interface{}, paths ...string) interface{} {
	if obj == nil {
		return nil
	}
	if val, ok := obj[paths[0]]; ok {
		if len(paths) == 1 {
			return val
		}
		switch obj := val.(type) {
		case map[string]interface{}:
			return getFromMap(obj, paths[1:]...)
		case Object:
			return getFromMap(obj, paths[1:]...)
		}
	}
	return nil
}
