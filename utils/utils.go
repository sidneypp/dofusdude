package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

var Languages = []string{"de", "en", "es", "fr", "it", "pt"}

var ApiHostName string
var ApiPort string
var ApiScheme string

type VersionT struct {
	Search bool
	MemDb  bool
}

func ReadEnvs() (string, string) {
	apiScheme, ok := os.LookupEnv("API_SCHEME")
	if !ok {
		apiScheme = "http"
	}
	ApiScheme = apiScheme

	apiHostName, ok := os.LookupEnv("API_HOSTNAME")
	if !ok {
		apiHostName = "localhost"
	}
	ApiHostName = apiHostName

	apiPort, ok := os.LookupEnv("API_PORT")
	if !ok {
		apiPort = "3000"
	}
	ApiPort = apiPort

	return ApiHostName, ApiPort
}

func ItemImageUrl(itemId int) string {
	return fmt.Sprintf("%s/img/item/%d.png", ApiHostName, itemId)
}

func MountImageUrl(itemId int) string {
	return fmt.Sprintf("%s/img/item/%d.png", ApiHostName, itemId)
}

func CreateMeiliClient() *meilisearch.Client {
	meiliPort, ok := os.LookupEnv("MEILISEARCH_PORT")
	if !ok {
		meiliPort = "7700"
	}

	meiliKey, ok := os.LookupEnv("MEILISEARCH_API_KEY")
	if !ok {
		meiliKey = "masterKey"
	}

	meiliHost, ok := os.LookupEnv("MEILISEARCH_HOST")
	if !ok {
		meiliHost = "http://127.0.0.1"
	}

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   fmt.Sprintf("%s:%s", meiliHost, meiliPort),
		APIKey: meiliKey,
	})

	return client
}

func CreateDataDirectoryStructure() {
	os.MkdirAll("data/tmp", 0755)
	os.MkdirAll("data/img/item", 0755)
	os.MkdirAll("data/languages", 0755)
}

type Config struct {
	CurrentVersion string `json:"currentDofusVersion"`
}

func GetConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}

	return config
}

func SaveConfig(config Config, path string) error {
	configJson, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, configJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

func CurrentRedBlueVersionStr(redBlueValue bool) string {
	if redBlueValue {
		return "red"
	}
	return "blue"
}

func NextRedBlueVersionStr(redBlueValue bool) string {
	if redBlueValue {
		return "blue"
	}
	return "red"
}

type Pagination struct {
	PageNumber  int
	PageSize int

	BiggestPageSize int
}

type PaginationLinks struct {
	First string `json:"first"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	Last  string `json:"last"`
}

func (p *Pagination) ValidatePagination(listSize int) int {
	if p.PageSize > p.BiggestPageSize {
		return -1
	}
	if p.BiggestPageSize * p.PageNumber > listSize + p.BiggestPageSize {
		return 1
	}
	return 0
}

func (p *Pagination) BuildLinks(listSize int) (PaginationLinks, bool) {
	mainUrl := url.URL{
        Scheme: ApiScheme,
        Host:   ApiHostName,
    }

	firstPage := 1
	var lastPage int

	lastPageSize := listSize % p.PageSize
	if lastPageSize == 0 {
		lastPage = listSize / p.PageSize
	} else {
		lastPage = (listSize / p.PageSize) + 1
	}
	
	firstUrlQuery := mainUrl.Query()
	firstUrlQuery.Set("pnum", fmt.Sprint(firstPage))
	firstUrlQuery.Set("psize", fmt.Sprint(p.PageSize))

	prevUrlQuery := mainUrl.Query()
	prevUrlQuery.Set("pnum", fmt.Sprint(p.PageNumber - 1))
	prevUrlQuery.Set("psize", fmt.Sprint(p.PageSize))

	nextUrlQuery := mainUrl.Query()
	nextUrlQuery.Set("pnum", fmt.Sprint(p.PageNumber + 1))
	nextUrlQuery.Set("psize", fmt.Sprint(p.PageSize))

	lastUrlQuery := mainUrl.Query()
	lastUrlQuery.Set("pnum", fmt.Sprint(lastPage))
	lastUrlQuery.Set("psize", fmt.Sprint(p.PageSize))

	firstUrl := mainUrl
	prevUrl := mainUrl
	nextUrl := mainUrl
	lastUrl := mainUrl

	firstUrl.RawQuery = firstUrlQuery.Encode()
	prevUrl.RawQuery = prevUrlQuery.Encode()
	nextUrl.RawQuery = nextUrlQuery.Encode()
	lastUrl.RawQuery = lastUrlQuery.Encode()

	firstUrlStr := firstUrl.String()
	prevUrlStr := prevUrl.String()
	nextUrlStr := nextUrl.String()
	lastUrlStr := lastUrl.String()

	if p.PageNumber == firstPage {
		prevUrlStr = ""
	}
	if p.PageNumber == lastPage {
		nextUrlStr = ""
	}

	return PaginationLinks{
		First: firstUrlStr,
		Prev:  prevUrlStr,
		Next:  nextUrlStr,
		Last:  lastUrlStr,
	}, firstUrl == lastUrl
}

// calculate the start and end index for the pagination
func (p *Pagination) CalculateStartEndIndex(listSize int) (int, int) {
	startIndex := (p.PageNumber - 1) * p.PageSize
	endIndex := startIndex + p.PageSize
	if endIndex > listSize {
		endIndex = listSize
	}
	return startIndex, endIndex
}

func WithValues(ctx context.Context, kv ...interface{}) context.Context {
    if len(kv)%2 != 0 {
        panic("odd numbers of key-value pairs")
    }
    for i := 0; i < len(kv); i = i + 2 {
        ctx = context.WithValue(ctx, kv[i], kv[i+1])
    }
    return ctx
}
