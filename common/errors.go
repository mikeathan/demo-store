package common

import "errors"

var ErrorAuthorizationHeaderMissing = errors.New("Authorization header missing")
var ErrorInvalidAuthorizationHeader = errors.New("Invalid authorization header")
var ErrorAuthorizationFailed = errors.New("Authorization failed")
var ErrorKeyNotSet = errors.New("Store key not specified")
var ErrorStoreValueNotSet = errors.New("Store value not specified")
var ErrorKeyNotFound = errors.New("Key not found")
var ErrorUnauthorisedOwner = errors.New("Owner not authorised to update value")
var ErrorCreatingJwtToken error = errors.New("Error creating the token")
var ErrorParsigJwtToken error = errors.New("Error parsing the token")
var ErrorValidatingJwtToken error = errors.New("Error validating token")
