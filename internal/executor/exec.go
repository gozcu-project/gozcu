package executor

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Run, onaylanan komutu exec() ile gerçek binary'ye geçirir.
// PID korunur, audit/process tree doğal görünür.
func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("çalıştırılacak komut yok")
	}
	binPath, err := exec.LookPath(args[0])
	if err != nil {
		return fmt.Errorf("komut bulunamadı: %s", args[0])
	}
	env := os.Environ()
	return syscall.Exec(binPath, args, env)
}