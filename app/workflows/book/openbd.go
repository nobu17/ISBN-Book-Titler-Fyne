package book

import (
	"encoding/json"
)

type openBDReader struct {
	client *client
}

type openBDData struct {
	Onix struct {
		RecordReference   string `json:"RecordReference"`
		NotificationType  string `json:"NotificationType"`
		ProductIdentifier struct {
			ProductIDType string `json:"ProductIDType"`
			IDValue       string `json:"IDValue"`
		} `json:"ProductIdentifier"`
		DescriptiveDetail struct {
			ProductComposition string `json:"ProductComposition"`
			ProductForm        string `json:"ProductForm"`
			ProductFormDetail  string `json:"ProductFormDetail"`
			Collection         struct {
				CollectionType     string `json:"CollectionType"`
				CollectionSequence struct {
					CollectionSequenceType     string `json:"CollectionSequenceType"`
					CollectionSequenceTypeName string `json:"CollectionSequenceTypeName"`
					CollectionSequenceNumber   string `json:"CollectionSequenceNumber"`
				} `json:"CollectionSequence"`
				CollectionSequenceArray []struct {
					CollectionSequenceType     string `json:"CollectionSequenceType"`
					CollectionSequenceTypeName string `json:"CollectionSequenceTypeName"`
					CollectionSequenceNumber   string `json:"CollectionSequenceNumber"`
				} `json:"CollectionSequenceArray"`
			} `json:"Collection"`
			TitleDetail struct {
				TitleType    string `json:"TitleType"`
				TitleElement struct {
					TitleElementLevel string `json:"TitleElementLevel"`
					TitleText         struct {
						Collationkey string `json:"collationkey"`
						Content      string `json:"content"`
					} `json:"TitleText"`
					Subtitle struct {
						Collationkey string `json:"collationkey"`
						Content      string `json:"content"`
					} `json:"Subtitle"`
				} `json:"TitleElement"`
			} `json:"TitleDetail"`
			Contributor []struct {
				SequenceNumber  string   `json:"SequenceNumber"`
				ContributorRole []string `json:"ContributorRole"`
				PersonName      struct {
					Collationkey string `json:"collationkey"`
					Content      string `json:"content"`
				} `json:"PersonName"`
			} `json:"Contributor"`
			Language []struct {
				LanguageRole string `json:"LanguageRole"`
				LanguageCode string `json:"LanguageCode"`
			} `json:"Language"`
			Extent []struct {
				ExtentType  string `json:"ExtentType"`
				ExtentValue string `json:"ExtentValue"`
				ExtentUnit  string `json:"ExtentUnit"`
			} `json:"Extent"`
			Subject []struct {
				MainSubject             string `json:"MainSubject,omitempty"`
				SubjectSchemeIdentifier string `json:"SubjectSchemeIdentifier"`
				SubjectCode             string `json:"SubjectCode"`
			} `json:"Subject"`
			Audience []struct {
				AudienceCodeType  string `json:"AudienceCodeType"`
				AudienceCodeValue string `json:"AudienceCodeValue"`
			} `json:"Audience"`
		} `json:"DescriptiveDetail"`
		CollateralDetail struct {
			TextContent []struct {
				TextType        string `json:"TextType"`
				ContentAudience string `json:"ContentAudience"`
				Text            string `json:"Text"`
			} `json:"TextContent"`
			SupportingResource []struct {
				ResourceContentType string `json:"ResourceContentType"`
				ContentAudience     string `json:"ContentAudience"`
				ResourceMode        string `json:"ResourceMode"`
				ResourceVersion     []struct {
					ResourceForm           string `json:"ResourceForm"`
					ResourceVersionFeature []struct {
						ResourceVersionFeatureType string `json:"ResourceVersionFeatureType"`
						FeatureValue               string `json:"FeatureValue"`
					} `json:"ResourceVersionFeature"`
					ResourceLink string `json:"ResourceLink"`
				} `json:"ResourceVersion"`
			} `json:"SupportingResource"`
		} `json:"CollateralDetail"`
		PublishingDetail struct {
			Imprint struct {
				ImprintIdentifier []struct {
					ImprintIDType string `json:"ImprintIDType"`
					IDValue       string `json:"IDValue"`
				} `json:"ImprintIdentifier"`
				ImprintName string `json:"ImprintName"`
			} `json:"Imprint"`
			Publisher struct {
				PublishingRole      string `json:"PublishingRole"`
				PublisherIdentifier []struct {
					PublisherIDType string `json:"PublisherIDType"`
					IDValue         string `json:"IDValue"`
				} `json:"PublisherIdentifier"`
				PublisherName string `json:"PublisherName"`
			} `json:"Publisher"`
			PublishingDate []struct {
				PublishingDateRole string `json:"PublishingDateRole"`
				Date               string `json:"Date"`
			} `json:"PublishingDate"`
		} `json:"PublishingDetail"`
		ProductSupply struct {
			MarketPublishingDetail struct {
				MarketPublishingStatus string `json:"MarketPublishingStatus"`
			} `json:"MarketPublishingDetail"`
			SupplyDetail struct {
				ProductAvailability string `json:"ProductAvailability"`
				Price               []struct {
					PriceType    string `json:"PriceType"`
					PriceAmount  string `json:"PriceAmount"`
					CurrencyCode string `json:"CurrencyCode"`
				} `json:"Price"`
			} `json:"SupplyDetail"`
		} `json:"ProductSupply"`
	} `json:"onix"`
	Hanmoto struct {
		Datemodified string `json:"datemodified"`
		Datecreated  string `json:"datecreated"`
		Datekoukai   string `json:"datekoukai"`
	} `json:"hanmoto"`
	Summary struct {
		Isbn      string `json:"isbn"`
		Title     string `json:"title"`
		Volume    string `json:"volume"`
		Series    string `json:"series"`
		Publisher string `json:"publisher"`
		Pubdate   string `json:"pubdate"`
		Cover     string `json:"cover"`
		Author    string `json:"author"`
	} `json:"summary"`
}

var openBDCodes = map[string]string{
	"01": "文芸",
	"02": "新書",
	"03": "社会一般",
	"04": "資格.試験",
	"05": "ビジネス",
	"06": "スポーツ.健康",
	"07": "趣味.実用",
	"09": "ゲーム",
	"10": "芸能.タレント",

	"11": "テレビ.映画化",
	"12": "芸術",
	"13": "哲学.宗教",
	"14": "歴史.地理",
	"15": "社会科学",
	"16": "教育",
	"17": "自然科学",
	"18": "医学",
	"19": "工業.工学",
	"20": "コンピュータ",

	"21": "語学.辞事典",
	"22": "学参",
	"23": "児童図書",
	"24": "ヤングアダルト",
	"29": "新刊セット",
	"30": "全集",

	"31": "文庫",
	"36": "コミック文庫",

	"41": "コミックス(欠番扱)",
	"42": "コミックス(雑誌扱)",
	"43": "コミックス(書籍)",
	"44": "コミックス(廉価版)",

	"51": "ムック",
}

func NewOpenBDReader() *openBDReader {
	cli, _ := NewClient("https://api.openbd.jp/v1")
	return &openBDReader{
		client: cli,
	}
}

func (o *openBDReader) GetBookInfo(isbn13 string) (*BookInfo, error) {
	params := map[string]string{
		"isbn": isbn13,
	}
	byteArray, err := o.client.Get("/get", params)
	if err != nil {
		return nil, err
	}
	var data []*openBDData

	if err := json.Unmarshal(byteArray, &data); err != nil {
		return nil, err
	}
	return o.getBookInfoFromBDData(data[0])
}

func (o *openBDReader) getBookInfoFromBDData(data *openBDData) (*BookInfo, error) {
	// *(data).Onix
	title := data.Summary.Title
	publisher := data.Summary.Publisher
	pubdate := data.Summary.Pubdate // yyyyMMdd

	var authors []string
	for _, cont := range data.Onix.DescriptiveDetail.Contributor {
		authors = append(authors, cont.PersonName.Content)
	}

	kind := ""
	genre := ""
	codes := data.Onix.DescriptiveDetail.Subject
	for _, code := range codes {
		if code.SubjectSchemeIdentifier == "78" {
			ccode := code.SubjectCode
			kind = getBookKindNameByCCode(ccode[1:2])
			// can get genre from ccode but it is pending currently.
		} else if code.SubjectSchemeIdentifier == "79" {
			genre = getGenreNameByGenreCode(code.SubjectCode)
		}
	}
	return NewBookInfo(title, authors, publisher, pubdate, kind, genre), nil
}

func getGenreNameByGenreCode(genreCode string) string {
	v, ok := openBDCodes[genreCode]
	if ok {
		return v
	}
	return "不明"
}
