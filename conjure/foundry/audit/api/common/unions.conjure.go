// Deprecated: exists only for backwards compatibility.
// New usages should use the github.com/palantir/witchcraft-go-logging/conjure/audit-api/foundry/audit/api/common package.

package common

import (
	"github.com/palantir/witchcraft-go-logging/conjure/audit-api/foundry/audit/api/common"
)

type Resource = common.Resource

type ResourceVisitor = common.ResourceVisitor

type ResourceVisitorWithContext = common.ResourceVisitorWithContext

func NewResourceFromGotham(v GothamResource) Resource {
	return common.NewResourceFromGotham(v)
}

func NewResourceFromFoundry(v FoundryResource) Resource {
	return common.NewResourceFromFoundry(v)
}
