package versions

const (
	currentMain  = 0
	currentMinor   = 1
	currentPatch = 0
)

func getCurrent() *Version {
	return &Version{
		currentMain,
		currentMinor,
		currentPatch,
	}
}
