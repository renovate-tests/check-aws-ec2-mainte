package checkawsec2mainte

import (
	"context"
	"os"

	"github.com/mackerelio/checkers"
)

// Do ...
func Do() {
	var ckr *checkers.Checker
	defer func() {
		if ckr != nil {
			ckr.Name = "EC2 Mainte"
			ckr.Exit()
		}
	}()

	c, err := NewCli(os.Args)
	if err != nil {
		ckr = checkers.Unknown(err.Error())
		return
	}

	events, err := c.Fetch(context.Background())
	if err != nil {
		ckr = checkers.Unknown(err.Error())
		return
	}

	ckr = c.Evaluate(events)
}
