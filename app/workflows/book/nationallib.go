package book

import (
	"encoding/xml"
	"strings"

	"isbnbook/app/repos"
)

type nationalLibReader struct {
	client repos.Client
}

type nationalLibData struct {
	XMLName            xml.Name `xml:"searchRetrieveResponse"`
	Text               string   `xml:",chardata"`
	Xmlns              string   `xml:"xmlns,attr"`
	Xsi                string   `xml:"xsi,attr"`
	Version            string   `xml:"version"`
	NumberOfRecords    string   `xml:"numberOfRecords"`
	NextRecordPosition string   `xml:"nextRecordPosition"`
	ExtraResponseData  struct {
		Text   string `xml:",chardata"`
		Facets struct {
			Text string `xml:",chardata"`
			Lst  []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
				Int  []struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
				} `xml:"int"`
			} `xml:"lst"`
		} `xml:"facets"`
	} `xml:"extraResponseData"`
	Records struct {
		Text   string `xml:",chardata"`
		Record struct {
			Text          string `xml:",chardata"`
			RecordSchema  string `xml:"recordSchema"`
			RecordPacking string `xml:"recordPacking"`
			RecordData    struct {
				Text string `xml:",chardata"`
				Dc   struct {
					Text        string `xml:",chardata"`
					DcndlSimple string `xml:"dcndl_simple,attr"`
					Dc          string `xml:"dc,attr"`
					Dcterms     string `xml:"dcterms,attr"`
					Dcndl       string `xml:"dcndl,attr"`
					Rdf         string `xml:"rdf,attr"`
					Rdfs        string `xml:"rdfs,attr"`
					Foaf        string `xml:"foaf,attr"`
					Owl         string `xml:"owl,attr"`
					Identifier  []struct {
						Text string `xml:",chardata"`
						Type string `xml:"type,attr"`
					} `xml:"identifier"`
					Title                    string   `xml:"title"`
					TitleTranscription       string   `xml:"titleTranscription"`
					Creator                  []string `xml:"creator"`
					CreatorTranscription     string   `xml:"creatorTranscription"`
					PartTitle                []string `xml:"partTitle"`
					PartCreator              []string `xml:"partCreator"`
					SeriesTitle              string   `xml:"seriesTitle"`
					SeriesTitleTranscription string   `xml:"seriesTitleTranscription"`
					Publisher                []string `xml:"publisher"`
					PublicationPlace         []struct {
						Text string `xml:",chardata"`
						Type string `xml:"type,attr"`
					} `xml:"publicationPlace"`
					Issued struct {
						Text string `xml:",chardata"`
						Type string `xml:"type,attr"`
					} `xml:"issued"`
					Price   string `xml:"price"`
					Extent  string `xml:"extent"`
					Subject []struct {
						Text string `xml:",chardata"`
						Type string `xml:"type,attr"`
					} `xml:"subject"`
					MaterialType string `xml:"materialType"`
					AccessRights string `xml:"accessRights"`
					SeeAlso      []struct {
						Text     string `xml:",chardata"`
						Resource string `xml:"resource,attr"`
					} `xml:"seeAlso"`
					Language struct {
						Text string `xml:",chardata"`
						Type string `xml:"type,attr"`
					} `xml:"language"`
				} `xml:"dc"`
			} `xml:"recordData"`
			RecordPosition string `xml:"recordPosition"`
		} `xml:"record"`
	} `xml:"records"`
}

var ndCodeList = map[string]string{
	"00": "総記",
	"01": "図書館、図書館学",
	"02": "図書、書誌学",
	"03": "図書館、百科事典",
	"04": "一般論文集、一般講演集",
	"05": "逐次刊行物",
	"06": "団体",
	"07": "ジャーナリズム、新聞",
	"08": "叢書、全集、選集",
	"09": "貴重書、郷土資料、その他の特別コレクション",

	"10": "哲学",
	"11": "哲学各論",
	"12": "東洋思想",
	"13": "西洋哲学",
	"14": "心理学",
	"15": "倫理学、道徳",
	"16": "宗教",
	"17": "神道",
	"18": "仏教",
	"19": "キリスト教",

	"20": "歴史",
	"21": "日本史",
	"22": "アジア史、東洋史",
	"23": " ヨーロッパ史、西洋史",
	"24": "アフリカ史",
	"25": "北アメリカ史",
	"26": "南アメリカ史",
	"27": "オセアニア史、両極地方史",
	"28": "伝記",
	"29": "地理、地誌、紀行",

	"30": "社会科学",
	"31": "政治",
	"32": "法律",
	"33": "経済",
	"34": "財政",
	"35": "統計",
	"36": "社会",
	"37": "教育",
	"38": "風俗習慣、民俗学、民族学",
	"39": "国防、軍事",

	"40": "自然科学",
	"41": "数学",
	"42": "物理学",
	"43": "化学",
	"44": "天文学、宇宙科学",
	"45": "地球科学、地学",
	"46": "生物科学、一般生物学",
	"47": "植物学",
	"48": "動物学",
	"49": "医学、薬学",

	"50": "技術、工学",
	"51": "建設工学、土木工事",
	"52": "建築学",
	"53": "機械工学、原子力工学",
	"54": "電気工学、電子工学",
	"55": "海洋工学、船舶工学、兵器",
	"56": "金属工学、鉱山工学",
	"57": "化学工業",
	"58": "製造工業",
	"59": "家政学、生活科学",

	"60": "産業",
	"61": "農業",
	"62": "園芸",
	"63": "蚕糸業",
	"64": "畜産業、獣医学",
	"65": "林業",
	"66": "水産業",
	"67": "商業",
	"68": "運輸、交通",
	"69": "通信事業",

	"70":  "芸術、美術",
	"71":  "彫刻",
	"72":  "絵画、書道",
	"726": "コミック",
	"73":  "版画",
	"74":  "写真、印刷",
	"75":  "工芸",
	"76":  "音楽、舞踊",
	"77":  "演劇、映画",
	"78":  "スポーツ、体育",
	"79":  "諸芸、娯楽",

	"80": "言語",
	"81": "日本語",
	"82": "中国語、その他の東洋の諸言語",
	"83": "英語",
	"84": "ドイツ語",
	"85": "フランス語",
	"86": "スペイン語",
	"87": "イタリア語",
	"88": "ロシア語",
	"89": "その他の諸言語",

	"90": "文学",
	"91": "日本文学",
	"92": "中国文学、その他の東洋文学",
	"93": "英米文学",
	"94": "ドイツ文学",
	"95": "フランス文学",
	"96": "スペイン文学",
	"97": "イタリア文学",
	"98": "ロシア、ソヴィエト文学",
	"99": "その他の諸文学",
}

func NewNationalLibReader() *nationalLibReader {
	cli, _ := repos.NewClient("http://iss.ndl.go.jp/api")
	return &nationalLibReader{
		client: cli,
	}
}

func (o *nationalLibReader) GetBookInfo(isbn13 string) (*BookInfo, error) {
	params := map[string]string{
		"operation":      "searchRetrieve",
		"version":        "1.2",
		"maximumRecords": "1",
		"query":          "isbn= \"" + isbn13 + "\"",
		"recordSchema":   "dcndl_simple",
	}
	byteArray, err := o.client.Get("/sru", params)
	if err != nil {
		return nil, err
	}
	var data *nationalLibData

	if err := xml.Unmarshal(byteArray, &data); err != nil {
		return nil, err
	}
	return o.getBookInfoFromData(data)
}

func (o *nationalLibReader) getBookInfoFromData(data *nationalLibData) (*BookInfo, error) {
	title := data.Records.Record.RecordData.Dc.Title

	publisher := ""
	if data.Records.Record.RecordData.Dc.Publisher != nil && len(data.Records.Record.RecordData.Dc.Publisher) > 0 {
		publisher = data.Records.Record.RecordData.Dc.Publisher[0]
	}

	pubdate := data.Records.Record.RecordData.Dc.Issued.Text

	var authors []string
	for _, rec := range data.Records.Record.RecordData.Dc.Creator {
		rec = strings.Replace(rec, " 著", "", 1)
		rec = strings.Trim(rec, " ")
		authors = append(authors, rec)
	}

	// national lib has no kind
	kind := ""
	genre := ""
	for _, code := range data.Records.Record.RecordData.Dc.Subject {
		if code.Type == "dcndl:NDC9" {
			for key := range ndCodeList {
				if strings.HasPrefix(code.Text, key) {
					genre = ndCodeList[key]
					break
				}
			}
		}
	}
	return NewBookInfo(title, authors, publisher, pubdate, kind, genre), nil
}
