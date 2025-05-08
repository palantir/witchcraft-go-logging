package categories

import (
	"github.com/palantir/pkg/rid"
	v2 "github.com/palantir/witchcraft-go-logging/conjure/foundry/audit/api/common/v2"
)

func EntitiesFromResources(resource ...v2.Resource) ([]rid.ResourceIdentifier, error) {
	var ids []rid.ResourceIdentifier
	for _, r := range resource {
		rids, err := entitiesFromResource(r)
		if err != nil {
			return nil, err
		}
		ids = append(ids, rids...)
	}
	return ids, nil
}

func entitiesFromResource(resource v2.Resource) ([]rid.ResourceIdentifier, error) {
	var ids []rid.ResourceIdentifier

	ridsFromIDFn := func(id v2.Identifier) error {
		rids, err := ridsFromIdentifier(id)
		if err != nil {
			return err
		}
		ids = append(ids, rids...)
		return nil
	}
	ridsFn := func(rids []rid.ResourceIdentifier, err error) error {
		if err != nil {
			return err
		}
		ids = append(ids, rids...)
		return nil
	}

	if err := resource.AcceptFuncs(
		func(v v2.ApplicationResource) error {
			return ridsFromIDFn(v.Id)
		},
		func(v v2.DataResource) error {
			return ridsFromIDFn(v.Id)
		},
		func(v v2.MonitorResource) error {
			return ridsFromIDFn(v.Id)
		},
		resource.SystemNoopSuccess,
		func(v v2.LogicResource) error {
			return ridsFromIDFn(v.Id)
		},
		func(v v2.RequestResource) error {
			return ridsFromIDFn(v.Id)
		},
		func(v v2.OntologyDataResource) error {
			return ridsFn(ridsFromOntologyDataResource(v))
		},
		func(v v2.OntologyDataResourceList) error {
			return ridsFn(ridsFromOntologyDataResourceList(v))
		},
		func(v v2.OntologyLogicResource) error {
			return ridsFn(ridsFromOntologyLogicResource(v))
		},
		func(v v2.OntologyMetaDataResource) error {
			return ridsFromIDFn(v.Id)
		},
		func(v v2.Identifier) error {
			return ridsFromIDFn(v)
		},
		resource.ErrorOnUnknown,
	); err != nil {
		return nil, err
	}
	return ids, nil
}

func ridsFromIdentifier(identifier v2.Identifier) ([]rid.ResourceIdentifier, error) {
	var ids []rid.ResourceIdentifier
	if err := identifier.AcceptFuncs(
		func(v rid.ResourceIdentifier) error {
			ids = append(ids, v)
			return nil
		},
		func(v []rid.ResourceIdentifier) error {
			ids = append(ids, v...)
			return nil
		},
		identifier.OtherNoopSuccess,
		identifier.OthersNoopSuccess,
		identifier.ErrorOnUnknown,
	); err != nil {
		return nil, err
	}
	return ids, nil
}

func ridsFromOntologyDataResource(resource v2.OntologyDataResource) ([]rid.ResourceIdentifier, error) {
	var ids []rid.ResourceIdentifier
	for _, k := range resource.ObjectPrimaryKey {
		rids, err := ridsFromObjectPropertyIdentifier(k.PropertyIdentifier)
		if err != nil {
			return nil, err
		}
		ids = append(ids, rids...)
	}
	if resource.AdditionalObjectProperties != nil {
		for _, prop := range *resource.AdditionalObjectProperties {
			rids, err := ridsFromObjectPropertyIdentifier(prop.PropertyIdentifier)
			if err != nil {
				return nil, err
			}
			ids = append(ids, rids...)
		}
	}
	if resource.OntologyContext != nil {
		rids, err := ridsFromOntologyContext(*resource.OntologyContext)
		if err != nil {
			return nil, err
		}
		ids = append(ids, rids...)
	}
	return ids, nil
}

func ridsFromObjectPropertyIdentifier(identifier v2.ObjectPropertyIdentifier) ([]rid.ResourceIdentifier, error) {
	var ids []rid.ResourceIdentifier
	if err := identifier.AcceptFuncs(
		func(v rid.ResourceIdentifier) error {
			ids = append(ids, v)
			return nil
		},
		identifier.IdNoopSuccess,
		identifier.ErrorOnUnknown,
	); err != nil {
		return nil, err
	}
	return ids, nil
}

func ridsFromOntologyContext(ontologyContext v2.OntologyContext) ([]rid.ResourceIdentifier, error) {
	if ontologyContext.EntityType == nil {
		return nil, nil
	}
	var ids []rid.ResourceIdentifier
	if err := ontologyContext.EntityType.AcceptFuncs(
		func(v rid.ResourceIdentifier) error {
			ids = append(ids, v)
			return nil
		},
		ontologyContext.EntityType.IdNoopSuccess,
		ontologyContext.EntityType.ErrorOnUnknown,
	); err != nil {
		return nil, err
	}
	return ids, nil
}

func ridsFromOntologyDataResourceList(resource v2.OntologyDataResourceList) ([]rid.ResourceIdentifier, error) {
	var ids []rid.ResourceIdentifier

	rids, err := ridsFromOntologyContext(resource.SharedOntologyResourceContext.OntologyContext)
	if err != nil {
		return nil, err
	}
	ids = append(ids, rids...)

	for _, r := range resource.OntologyDataResources {
		rids, err := ridsFromOntologyDataResource(r)
		if err != nil {
			return nil, err
		}
		ids = append(ids, rids...)
	}
	return ids, nil
}

func ridsFromOntologyLogicResource(resource v2.OntologyLogicResource) ([]rid.ResourceIdentifier, error) {
	var ids []rid.ResourceIdentifier
	for _, ontologyContext := range resource.OntologyContext {
		rids, err := ridsFromOntologyContext(ontologyContext)
		if err != nil {
			return nil, err
		}
		ids = append(ids, rids...)
	}
	if resource.Id != nil {
		rids, err := ridsFromIdentifier(*resource.Id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, rids...)
	}
	return ids, nil
}
