package isbn

type pageRange struct {
	start int
	end int
}

func GetPageRange(totalPage int, pageOffset int) []pageRange {
	pRanges := make([]pageRange, 0)

	// get start pages
	if (totalPage > pageOffset) {
		pRanges = append(pRanges, pageRange{1, pageOffset})
	} else {
		pRanges = append(pRanges, pageRange{1, totalPage})
	}

	// get end pages
	var startPage = totalPage - pageOffset
	if (startPage > 0) {
		pRanges = append(pRanges, pageRange{startPage + 1, totalPage})
	}
	return pRanges
}