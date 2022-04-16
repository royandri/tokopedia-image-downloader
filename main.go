package main

import (
	"fmt"
	"tokopedia/tokopedia"
)

func main() {
	limit := 200
	shopId := "11422428"
	tokopedia := tokopedia.TokopediaImageScrapper(shopId, &limit)
	err := tokopedia.DownloadProductImages()

	if err != nil {
		fmt.Println(err)
	}
}
