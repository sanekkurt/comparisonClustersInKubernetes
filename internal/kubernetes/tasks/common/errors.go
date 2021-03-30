package common

import "errors"

var (
	ErrorBackoffLimitDifferent  = errors.New("backoffLimit in tasks is different")
	ErrorRestartPolicyDifferent = errors.New("restartPolicy in tasks is different")

	ErrorScheduleDifferent = errors.New("schedule in cronJobs is different")
)
