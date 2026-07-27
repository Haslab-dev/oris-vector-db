// Package distance provides similarity/distance computation kernels.
//
// Oris uses these kernels for all similarity computations. The implementation
// auto-selects the fastest available backend based on CPU capabilities.
package distance

// CPUFeature flags for distance kernel capability detection.
const (
	CPUFeatureNone    uint64 = 0
	CPUFeatureNEON    uint64 = 1 << iota // ARM NEON
	CPUFeatureAVX2                        // x86 AVX2
	CPUFeatureAVX512F                     // x86 AVX-512 Foundation
	CPUFeatureAVX512BW                    // x86 AVX-512 Byte/Word
	CPUFeatureSVE                         // ARM Scalable Vector Extensions
)

// CPUInfo exposes detected CPU capabilities. Populated at init time.
var CPUInfo struct {
	Features uint64
	HasNEON  bool
	HasAVX2  bool
	HasAVX512F  bool
	HasAVX512BW bool
	HasSVE   bool
}

func featureSet(features uint64) {
	CPUInfo.Features = features
	CPUInfo.HasNEON = features&CPUFeatureNEON != 0
	CPUInfo.HasAVX2 = features&CPUFeatureAVX2 != 0
	CPUInfo.HasAVX512F = features&CPUFeatureAVX512F != 0
	CPUInfo.HasAVX512BW = features&CPUFeatureAVX512BW != 0
	CPUInfo.HasSVE = features&CPUFeatureSVE != 0
}
