package goarmorapi

func NewJSONMsgDefautlDebugs(errs ...error) []*JSONMsg {
	return newJSONMsgs(MsgSeverityDebug, errs...)
}

func NewJSONMsgDefautlInfos(errs ...error) []*JSONMsg {
	return newJSONMsgs(MsgSeverityInfo, errs...)
}

func NewJSONMsgDefautlWarns(errs ...error) []*JSONMsg {
	return newJSONMsgs(MsgSeverityWarn, errs...)
}

func NewJSONMsgDefautlErrors(errs ...error) []*JSONMsg {
	return newJSONMsgs(MsgSeverityError, errs...)
}

func NewJSONMsgDefautlFatals(errs ...error) []*JSONMsg {
	return newJSONMsgs(MsgSeverityFatal, errs...)
}

func NewJSONMsgDefautlPanics(errs ...error) []*JSONMsg {
	return newJSONMsgs(MsgSeverityPanic, errs...)
}

func newJSONMsgs(s JSONMsgSeverity, errs ...error) []*JSONMsg {
	var e []*JSONMsg

	for _, err := range errs {
		e = append(e, &JSONMsg{
			Code:     s.ErrorDefaultCode(),
			Severity: uint64(s),
			Err:      err})
	}

	return e
}
