//go:build !windows
// +build !windows

package engine

// check if a file is hidden on unix
func isHidden(path string) (bool, error) {
	return path[0] == '.', nil
}
