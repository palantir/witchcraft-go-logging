package categories

import (
	"fmt"
	werror "github.com/palantir/witchcraft-go-error"
	v2 "github.com/palantir/witchcraft-go-logging/conjure/foundry/audit/api/common/v2"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"reflect"
)

type Classification string

const (
	Classification_RESOURCE     = Classification("RESOURCE")
	Classification_TOKEN        = Classification("TOKEN")
	Classification_UID          = Classification("UID")
	Classification_DATA         = Classification("DATA")
	Classification_METADATA     = Classification("METADATA")
	Classification_USER_INPUT   = Classification("USER_INPUT")
	Classification_CONSTANT     = Classification("CONSTANT")
	Classification_PASS_THROUGH = Classification("PASS_THROUGH")
)

// Converts the provided data to a slice of v2.Resource. Based on the Java "promote" method.
func toResource(data any) ([]v2.Resource, error) {
	switch val := data.(type) {
	case v2.Resource:
		return []v2.Resource{val}, nil
	case v2.ApplicationResource:
		return []v2.Resource{v2.NewResourceFromApplication(val)}, nil
	case v2.DataResource:
		return []v2.Resource{v2.NewResourceFromData(val)}, nil
	case v2.MonitorResource:
		return []v2.Resource{v2.NewResourceFromMonitor(val)}, nil
	case v2.SystemResource:
		return []v2.Resource{v2.NewResourceFromSystem(val)}, nil
	case v2.LogicResource:
		return []v2.Resource{v2.NewResourceFromLogic(val)}, nil
	case v2.RequestResource:
		return []v2.Resource{v2.NewResourceFromRequest(val)}, nil
	case v2.OntologyLogicResource:
		return []v2.Resource{v2.NewResourceFromOntologyLogic(val)}, nil
	case v2.OntologyDataResource:
		return []v2.Resource{v2.NewResourceFromOntologyData(val)}, nil
	case v2.OntologyDataResourceList:
		return []v2.Resource{v2.NewResourceFromOntologyDataList(val)}, nil
	case v2.OntologyMetaDataResource:
		return []v2.Resource{v2.NewResourceFromOntologyMetaData(val)}, nil
	case v2.Identifier:
		return []v2.Resource{v2.NewResourceFromExternal(val)}, nil
	case v2.ImportDestination:
		resources := []v2.Resource{v2.NewResourceFromData(val.Id)}
		if val.Parent != nil {
			resources = append(resources, v2.NewResourceFromData(*val.Parent))
		}
		return resources, nil
	}
	return nil, werror.Error("Cannot convert non-resource type to resource", werror.SafeParam("type", fmt.Sprintf("%T", data)))
}

// CheckAndExtractResources takes an input and checks if it is a resource or a collection of resources and, if so,
// returns the v2.Resource objects for the input. Returns an error if the provided input is not a recognized resource
// type (or a slice of recognized resource types).
//
// Based on the Java "public static Set<Resource> checkAndExtractResources(Object data)" function.
func CheckAndExtractResources(input any) ([]v2.Resource, error) {
	if input == nil {
		return nil, nil
	}

	// if input is a slice, treat it as a collection of resources
	if valViaReflection := reflect.ValueOf(input); valViaReflection.Kind() == reflect.Slice {
		var collectedResources []v2.Resource
		for i := 0; i < valViaReflection.Len(); i++ {
			sliceVal := valViaReflection.Index(i).Interface()

			resources, err := toResource(sliceVal)
			if err != nil {
				return nil, werror.Wrap(err, "Data is not a Resource in a collection of Resources")
			}
			collectedResources = append(collectedResources, resources...)
		}
		return collectedResources, nil
	}

	// otherwise, treat it as a single resource
	return toResource(input)
}

// CheckAndExtractUIDs takes an input and checks if it is a UserId or a slice of UserIds and, if so, returns a slice of
// UserIds that contains all the UserIds. Returns an error if the provided input is not a UserId or a slice of UserIds.
//
// Based on the Java "public static Set<UserId> checkAndExtractUids(Object data)" function.
func CheckAndExtractUIDs(input any) ([]logging.UserId, error) {
	var uids []logging.UserId

	// if input is a slice, treat it as a collection of resources
	if valViaReflection := reflect.ValueOf(input); valViaReflection.Kind() == reflect.Slice {
		for i := 0; i < valViaReflection.Len(); i++ {
			sliceVal := valViaReflection.Index(i).Interface()
			userID, err := toUserID(sliceVal)
			if err != nil {
				return nil, werror.Wrap(err, "Data is not a UserId in a collection of UserIds")
			}
			uids = append(uids, userID)
		}
	} else {
		userID, err := toUserID(input)
		if err != nil {
			return nil, err
		}
		uids = append(uids, userID)
	}

	return uids, nil
}

func toUserID(in any) (logging.UserId, error) {
	userID, ok := in.(logging.UserId)
	if !ok {
		return "", werror.Error("Cannot convert non-UserId type to UserId", werror.SafeParam("type", fmt.Sprintf("%T", in)))
	}
	return userID, nil
}

func ExtractUserIDsFromTokens(tokens []v2.Token) []logging.UserId {
	var userIDs []logging.UserId
	for _, token := range tokens {
		if token.UserId == nil {
			continue
		}
		userIDs = append(userIDs, logging.UserId(*token.UserId))
	}
	return userIDs
}

//private static Set<UserId> extractUserIdsFromTokens(List<Token> tokens) {
//return tokens.stream()
//.map(Token::getUserId)
//.filter(Optional::isPresent)
//.map(Optional::get)
//.collect(Collectors.toSet());
//}
