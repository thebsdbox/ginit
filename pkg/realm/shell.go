package realm

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Shell will Start a userland shell
func Shell() {
	// Shell stuff
	log.Println("Starting Shell")

	// TTY hack to support ctrl+c
	cmd := exec.Command("/usr/bin/setsid", "cttyhack", "/bin/sh")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Errorf("Shell error [%v]", err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	if err != nil {
		log.Errorf("Shell error [%v]", err)
	}
}
