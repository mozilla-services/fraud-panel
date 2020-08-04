package open

import "os/exec"

func open(this string) {
	exec.Command("xdg-open", this).Run()
}
