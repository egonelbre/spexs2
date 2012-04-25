package debugger

import (
	"bufio"
	"os"
	"strings"
)

func RunShell(d *Debugger) {
	r := bufio.NewReader(os.Stdin)

	readLine := func() string {
		ret, err := r.ReadString('\n')
		if err != nil {
			os.Exit(1)
		}
		return ret
	}

	go func() {
		for {
			cmd := readLine()
			cmd = strings.TrimSpace(cmd)
			d.Commands <- cmd
		}
	}()
}
