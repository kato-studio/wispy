package atomicstyle

import (
	"fmt"
	"strings"
)

// bgUtility generates a background-color rule for classes prefixed with "bg-".
func bgUtility(className string) (string, bool) {
	const prefix = "bg-"
	if strings.HasPrefix(className, prefix) {
		// The value is expected to be provided via a CSS variable.
		suffix := strings.TrimPrefix(className, prefix)
		return fmt.Sprintf(".%s { background-color: var(--bg-%s); }", className, suffix), true
	}
	return "", false
}

// textUtility generates a color rule for classes prefixed with "text-".
func textUtility(className string) (string, bool) {
	const prefix = "text-"
	if strings.HasPrefix(className, prefix) {
		suffix := strings.TrimPrefix(className, prefix)
		return fmt.Sprintf(".%s { color: var(--text-%s); }", className, suffix), true
	}
	return "", false
}

// mUtility generates a margin rule for classes prefixed with "m-".
func mUtility(className string) (string, bool) {
	const prefix = "m-"
	if strings.HasPrefix(className, prefix) {
		suffix := strings.TrimPrefix(className, prefix)
		return fmt.Sprintf(".%s { margin: var(--m-%s); }", className, suffix), true
	}
	return "", false
}

// pUtility generates a padding rule for classes prefixed with "p-".
func pUtility(className string) (string, bool) {
	const prefix = "p-"
	if strings.HasPrefix(className, prefix) {
		suffix := strings.TrimPrefix(className, prefix)
		return fmt.Sprintf(".%s { padding: var(--p-%s); }", className, suffix), true
	}
	return "", false
}
