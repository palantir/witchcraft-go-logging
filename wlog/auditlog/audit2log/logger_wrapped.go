package audit2log

type wrappedLogger struct {
	logger Logger
	params []Param
}

func (w wrappedLogger) Audit(name string, result AuditResultType, params ...Param) {
	w.logger.Audit(name, result, append(append([]Param{}, w.params...), params...)...)
}
