package assert

import (
	"fmt"
	"os"
)

func Assert(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: (%s)\n", msg, err.Error())
		os.Exit(1)
	}
}
