package book

func getBookKindNameByCCode(code string) string {
	v, ok := codes[code]
	if ok {
		return v
	}
	return "不明"
}

var codes = map[string]string{
	"0": "単行本",
	"1": "文庫",
	"2": "新書",
	"3": "全集・双書",
	"4": "ムック・その他",
	"5": "事・辞典",
	"6": "図鑑",
	"7": "絵本",
	"8": "磁性媒体など",
	"9": "コミック",
}
