package pacsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Param struct {
	Query      string    `name:"q"`
	Name       string    `name:"name"`
	Repo       Repo      `name:"repo"`
	Desc       string    `name:"desc"`
	Arch       Arch      `name:"arch"`
	Maintainer string    `name:"maintainer"`
	Packager   string    `name:"packager"`
	OutOfDate  OutOfDate `name:"flagged"`
}

type OutOfDate int

const (
	OutOfDateFlagged OutOfDate = iota + 1
	OutOfDateNotFlagged
)

func (o OutOfDate) String() string {
	switch o {
	case OutOfDateFlagged:
		return "Flagged"
	case OutOfDateNotFlagged:
		return "Not Flagged"
	default:
		return ""
	}
}

type Arch int

const (
	ArchAny   Arch = iota + 1
	ArchX8664 Arch = iota + 1
	ArchI686  Arch = iota + 1
	ArchNone  Arch = iota + 1
)

func (a Arch) String() string {
	switch a {
	case ArchAny:
		return "any"
	case ArchX8664:
		return "x86_64"
	case ArchI686:
		return "i686"
	case ArchNone:
		return ""
	default:
		return ""
	}
}

type Repo int

const (
	RepoCore Repo = iota
	RepoExtra
	RepoTesting
	RepoMultilib
	RepoMultilibTesting
	RepoCommunity
	RepoCommunityTesting
	RepoNone
)

func (r Repo) String() string {
	switch r {
	case RepoCore:
		return "Core"
	case RepoExtra:
		return "Extra"
	case RepoTesting:
		return "Testing"
	case RepoMultilib:
		return "Multilib"
	case RepoMultilibTesting:
		return "Multilib-Testing"
	case RepoCommunity:
		return "Communitiy"
	case RepoCommunityTesting:
		return "Communitiy-Testing"
	case RepoNone:
		return ""
	default:
		return ""
	}
}

//nolint:tagliatelle
type SearchResult struct {
	Limit    int64 `json:"limit"`
	NumPages int64 `json:"num_pages"`
	Page     int64 `json:"page"`
	Results  []struct {
		Arch           string   `json:"arch"`
		BuildDate      string   `json:"build_date"`
		Checkdepends   []string `json:"checkdepends"`
		CompressedSize int64    `json:"compressed_size"`
		Conflicts      []string `json:"conflicts"`
		Depends        []string `json:"depends"`
		Epoch          int64    `json:"epoch"`
		Filename       string   `json:"filename"`
		FlagDate       string   `json:"flag_date"`
		Groups         []string `json:"groups"`
		InstalledSize  int64    `json:"installed_size"`
		LastUpdate     string   `json:"last_update"`
		Licenses       []string `json:"licenses"`
		Maintainers    []string `json:"maintainers"`
		Makedepends    []string `json:"makedepends"`
		Optdepends     []string `json:"optdepends"`
		Packager       string   `json:"packager"`
		Pkgbase        string   `json:"pkgbase"`
		Pkgdesc        string   `json:"pkgdesc"`
		Pkgname        string   `json:"pkgname"`
		Pkgrel         string   `json:"pkgrel"`
		Pkgver         string   `json:"pkgver"`
		Provides       []string `json:"provides"`
		Replaces       []string `json:"replaces"`
		Repo           string   `json:"repo"`
		URL            string   `json:"url"`
	} `json:"results"`
	Valid   bool  `json:"valid"`
	Version int64 `json:"version"`
}

var (
	ErrQueryIsEmpty    = errors.New("query is empty")
	ErrQueryIsTooShort = errors.New("query is too short")
)

const minQueryLength = 2

//nolint:varnamelen
func NewParam(query string) (*Param, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, ErrQueryIsEmpty
	}

	if len(q) <= minQueryLength {
		return nil, ErrQueryIsTooShort
	}

	return &Param{
		Query:      q,
		Name:       "",
		Desc:       "",
		Repo:       RepoNone,
		Arch:       ArchNone,
		Maintainer: "",
		Packager:   "",
		OutOfDate:  OutOfDateNotFlagged,
	}, nil
}

func (param *Param) ToRawQuery() string {
	if param == nil {
		return ""
	}

	query := url.Values{}

	if !isEmptyString(param.Query) {
		query.Add("q", param.Query)
	}

	if !isEmptyString(param.Name) {
		query.Add("name", param.Name)
	}

	if !isEmptyString(param.Repo.String()) {
		query.Add("repo", param.Repo.String())
	}

	if !isEmptyString(param.Arch.String()) {
		query.Add("arch", param.Arch.String())
	}

	if !isEmptyString(param.Desc) {
		query.Add("desc", param.Desc)
	}

	if !isEmptyString(param.Maintainer) {
		query.Add("maintainer", param.Maintainer)
	}

	if !isEmptyString(param.Packager) {
		query.Add("packager", param.Packager)
	}

	if !isEmptyString(param.OutOfDate.String()) {
		query.Add("flagged", param.OutOfDate.String())
	}

	return query.Encode()
}

func (param *Param) ToURL() string {
	//nolint:exhaustivestruct,exhaustruct,varnamelen
	u := &url.URL{
		Scheme:   "https",
		Host:     "www.archlinux.org",
		Path:     "packages/search/json",
		RawQuery: param.ToRawQuery(),
	}

	return u.String()
}

func (param *Param) Search() (*SearchResult, error) {
	body, err := fetch(param.ToURL())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the result: %w", err)
	}
	defer body.Close()

	result, err := decodeToSearchResult(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the search result: %w", err)
	}

	return result, nil
}

const timeout = 10 * time.Second

func fetch(url string) (io.ReadCloser, error) {
	//nolint:exhaustivestruct,exhaustruct
	client := &http.Client{Timeout: timeout}

	//nolint:noctx
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		if res.StatusCode < 200 || 300 < res.StatusCode {
			return nil, fmt.Errorf("HTTP error (%s): %w", res.Status, err)
		}

		return nil, fmt.Errorf("failed to get the response: %w", err)
	}

	return res.Body, nil
}

func decodeToSearchResult(body io.Reader) (*SearchResult, error) {
	var result SearchResult
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &result, nil
}

func isEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
