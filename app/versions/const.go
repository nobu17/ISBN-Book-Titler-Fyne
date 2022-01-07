package versions

const (
	currentMain  = 0
	currentSub   = 1
	currentPatch = 0
)

func getCurrent() *Version {
	return &Version{
		currentMain,
		currentSub,
		currentPatch,
	}
}
