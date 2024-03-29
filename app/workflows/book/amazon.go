package book

import (
	"fmt"
	"net/http"
	"strings"

	"isbnbook/app/utils"

	paapi5 "github.com/spiegel-im-spiegel/pa-api"
	"github.com/spiegel-im-spiegel/pa-api/entity"
	"github.com/spiegel-im-spiegel/pa-api/query"
)

type amazonPAReader struct {
	client paapi5.Client
}

func NewAmazonPAReader(associateId, accessKey, secrectKey string) *amazonPAReader {
	client := paapi5.New(
		paapi5.WithMarketplace(paapi5.LocaleJapan),
	).CreateClient(
		associateId,
		accessKey,
		secrectKey,
		paapi5.WithHttpClient(&http.Client{}),
	)
	return &amazonPAReader{
		client: client,
	}
}

func (a *amazonPAReader) GetBookInfo(isbn13 string) (*BookInfo, error) {
	if strings.HasPrefix(isbn13, "978") {
		return a.getBookInfoFromIsbn13(isbn13)
	}
	return a.getBookInfoFromEan13(isbn13)
}

func (a *amazonPAReader) getBookInfoFromEan13(ean13 string) (*BookInfo, error) {
	asin, err := a.getASINFromEan13(ean13)
	if err != nil {
		return nil, err
	}
	q := query.NewGetItems(a.client.Marketplace(), a.client.PartnerTag(), a.client.PartnerType()).
		ASINs([]string{*asin}).
		EnableItemInfo().
		EnableBrowseNodeInfo()

	body, err := a.client.Request(q)
	if err != nil {
		return nil, err
	}
	res, err := entity.DecodeResponse(body)
	if err != nil {
		return nil, err
	}
	return a.getBookInfoFromPAData(res)
}

func (a *amazonPAReader) getASINFromEan13(ean13 string) (*string, error) {
	q := query.NewSearchItems(a.client.Marketplace(), a.client.PartnerTag(), a.client.PartnerType()).
		Search(query.Keywords, ean13).
		EnableItemInfo().
		EnableBrowseNodeInfo()

	body, err := a.client.Request(q)
	if err != nil {
		return nil, err
	}
	res, err := entity.DecodeResponse(body)
	if err != nil {
		return nil, err
	}
	return a.getAsinFromSearchData(res)
}

func (a *amazonPAReader) getBookInfoFromIsbn13(isbn13 string) (*BookInfo, error) {
	// convert to isbn10
	isbn10, err := utils.ConvertToISBN10(isbn13)
	if err != nil {
		return nil, err
	}

	q := query.NewGetItems(a.client.Marketplace(), a.client.PartnerTag(), a.client.PartnerType()).
		ASINs([]string{isbn10}).
		EnableItemInfo().
		EnableBrowseNodeInfo()

	body, err := a.client.Request(q)
	if err != nil {
		return nil, err
	}
	res, err := entity.DecodeResponse(body)
	if err != nil {
		return nil, err
	}
	return a.getBookInfoFromPAData(res)
}

func (a *amazonPAReader) getBookInfoFromPAData(data *entity.Response) (*BookInfo, error) {
	if data.Errors != nil && len(data.Errors) > 0 {
		return nil, fmt.Errorf(data.Errors[0].Message)
	}
	if data.ItemsResult == nil || data.ItemsResult.Items == nil || len(data.ItemsResult.Items) < 1 {
		return nil, fmt.Errorf("no result")
	}
	item := data.ItemsResult.Items[0]

	title := ""
	if item.ItemInfo.Title != nil {
		title = item.ItemInfo.Title.DisplayValue
	}

	publisher := ""
	if item.ItemInfo.ByLineInfo != nil && item.ItemInfo.ByLineInfo.Manufacturer != nil {
		publisher = item.ItemInfo.ByLineInfo.Manufacturer.DisplayValue
	}

	pubDate := ""
	if item.ItemInfo.ContentInfo != nil {
		pubDate = item.ItemInfo.ContentInfo.PublicationDate.DisplayValue.Time.Format("2006-01-02")
	}

	var authors []string
	if item.ItemInfo.ByLineInfo != nil {
		for _, cont := range item.ItemInfo.ByLineInfo.Contributors {
			authors = append(authors, cont.Name)
		}
	}

	kind := ""
	if item.ItemInfo.Classifications != nil {
		kind = item.ItemInfo.Classifications.Binding.DisplayValue
	}

	genre := ""
	if item.BrowseNodeInfo != nil {
		for _, node := range item.BrowseNodeInfo.BrowseNodes {
			if node.DisplayName != "" {
				genre = node.DisplayName
				break
			}
		}
	}

	return NewBookInfo(title, authors, publisher, pubDate, kind, genre), nil
}

func (a *amazonPAReader) getAsinFromSearchData(data *entity.Response) (*string, error) {
	if data.Errors != nil && len(data.Errors) > 0 {
		return nil, fmt.Errorf(data.Errors[0].Message)
	}
	if data.SearchResult == nil || data.SearchResult.Items == nil || len(data.SearchResult.Items) < 1 {
		return nil, fmt.Errorf("no result")
	}
	item := data.SearchResult.Items[0]
	return &item.ASIN, nil
}
