package tokopedia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

type Tokopedia struct {
	StoreId string
	Limit   *int
}

type response struct {
	Data productData `json:"data"`
}

type productData struct {
	GetShopProduct getShopProduct `json:"GetShopProduct"`
}

type getShopProduct struct {
	Data []product `json:"data"`
}

type product struct {
	PrimaryImage primaryImage `json:"primary_image"`
	Name         string       `json:"name"`
}

type primaryImage struct {
	Original string `json:"original"`
}

type image struct {
	url  string
	name string
}

func (t Tokopedia) GetProducts() ([]product, error) {
	fmt.Println("Mengambil data produk")
	url := "https://gql.tokopedia.com/graphql/ShopProducts"

	client := &http.Client{}
	client.Timeout = time.Second * 15

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(t.ConstructBody()))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "_abck=B1B6E4F6AFF9C079B2E81E63450454ED~-1~YAAQpCE1F7pzBiOAAQAAMOLtKAdJ0sg7yWaPBqXrhBMLmqUDNVJ/hSRqyWSRFXZrBGQdW3irW//oPVR3rJhVp2bzu7jUytpH9SND2cOk0nohqM2gVNj0LMtjwBh6faSgHzDMI80feLLGCFQ2WzayZkpMvy86/mf1xyAxHP1BDpQySr3P4KHIU9Ay0xeXoTr/OWcR56LCsIAk6wQz7X3TOtOAJkP8z0TqCQ/jEsUgMpyE5NfYt0IB1ko3r+7pHKZIz2CHkYoU+IN63mWZoOfBt/Po1mbul+k9S9uRQ3ni1wMjXQ1tv7nuOBO3Q/a3XbXxEwt+mRIXYpaQxnldcRYtg4j5wherX+0Di6OzPNrQFB+qyz1e6xGECGe2EKQ=~-1~-1~-1; bm_sz=3ADAEB3B679EC05A9DB2F20326B6751F~YAAQpCE1F7tzBiOAAQAAMOLtKA/IM219qvzrw4kZUduJFasVBo0Dw7zwk0paI4wI3t2QF3UOjc2aH+GVU+/+QQjv4w3ClQ83ae8FArAY6YdBeZqAfJZWSNJHBMdNbQ5kVSIsYqlRGzwQyvTlci7N90U78PfFTw1Bikcg032OkQewIAcjyJMbIutKZjtrG5MPodQbAWHf7gedDcIMh4IFqLxHe+6dBFhglQ2Wv862V/8BAYaI2gQq2WQLCoz7pkHtQVvs16OKhOcnog8HPyMHR2RdoxD+CzbRcKtoCFGlPayzFoBFaxA=~3360309~3420728")
	res, errRes := client.Do(req)
	if errRes != nil {
		return nil, errRes
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var p []response
	err = json.Unmarshal(resBody, &p)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Berhasil mengambil %d data produk\n", len(p[0].Data.GetShopProduct.Data))
	return p[0].Data.GetShopProduct.Data, nil
}

func (t Tokopedia) ConstructBody() []byte {
	limit := 10
	if t.Limit != nil {
		limit = *t.Limit
	}

	body := []byte(fmt.Sprintf(`[
		{
		  "operationName": "ShopProducts",
		  "variables": {
			"sid": "%s",
			"page": 0,
			"perPage": %d,
			"etalaseId": "etalase",
			"sort": 1
		  },
		  "query": "query ShopProducts($sid: String!, $page: Int, $perPage: Int, $keyword: String, $etalaseId: String, $sort: Int, $user_districtId: String, $user_cityId: String, $user_lat: String, $user_long: String) {\n  GetShopProduct(shopID: $sid, filter: {page: $page, perPage: $perPage, fkeyword: $keyword, fmenu: $etalaseId, sort: $sort, user_districtId: $user_districtId, user_cityId: $user_cityId, user_lat: $user_lat, user_long: $user_long}) {\n    status\n    errors\n    links {\n      prev\n      next\n      __typename\n    }\n    data {\n      name\n      product_url\n      product_id\n      price {\n        text_idr\n        __typename\n      }\n      primary_image {\n        original\n        thumbnail\n        resize300\n        __typename\n      }\n      flags {\n        isSold\n        isPreorder\n        isWholesale\n        isWishlist\n        __typename\n      }\n      campaign {\n        discounted_percentage\n        original_price_fmt\n        start_date\n        end_date\n        __typename\n      }\n      label {\n        color_hex\n        content\n        __typename\n      }\n      label_groups {\n        position\n        title\n        type\n        url\n        __typename\n      }\n      badge {\n        title\n        image_url\n        __typename\n      }\n      stats {\n        reviewCount\n        rating\n        __typename\n      }\n      category {\n        id\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n"
		}
	  ]`, t.StoreId, limit))

	return body
}

func (t Tokopedia) ConstructProductImageURLs() ([]image, error) {
	products, err := t.GetProducts()

	if err != nil {
		fmt.Println("Gagal mengambil data produk")
		return nil, err
	}

	fmt.Println("Rekonstruksi gambar produk")
	var finalProducts []image
	for _, product := range products {
		finalProducts = append(finalProducts, image{
			name: slug.Make(product.Name),
			url:  strings.Replace(product.PrimaryImage.Original, "300-square", "900", -1),
		})
	}

	return finalProducts, nil
}

func (t Tokopedia) DownloadProductImages() error {
	productImages, err := t.ConstructProductImageURLs()

	if err != nil {
		return err
	}

	for i, image := range productImages {
		err := t.DownloadImage(image)
		if err != nil {
			fmt.Printf("Gagal mendownload image: %s, cause: %v\n", image.url, err)
		} else {
			fmt.Printf("%d/%d gambar selesai di proses\n", i+1, len(productImages))
		}
	}

	return nil
}

func (t Tokopedia) DownloadImage(image image) error {
	resp, err := http.Get(image.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fmt.Sprintf("./result/tokopedia/%s.jpg", image.name))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func TokopediaImageScrapper(storeId string, limit *int) Tokopedia {
	return Tokopedia{
		StoreId: storeId,
		Limit:   limit,
	}
}
