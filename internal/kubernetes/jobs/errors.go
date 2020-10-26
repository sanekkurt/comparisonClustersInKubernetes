package jobs

import "errors"

var (
	ErrorBackoffLimitDifferent  = errors.New("backoffLimit in jobs is different")
	ErrorRestartPolicyDifferent = errors.New("restartPolicy in jobs is different")

	ErrorScheduleDifferent = errors.New("schedule in cronJobs is different")
)
