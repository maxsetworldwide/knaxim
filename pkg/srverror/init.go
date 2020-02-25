package srverror

import "flag"

func init() {
	flag.BoolVar(&DEBUG, "dev", DEBUG, "development mode, server errors send full report")
}
