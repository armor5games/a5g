package goarmorapi

func NewErrorJSONDefautlDebugs(errs ...error) []*ErrorJSON {
	return newErrorJSONs(ErrSeverityDebug, errs...)
}

func NewErrorJSONDefautlInfos(errs ...error) []*ErrorJSON {
	return newErrorJSONs(ErrSeverityInfo, errs...)
}

func NewErrorJSONDefautlWarns(errs ...error) []*ErrorJSON {
	return newErrorJSONs(ErrSeverityWarn, errs...)
}

func NewErrorJSONDefautlErrors(errs ...error) []*ErrorJSON {
	return newErrorJSONs(ErrSeverityError, errs...)
}

func NewErrorJSONDefautlFatals(errs ...error) []*ErrorJSON {
	return newErrorJSONs(ErrSeverityFatal, errs...)
}

func NewErrorJSONDefautlPanics(errs ...error) []*ErrorJSON {
	return newErrorJSONs(ErrSeverityPanic, errs...)
}

func newErrorJSONs(s ErrorJSONSeverity, errs ...error) []*ErrorJSON {
	var e []*ErrorJSON

	for _, err := range errs {
		e = append(e, &ErrorJSON{
			Code:     s.ErrorDefaultCode(),
			Severity: uint64(s),
			Err:      err})
	}

	return e
}
