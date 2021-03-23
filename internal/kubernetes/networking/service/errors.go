package service

import (
	"errors"
)

var (
	ErrorPortsCountDifferent     = errors.New("the ports count are different in services specs")
	ErrorPortInServicesDifferent = errors.New("the port is different in services specs")

	ErrorSelectorsCountDifferent     = errors.New("the selectors count in the services specs is different")
	ErrorSelectorInServicesDifferent = errors.New("the selector in the services specs is different")

	ErrorTypeInServicesDifferent = errors.New("the type in the services specs is different")
)
