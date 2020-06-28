package main

// Import declaration declares library packages referenced in this file.
import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

func init() {
	// Adds a hook to setu up container namespaces
	reexec.Register("namespaceSetup", namespaceSetup)
	if reexec.Init() {
		os.Exit(0)
	}
}

func namespaceSetup() {
	fmt.Println("hi")
	args := os.Args[1:]
	usercmd := args[0]

	cmd := exec.Command(usercmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{"PS1=containerville # "}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the %s command - %s\n", usercmd, err)
		os.Exit(1)
	}
}

func main() {
	cmd := reexec.Command("namespaceSetup")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

}
