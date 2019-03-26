package metadata

import (
	"encoding/json"
	"time"
)

const MetaMainteTimeLayout = "2 Jan 2006 15:04:05 MST"

type MetaMainteTime time.Time

func (ct *MetaMainteTime) UnmarshalJSON(b []byte) error {
	tm, err := time.Parse(`"`+MetaMainteTimeLayout+`"`, string(b))
	*ct = MetaMainteTime(tm)
	return err
}

func (ct MetaMainteTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(ct).Format(MetaMainteTimeLayout))
}
