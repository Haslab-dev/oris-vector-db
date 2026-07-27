//go:build amd64

package distance

func init() {
	// Register backends in priority order (most capable first).
	if CPUInfo.HasAVX512BW || CPUInfo.HasAVX512F {
		RegisterBackend(Backend{
			Name:    "avx512",
			Kernel:  &avx512Kernel{},
			Checker: func() bool { return CPUInfo.HasAVX512BW || CPUInfo.HasAVX512F },
		})
	}
	if CPUInfo.HasAVX2 {
		RegisterBackend(Backend{
			Name:    "avx2",
			Kernel:  &avx2Kernel{},
			Checker: func() bool { return CPUInfo.HasAVX2 },
		})
	}
}
