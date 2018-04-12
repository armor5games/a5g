package a5gapi

func NewJSONMsgDefautlDebugs(errs ...error) []*APIErr {
	return newJSONMsgs(ErrSeverityDebug, errs...)
}

func NewJSONMsgDefautlInfos(errs ...error) []*APIErr {
	return newJSONMsgs(ErrSeverityInfo, errs...)
}

func NewJSONMsgDefautlWarns(errs ...error) []*APIErr {
	return newJSONMsgs(ErrSeverityWarn, errs...)
}

func NewJSONMsgDefautlErrors(errs ...error) []*APIErr {
	return newJSONMsgs(ErrSeverityError, errs...)
}

func NewJSONMsgDefautlFatals(errs ...error) []*APIErr {
	return newJSONMsgs(ErrSeverityFatal, errs...)
}

func NewJSONMsgDefautlPanics(errs ...error) []*APIErr {
	return newJSONMsgs(ErrSeverityPanic, errs...)
}

func newJSONMsgs(s ErrSeverity, errs ...error) []*APIErr {
	var e []*APIErr

	for _, err := range errs {
		e = append(e, &APIErr{
			Code:     s.ErrorDefaultCode(),
			Severity: uint64(s),
			Err:      err})
	}

	return e
}
