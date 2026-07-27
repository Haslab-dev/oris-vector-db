//go:build arm64

package distance

import "golang.org/x/sys/cpu"

func init() {
	var features uint64
	if cpu.ARM64.HasASIMD {
		features |= CPUFeatureNEON
	}
	if cpu.ARM64.HasSVE {
		features |= CPUFeatureSVE
	}
	featureSet(features)
}
