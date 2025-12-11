// Deprecated: exists only for backwards compatibility.
// New usages should use the github.com/palantir/witchcraft-go-logging/conjure/audit-api/foundry/audit/api/common package.

package common

import (
	"github.com/palantir/witchcraft-go-logging/conjure/audit-api/foundry/audit/api/common"
)

type FileType = common.FileType

type FileType_Value = common.FileType_Value

const (
	FileType_DATASET = common.FileType_DATASET
	FileType_BINARY  = common.FileType_BINARY
	FileType_UNKNOWN = common.FileType_UNKNOWN
)

// FileType_Values returns all known variants of FileType.
func FileType_Values() []FileType_Value {
	return common.FileType_Values()
}

func New_FileType(value FileType_Value) FileType {
	return common.New_FileType(value)
}
