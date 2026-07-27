//go:build arm64

package distance

func init() {
	if CPUInfo.HasNEON {
		RegisterBackend(Backend{
			Name:    "neon",
			Kernel:  &neonKernel{},
			Checker: func() bool { return CPUInfo.HasNEON },
		})
	}
}
