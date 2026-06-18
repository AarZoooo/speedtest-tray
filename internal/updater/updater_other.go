//go:build !windows && !darwin

package updater

import (
	"fmt"
	"runtime"
)

func Apply(info UpdateInfo, onProgress func(percent int)) error {
	return fmt.Errorf("updates are not supported on %s", runtime.GOOS)
}
