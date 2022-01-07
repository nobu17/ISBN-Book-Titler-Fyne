package book

import (
	"encoding/json"
	"fmt"
	"strings"

	"isbnbook/app/log"
	"isbnbook/app/repos"
)

type rakutenReader struct {
	apikey string
	client repos.Client
	logger log.AppLogger
}

type rakutenData struct {
	GenreInformation []interface{} `json:"GenreInformation"`
	Items            []struct {
		Item struct {
			AffiliateURL   string `json:"affiliateUrl"`
			Author         string `json:"author"`
			AuthorKana     string `json:"authorKana"`
			Availability   string `json:"availability"`
			BooksGenreID   string `json:"booksGenreId"`
			ChirayomiURL   string `json:"chirayomiUrl"`
			Contents       string `json:"contents"`
			DiscountPrice  int    `json:"discountPrice"`
			DiscountRate   int    `json:"discountRate"`
			Isbn           string `json:"isbn"`
			ItemCaption    string `json:"itemCaption"`
			ItemPrice      int    `json:"itemPrice"`
			ItemURL        string `json:"itemUrl"`
			LargeImageURL  string `json:"largeImageUrl"`
			LimitedFlag    int    `json:"limitedFlag"`
			ListPrice      int    `json:"listPrice"`
			MediumImageURL string `json:"mediumImageUrl"`
			PostageFlag    int    `json:"postageFlag"`
			PublisherName  string `json:"publisherName"`
			ReviewAverage  string `json:"reviewAverage"`
			ReviewCount    int    `json:"reviewCount"`
			SalesDate      string `json:"salesDate"`
			SeriesName     string `json:"seriesName"`
			SeriesNameKana string `json:"seriesNameKana"`
			Size           string `json:"size"`
			SmallImageURL  string `json:"smallImageUrl"`
			SubTitle       string `json:"subTitle"`
			SubTitleKana   string `json:"subTitleKana"`
			Title          string `json:"title"`
			TitleKana      string `json:"titleKana"`
		} `json:"Item"`
	} `json:"Items"`
	Carrier   int `json:"carrier"`
	Count     int `json:"count"`
	First     int `json:"first"`
	Hits      int `json:"hits"`
	Last      int `json:"last"`
	Page      int `json:"page"`
	PageCount int `json:"pageCount"`
}

type rakutenGenreData struct {
	Children []interface{} `json:"children"`
	Current  struct {
		BooksGenreID   string `json:"booksGenreId"`
		BooksGenreName string `json:"booksGenreName"`
		GenreLevel     int    `json:"genreLevel"`
	} `json:"current"`
	Parents []struct {
		Parent struct {
			BooksGenreID   string `json:"booksGenreId"`
			BooksGenreName string `json:"booksGenreName"`
			GenreLevel     int    `json:"genreLevel"`
		} `json:"parent"`
	} `json:"parents"`
}

func NewRakutenReader(apikey string) *rakutenReader {
	cli, _ := repos.NewClient("https://app.rakuten.co.jp/services/api")
	log := log.GetLogger()
	return &rakutenReader{
		apikey: apikey,
		client: cli,
		logger: log,
	}
}

func (r *rakutenReader) GetBookInfo(isbn13 string) (*BookInfo, error) {
	params := map[string]string{
		"format":        "json",
		"isbn":          isbn13,
		"applicationId": r.apikey,
	}
	byteArray, err := r.client.Get("/BooksBook/Search/20170404", params)
	if err != nil {
		return nil, err
	}
	var data *rakutenData

	if err := json.Unmarshal(byteArray, &data); err != nil {
		return nil, err
	}
	return r.getBookInfoFromRakutenData(data)
}

func (r *rakutenReader) getBookInfoFromRakutenData(data *rakutenData) (*BookInfo, error) {
	if data.Count == 0 || len(data.Items) == 0 {
		r.logger.Error("not data from api", nil)
		return nil, fmt.Errorf("no data from api")
	}

	item := data.Items[0].Item

	title := item.Title
	publisher := item.PublisherName
	pubdate := item.SalesDate // yyyy年MM月dd
	pubdate = strings.Replace(pubdate, "年", "-", -1)
	pubdate = strings.Replace(pubdate, "月", "-", -1)
	pubdate = strings.Replace(pubdate, "日", "-", -1)
	pubdate = strings.Replace(pubdate, "頃", "", -1)
	pubdate = strings.TrimSuffix(pubdate, "-")

	authors := strings.Split(item.Author, "/")

	kind := item.Size
	genre := r.getGenreString(item.BooksGenreID)

	return NewBookInfo(title, authors, publisher, pubdate, kind, genre), nil
}

func (r *rakutenReader) getGenreString(genreCodeStr string) string {
	genreCodes := strings.Split(genreCodeStr, "/")
	if len(genreCodes) < 1 {
		return ""
	}
	data, err := r.getGenre(genreCodes[0])
	if err != nil {
		return ""
	}
	if len(data.Parents) > 0 {
		return data.Parents[0].Parent.BooksGenreName
	}
	return data.Current.BooksGenreName
}

func (r *rakutenReader) getGenre(genreCode string) (*rakutenGenreData, error) {
	params := map[string]string{
		"format":        "json",
		"booksGenreId":  genreCode,
		"applicationId": r.apikey,
	}
	byteArray, err := r.client.Get("/BooksGenre/Search/20121128", params)
	if err != nil {
		return nil, err
	}
	var data *rakutenGenreData

	if err := json.Unmarshal(byteArray, &data); err != nil {
		return nil, err
	}
	return data, nil
}
