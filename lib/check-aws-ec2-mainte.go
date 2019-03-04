package checkawsec2mainte

import (
	"os"
	"fmt"

	"github.com/mackerelio/checkers"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "EC2 Mainte"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {

	mt, err := NewEC2Mainte(ec2.New(session.New()))
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	if mt.Length() != 0 {
		event := mt.GetCloseEvent()
		return checkers.Warning(fmt.Sprintf("%+v", event))
	}

	return checkers.Ok("")
}
