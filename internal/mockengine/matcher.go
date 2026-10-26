package mockengine

import "strings"

func matchPath(pattern, path string) bool {
	patternParts := splitPath(pattern)
	pathParts := splitPath(path)

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, patternPart := range patternParts {
		if strings.HasPrefix(patternPart, ":") {
			if len(patternPart) == 1 || pathParts[i] == "" {
				return false
			}

			continue
		}

		if patternPart != pathParts[i] {
			return false
		}
	}

	return true

}

func splitPath(path string) []string {
	if path == "/" {
		return []string{}
	}

	path = strings.Trim(path, "/")

	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}
