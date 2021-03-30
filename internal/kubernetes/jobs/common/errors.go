package common

import "errors"

var (
	//ErrorBackoffLimitDifferent  = errors.New("backoffLimit in jobs is different")
	ErrorBackoffLimitDifferent  = errors.New("job backoff limit is different")
	ErrorRestartPolicyDifferent = errors.New("job restartPolicy limit is different")

	ErrorScheduleDifferent = errors.New("schedule in cronJobs is different")
)
