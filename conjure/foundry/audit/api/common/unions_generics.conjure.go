// Deprecated: exists only for backwards compatibility.
// New usages should use the github.com/palantir/witchcraft-go-logging/conjure/audit-api/foundry/audit/api/common package.

//go:build go1.23

package common

import (
	"github.com/palantir/witchcraft-go-logging/conjure/audit-api/foundry/audit/api/common"
)

type ResourceWithT[T any] = common.ResourceWithT[T]

type ResourceVisitorWithT[T any] = common.ResourceVisitorWithT[T]
