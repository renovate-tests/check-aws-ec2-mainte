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
	version    = "indev"
	goversion  = ""
	commitHash = ""
	buildDate  = ""
)

var (
	app          = kingpin.New("check-aws-ec2-mainte", fmt.Sprintf("GoVer: %v\tCommitHash: %v\tBuildDate: %v", goversion, commitHash, buildDate)).Version(version).Author("ntrv")
	region       = app.Flag("region", "").Default("ap-northeast-1").String()
	warnDuration = app.Flag("warning-duration", "").Short('w').PlaceHolder("1h23m4s").Default("240h").Duration()
	critDuration = app.Flag("critical-duration", "").Short('c').PlaceHolder("5h56m7s").Default("120h").Duration()
	instanceIds  = app.Flag("instance-ids", "Available to specify multiple time").Short('i').PlaceHolder("i-0f456b937f33abe9e").Strings()
)

func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "EC2 Mainte"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {

	_, err := app.Parse(args)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	sess := session.New()
	mt, err := NewEC2Mainte(ec2.New(sess), *instanceIds...)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	if mt.Length() != 0 {
		event := mt.GetCloseEvent()
		return checkers.Warning(fmt.Sprintf("%+v", event))
	}

	return checkers.Ok("")
}
