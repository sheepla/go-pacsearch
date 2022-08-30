//nolint:paralleltest,testpackage,ifshort
package pacsearch

import (
	"errors"
	"testing"
)

func TestOutOfDateString(t *testing.T) {
	if OutOfDateFlagged.String() != "Flagged" {
		t.Error("missmatch `Flagged`")
	}

	if OutOfDateNotFlagged.String() != "Not Flagged" {
		t.Error("missmatch `Not Flagged`")
	}
}

func TestArchString(t *testing.T) {
	if ArchAny.String() != "any" {
		t.Error("missmatch `any`")
	}

	if ArchX8664.String() != "x86_64" {
		t.Error("missmatch `x86_64`")
	}

	if ArchI686.String() != "i686" {
		t.Error("missmatch `i686`")
	}

	if ArchNone.String() != "" {
		t.Error("missmatch ArchNone")
	}
}

func TestRepoString(t *testing.T) {
	if RepoCore.String() != "Core" {
		t.Error("missmatch `Core`")
	}

	if RepoNone.String() != "" {
		t.Error("missmatch RepoNone")
	}
}

func TestNewParam(t *testing.T) {
	_, err := NewParam("")
	if !errors.Is(err, ErrQueryIsEmpty) {
		t.Error("not occurred ErrQueryIsEmpty")
	}

	_, err = NewParam("xx")
	if !errors.Is(err, ErrQueryIsTooShort) {
		t.Error("not occurred ErrQueryIsTooShort")
	}

	moderateLengthQuery := "moderate length query"

	param, err := NewParam(moderateLengthQuery)
	if err != nil {
		t.Error(err)
	}

	if param.Query != moderateLengthQuery {
		t.Error("missmatch query")
	}
}

func TestParamToRawQuery(t *testing.T) {
	param := &Param{
		Query:      "QUERY",
		Name:       "NAME",
		Repo:       RepoCore,
		Desc:       "DESC",
		Arch:       ArchX8664,
		Maintainer: "MAINTAINER",
		Packager:   "PACKAGER",
		OutOfDate:  OutOfDateNotFlagged,
	}

	have := param.ToRawQuery()
	want := "arch=x86_64&desc=DESC&flagged=Not+Flagged&maintainer=MAINTAINER&name=NAME&packager=PACKAGER&q=QUERY&repo=Core"

	if have != want {
		t.Errorf("have=%s, want=%s", have, want)
	}

	param = &Param{
		Query:      "",
		Name:       "",
		Repo:       RepoCore,
		Desc:       "",
		Arch:       ArchAny,
		Maintainer: "",
		Packager:   "",
		OutOfDate:  OutOfDateFlagged,
	}

	have = param.ToRawQuery()
	want = "arch=any&flagged=Flagged&repo=Core"

	if have != want {
		t.Errorf("have=%s, want=%s", have, want)
	}
}

func TestParamToURL(t *testing.T) {
	param := &Param{
		Query:      "QUERY",
		Name:       "NAME",
		Repo:       RepoCore,
		Desc:       "DESC",
		Arch:       ArchX8664,
		Maintainer: "MAINTAINER",
		Packager:   "PACKAGER",
		OutOfDate:  OutOfDateNotFlagged,
	}

	have := param.ToURL()
	want := "https://www.archlinux.org/packages/search/json?arch=x86_64&desc=DESC&flagged=Not+Flagged&maintainer=MAINTAINER&name=NAME&packager=PACKAGER&q=QUERY&repo=Core"

	if have != want {
		t.Errorf("have=%s, want=%s", have, want)
	}
}

func TestFetch(t *testing.T) {
	body, err := fetch("https://www.archlinux.org/packages/search/json?q=vim&limit=3")
	if err != nil {
		t.Error(err)
	}
	defer body.Close()

	t.Log(body)
}

func TestSearch(t *testing.T) {
	body, err := fetch("https://www.archlinux.org/packages/search/json?q=vim&limit=3")
	if err != nil {
		t.Error(err)
	}
	defer body.Close()

	result, err := decodeToSearchResult(body)
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}
