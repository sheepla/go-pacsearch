package pacsearch

import (
	"errors"
	"net/url"
	"strings"
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

func isEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
