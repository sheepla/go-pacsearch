package main

import (
	"fmt"
	"log"

	"github.com/sheepla/pacsearch"
)

func main() {
	param, err := pacsearch.NewParam("vim")
	if err != nil {
		log.Fatalln(err)
	}

	results, err := pacsearch.Search(param)
	if err != nil {
		log.Fatalln(err)
	}

	for _, result := range results.Results {
		//nolint:forbidigo
		fmt.Printf(
			"%s/%s v%s\n%s\n\n",
			result.Repo,
			result.Pkgname,
			result.Pkgver,
			result.Pkgdesc,
		)
	}
}
