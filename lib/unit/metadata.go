package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	checkawsec2mainte "github.com/ntrv/check-aws-ec2-mainte/lib"
)

func StartTestServer(patterns map[string]string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if resp, ok := patterns[r.RequestURI]; ok {
				w.Write([]byte(resp))
				return
			}
			http.Error(w, "not found", http.StatusNotFound)
			return
		}),
	)
}

func CreateEventsMetadata(t *testing.T, instanceId string) checkawsec2mainte.EC2Events {
	ds := CreateTimes(t, []string{
		"2019-03-14T16:04:05+09:00",
		"2019-03-16T16:04:05+09:00",
		"2019-03-16T18:04:05+09:00",
		"2019-03-17T17:34:35+09:00",
	})

	return checkawsec2mainte.EC2Events{
		{
			Code:        ec2.EventCodeSystemReboot,
			InstanceId:  instanceId,
			NotBefore:   ds[2],
			Description: "scheduled reboot",
			State:       checkawsec2mainte.StateActive,
		},
		{
			Code:        ec2.EventCodeSystemMaintenance,
			InstanceId:  instanceId,
			NotBefore:   ds[0], // Closest Event
			NotAfter:    ds[1],
			Description: "[Completed] Scheduled System Maintenance",
			State:       checkawsec2mainte.StateCompleted,
		},
		{
			Code:        ec2.EventCodeInstanceRetirement,
			InstanceId:  instanceId,
			NotBefore:   ds[3],
			Description: "[Canceled] Scheduled Instance Retirement Maintenance",
			State:       checkawsec2mainte.StateCanceled,
		},
	}
}
