package checkawsec2mainte

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mackerelio/checkers"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	region       = kingpin.Flag("region", "").Default("ap-northeast-1").String()
	warnDuration = kingpin.Flag("warning-duration", "").Short('w').Default("240h30m").Duration()
	critDuration = kingpin.Flag("critical-duration", "").Short('c').Default("120h30s").Duration()
	instanceIds = kingpin.Flag("instance-ids", "").Short('i').Strings()
)

func init() {
	kingpin.Parse()
}

func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "EC2 Mainte"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {

	mt, err := NewEC2Mainte(ec2.New(session.New()), *instanceIds...)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	if mt.Length() != 0 {
		event := mt.GetCloseEvent()
	
		return checkers.Warning(fmt.Sprintf("%+v", event))
	}

	return checkers.Ok("")
}
