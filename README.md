# go-pacsearch

A Go library for searching Arch Linux packages via [Official Repository Web Interface](https://wiki.archlinux.org/title/Official_repositories_web_interface)

## Examples

```go
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
```

```
community/firefox-tridactyl v1.22.1
Replace Firefox's control mechanism with one modelled on Vim

community/grub-theme-vimix v20190605
A blur theme for grub

extra/gvim v9.0.0321
Vi Improved, a highly configurable, improved version of the vi text editor (with advanced features, such as a GUI)

community/neovide v0.10.1

...
```

## License

[MIT](./LICENSE)
