package healthcheck

import (
	"context"

	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func compareHTTPGetProbes(ctx context.Context, probe1, probe2 v1.Probe) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if probe1.HTTPGet != nil && probe2.HTTPGet != nil {
		if probe1.HTTPGet.Host != probe2.HTTPGet.Host {
			//logging.DiffLog(log, ErrorContainerHealthCheckDifferentHTTPGetHost, probe1.HTTPGet.Host, probe2.HTTPGet.Host)
			diffsBatch.Add(ctx, false, "%w: '%s' vs '%s'", ErrorContainerHealthCheckDifferentHTTPGetHost, probe1.HTTPGet.Host, probe2.HTTPGet.Host)
		}

		if probe1.HTTPGet.HTTPHeaders != nil && probe2.HTTPGet.HTTPHeaders != nil {
			if len(probe1.HTTPGet.HTTPHeaders) != len(probe2.HTTPGet.HTTPHeaders) {
				//logging.DiffLog(log, ErrorContainerHealthCheckDifferentHTTPGetHTTPHeaders, len(probe1.HTTPGet.HTTPHeaders), len(probe2.HTTPGet.HTTPHeaders))
				diffsBatch.Add(ctx, false, "%w: '%d' vs '%d'", ErrorContainerHealthCheckDifferentHTTPGetHTTPHeaders, len(probe1.HTTPGet.HTTPHeaders), len(probe2.HTTPGet.HTTPHeaders))
			} else {

				for index, value := range probe1.HTTPGet.HTTPHeaders {
					if value.Name != probe2.HTTPGet.HTTPHeaders[index].Name {
						//logging.DiffLog(log, ErrorContainerHealthCheckDifferentHTTPGetHeaderName, index+1, probe1.HTTPGet.HTTPHeaders[index].Name, probe2.HTTPGet.HTTPHeaders[index].Name)
						diffsBatch.Add(ctx, false, "%w, %d: '%s' vs '%s'", ErrorContainerHealthCheckDifferentHTTPGetHeaderName, index+1, probe1.HTTPGet.HTTPHeaders[index].Name, probe2.HTTPGet.HTTPHeaders[index].Name)
					}

					if value.Value != probe2.HTTPGet.HTTPHeaders[index].Value {
						//logging.DiffLog(log, ErrorContainerHealthCheckDifferentHTTPGetHeaderValue, value.Name, value.Value, probe2.HTTPGet.HTTPHeaders[index].Value)
						diffsBatch.Add(ctx, false, "%w, %s: '%s' vs '%s'", ErrorContainerHealthCheckDifferentHTTPGetHeaderValue, value.Name, value.Value, probe2.HTTPGet.HTTPHeaders[index].Value)
					}
				}

			}

		} else if probe1.HTTPGet.HTTPHeaders != nil || probe2.HTTPGet.HTTPHeaders != nil {
			//logging.DiffLog(log, ErrorContainerHealthCheckHTTPGetMissingHeader, probe1.HTTPGet.HTTPHeaders, probe2.HTTPGet.HTTPHeaders)
			diffsBatch.Add(ctx, false, "%w", ErrorContainerHealthCheckHTTPGetMissingHeader)
		}

		if probe1.HTTPGet.Path != probe2.HTTPGet.Path {
			//logging.DiffLog(log, ErrorContainerHealthCheckDifferentHTTPGetPath, probe1.HTTPGet.Path, probe2.HTTPGet.Path)
			diffsBatch.Add(ctx, false, "%w: '%s' vs '%s'", ErrorContainerHealthCheckDifferentHTTPGetPath, probe1.HTTPGet.Path, probe2.HTTPGet.Path)
		}

		if probe1.HTTPGet.Port.IntVal != probe2.HTTPGet.Port.IntVal {
			diffsBatch.Add(ctx, false, "%w: '%d' vs '%d'", ErrorContainerHealthCheckDifferentHTTPGetPortIntVal, probe1.HTTPGet.Port.IntVal, probe2.HTTPGet.Port.IntVal)
		}

		if probe1.HTTPGet.Port.StrVal != probe2.HTTPGet.Port.StrVal {
			diffsBatch.Add(ctx, false, "%w: '%s' vs '%s'", ErrorContainerHealthCheckDifferentHTTPGetPortStrVal, probe1.HTTPGet.Port.StrVal, probe2.HTTPGet.Port.StrVal)
		}

		if probe1.HTTPGet.Port.Type != probe2.HTTPGet.Port.Type {
			diffsBatch.Add(ctx, false, "%w: '%d' vs '%d'", ErrorContainerHealthCheckDifferentHTTPGetPortType, probe1.HTTPGet.Port.Type, probe2.HTTPGet.Port.Type)
		}

		//if probe1.HTTPGet.Port.IntVal != probe2.HTTPGet.Port.IntVal ||
		//	probe1.HTTPGet.Port.StrVal != probe2.HTTPGet.Port.StrVal ||
		//	probe1.HTTPGet.Port.Type != probe2.HTTPGet.Port.Type {
		//	//logging.DiffLog(log, ErrorContainerHealthCheckDifferentHTTPGetPort, probe1.HTTPGet.Port, probe2.HTTPGet.Port)
		//	diffsBatch.Add(ctx, false, "%w", ErrorContainerHealthCheckDifferentHTTPGetPort)
		//}

		if probe1.HTTPGet.Scheme != probe2.HTTPGet.Scheme {
			//logging.DiffLog(log, ErrorContainerHealthCheckDifferentHTTPGetScheme, probe1.HTTPGet.Scheme, probe2.HTTPGet.Scheme)
			diffsBatch.Add(ctx, false, "%w: '%s' vs '%s'", ErrorContainerHealthCheckDifferentHTTPGetScheme, probe1.HTTPGet.Scheme, probe2.HTTPGet.Scheme)
		}

	} else if probe1.HTTPGet != nil || probe2.HTTPGet != nil {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentHTTPGet, probe1.HTTPGet, probe2.HTTPGet)
		diffsBatch.Add(ctx, false, "%w", ErrorContainerHealthCheckDifferentHTTPGet)
	}

}
