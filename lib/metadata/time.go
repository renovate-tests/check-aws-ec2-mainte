package metadata

import (
	"encoding/json"
	"time"
)

// MetaMainteTimeLayout ...
const MetaMainteTimeLayout = "2 Jan 2006 15:04:05 MST"

// MetaMainteTime ...
type MetaMainteTime time.Time

// UnmarshalJSON ...
func (ct *MetaMainteTime) UnmarshalJSON(b []byte) error {
	tm, err := time.Parse(`"`+MetaMainteTimeLayout+`"`, string(b))
	*ct = MetaMainteTime(tm)
	return err
}

// MarshalJSON ...
func (ct MetaMainteTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(ct).Format(MetaMainteTimeLayout))
}
