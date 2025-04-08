package entrypoint

import (
	"goproject/internal/repl"
	"os"
)

func Entrypoint() {
	ret := repl.StartRepl()
	os.Exit(ret)
}
