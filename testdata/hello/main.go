// Command hello is what the release workflow is rehearsed on. It says hello, and with --version
// it says which release it was built as, which is all a rehearsal needs of a program.
package main

import (
	"fmt"
	"os"
)

// version is set when a release is built.
var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("hello", version)
		return
	}
	fmt.Println("hello")
}
