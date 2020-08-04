package open

import (
	"os"
	"os/exec"
	"path/filepath"
)

var (
	cmd      = "url.dll,FileProtocolHandler"
	runDll32 = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
)

func open(this string) {
	exec.Command(runDll32, cmd, this).Run()
}
