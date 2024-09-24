package evt2log

import (
	"encoding/json"
	"time"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
)

func ConvertWLogParams(params []evt2log.Param) []Param {
	out := make([]Param, len(params))
	for _, p := range params {
		entry := wlog.NewMapLogEntry()
		evt2log.ApplyParam(p, entry)
		entry.AllValues()
		jsonBytes, err := json.Marshal(entry)
		if err != nil {
			panic(err)
		}
		var logEntry logging.EventLogV2
		if err := json.Unmarshal(jsonBytes, &logEntry); err != nil {
			panic(err)
		}
		switch {
		case logEntry.Type != "":
			out = append(out, Type())
		case !time.Time(logEntry.Time).IsZero():
			out = append(out, Time(time.Time(logEntry.Time)))
		case logEntry.EventName != "":
			out = append(out, EventName(logEntry.EventName))
		case logEntry.Values != nil:
			out = append(out, SafeParams(logEntry.Values))
		case logEntry.Uid != nil:
			out = append(out, UID(string(*logEntry.Uid)))
		case logEntry.Sid != nil:
			out = append(out, SID(string(*logEntry.Sid)))
		case logEntry.TokenId != nil:
			out = append(out, TokenID(string(*logEntry.TokenId)))
		case logEntry.TraceId != nil:
			out = append(out, TraceID(string(*logEntry.TraceId)))
		case logEntry.UnsafeParams != nil:
			out = append(out, UnsafeParams(logEntry.UnsafeParams))
		case logEntry.Tags != nil:
			out = append(out, Tags(logEntry.Tags))
		}
	}
	return out
}
