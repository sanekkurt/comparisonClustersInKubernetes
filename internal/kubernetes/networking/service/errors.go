package service

import (
	"errors"
)

var (
	ErrorPortsCountDifferent     = errors.New("the ports count are different")
	ErrorPortInServicesDifferent = errors.New("the port in the services is different")

	ErrorSelectorsCountDifferent     = errors.New("the selectors count are different")
	ErrorSelectorInServicesDifferent = errors.New("the selector in the services is different")

	ErrorTypeInServicesDifferent = errors.New("the type in the services is different")
)
