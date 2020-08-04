package main

import "os/exec"

func open(this string) {
	exec.Command("open", this).Run()
}
