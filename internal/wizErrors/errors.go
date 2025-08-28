package wizErrors

import "errors"

var MaxCapacityReached = errors.New("Max Capacity Reached")
var UserAlreadyInSession = errors.New("User Already In Session")
var ResourceNotFound = errors.New("Resource Not Found")
