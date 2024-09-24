package req2log

import (
	"io"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/extractor"
)

type LoggerBuilder interface {
	LoggerCreator(creator wlog.LoggerCreator)
	IdsExtractor(idsExtractor extractor.IDsFromRequest)

	SafePathParams(safePathParams []string)
	ForbiddenPathParams(forbiddenPathParams []string)

	SafeQueryParams(safeQueryParams []string)
	ForbiddenQueryParams(forbiddenQueryParams []string)

	SafeHeaderParams(safeHeaderParams []string)
	ForbiddenHeaderParams(forbiddenHeaderParams []string)
}

type defaultLoggerBuilder struct {
	idsExtractor extractor.IDsFromRequest

	safePathParams      []string
	forbiddenPathParams []string

	safeQueryParams      []string
	forbiddenQueryParams []string

	safeHeaderParams      []string
	forbiddenHeaderParams []string
}

func (b *defaultLoggerBuilder) LoggerCreator(creator wlog.LoggerCreator) {
	b.loggerCreator = creator
}

func (b *defaultLoggerBuilder) IdsExtractor(idsExtractor extractor.IDsFromRequest) {
	b.idsExtractor = idsExtractor
}

func (b *defaultLoggerBuilder) SafePathParams(safePathParams []string) {
	b.safePathParams = append(b.safePathParams, safePathParams...)
}

func (b *defaultLoggerBuilder) ForbiddenPathParams(forbiddenPathParams []string) {
	b.forbiddenPathParams = append(b.forbiddenPathParams, forbiddenPathParams...)
}

func (b *defaultLoggerBuilder) SafeQueryParams(safeQueryParams []string) {
	b.safeQueryParams = append(b.safeQueryParams, safeQueryParams...)
}

func (b *defaultLoggerBuilder) ForbiddenQueryParams(forbiddenQueryParams []string) {
	b.forbiddenQueryParams = append(b.forbiddenQueryParams, forbiddenQueryParams...)
}

func (b *defaultLoggerBuilder) SafeHeaderParams(safeHeaderParams []string) {
	b.safeHeaderParams = append(b.safeHeaderParams, safeHeaderParams...)
}

func (b *defaultLoggerBuilder) ForbiddenHeaderParams(forbiddenHeaderParams []string) {
	b.forbiddenHeaderParams = append(b.forbiddenHeaderParams, forbiddenHeaderParams...)
}

func (b *defaultLoggerBuilder) build(w io.Writer) *defaultLogger {
	defaultParams := DefaultRequestParamPerms()
	return &defaultLogger{
		logger:           b.loggerCreator(w),
		idsExtractor:     b.idsExtractor,
		pathParamPerms:   CombinedParamPerms(defaultParams.PathParamPerms(), NewParamPerms(b.safePathParams, b.forbiddenPathParams)),
		queryParamPerms:  CombinedParamPerms(defaultParams.QueryParamPerms(), NewParamPerms(b.safeQueryParams, b.forbiddenQueryParams)),
		headerParamPerms: CombinedParamPerms(defaultParams.HeaderParamPerms(), NewParamPerms(b.safeHeaderParams, b.forbiddenHeaderParams)),
	}
}
