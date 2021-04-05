package logging

//func DiffLog(log *zap.SugaredLogger, e interface{}, v... interface{}) {
//	var (
//		diff string
//		vals = make([]string, len(v), len(v))
//	)
//
//	switch v := e.(type) {
//	case error:
//		diff = v.Error()
//	case string:
//		diff = v
//	}
//
//	for idx, vidx := range v {
//		switch val := vidx.(type) {
//		case string:
//			vals[idx] = val
//		default:
//			vals[idx] = fmt.Sprintln(val)
//		}
//	}
//
//	switch len(vals) {
//	case 0:
//		log.Warnf("%s", diff)
//	case 1:
//		log.Warnf("%s: %s", diff, vals[0])
//	case 2:
//		log.Warnf("%s: %s vs %s", diff, vals[0], vals[1])
//	default:
//		log.Warnf("%s, %s: %s vs %s", diff, strings.Join(vals[0:len(vals)-2], " "), vals[len(vals)-2], vals[len(vals)-1])
//	}
//}
