package discovery

//func IsApiServerReachable(ctx context.Context, clientSet kubernetes.Interface) (string, error) {
//	ver, err := clientSet.Discovery().RESTClient().
//	if err != nil {
//		return "", fmt.Errorf("cannot detect kube-apiserver version: %w", err)
//	}
//
//	return ver.String(), nil
//}
//
//func CheckApiServersAvailability(ctx context.Context) error {
//	var (
//		cfg = config.FromContext(ctx)
//		log = logging.FromContext(ctx)
//	)
//
//	ver1, err := getApiServerVersion(ctx, cfg.Connections.Cluster1.ClientSet)
//	if err != nil {
//		return err
//	}
//
//	ver2, err := getApiServerVersion(ctx, cfg.Connections.Cluster2.ClientSet)
//	if err != nil {
//		return err
//	}
//
//	if ver1 == ver2 {
//		log.Infof("discovered kube-apiserver versions: %s", ver1)
//	} else {
//		log.Warnf("discovered kube-apiserver version(s): %s vs %s", ver1, ver2)
//	}
//
//	return nil
//}