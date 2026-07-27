//go:build amd64

package distance

import "golang.org/x/sys/cpu"

func init() {
	var features uint64
	if cpu.X86.HasAVX512BW {
		features |= CPUFeatureAVX512BW
	}
	if cpu.X86.HasAVX512F {
		features |= CPUFeatureAVX512F
	}
	if cpu.X86.HasAVX2 {
		features |= CPUFeatureAVX2
	}
	featureSet(features)
}
