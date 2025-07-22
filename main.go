package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ScrapeTarget struct {
	Latitude      string
	Longitude     string
	CategoryID    int
	SubcategoryID int
}

type Product struct {
	ID       string
	Name     string
	ImageURL string
}

type BlinkitResponse struct {
	IsSuccess bool `json:"is_success"`
	Response  struct {
		Snippets []struct {
			Data struct {
				Identity struct {
					ID string `json:"id"`
				} `json:"identity"`
				Name struct {
					Text string `json:"text"`
				} `json:"name"`
				Image struct {
					URL string `json:"url"`
				} `json:"image"`
			} `json:"data"`
		} `json:"snippets"`
	} `json:"response"`
}

func fetchProducts(input ScrapeTarget, offset, limit int) ([]Product, error) {
	url := fmt.Sprintf(
		"https://blinkit.com/v1/layout/listing_widgets?offset=%d&limit=%d&exclude_combos=false&l0_cat=%d&l1_cat=%d&last_snippet_type=product_card_snippet_type_2&last_widget_type=product_container&oos_visibility=true&page_index=1&total_entities_processed=1&total_pagination_items=41",
		offset, limit, input.CategoryID, input.SubcategoryID,
	)

	payload := `{
        "applied_filters": null,
        "is_sr_rail_visible": false,
        "is_subsequent_page": false,
        "postback_meta": {},
        "sort": ""
    }`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("access_token", "null")
	req.Header.Set("app_client", "consumer_web")
	req.Header.Set("app_version", "1010101010")
	req.Header.Set("auth_key", "c761ec3633c22afad934fb17a66385c1c06c5472b4898b866b7306186d0bb477")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("device_id", "8fd421f4-e057-46ae-b8c3-4761645c565b")
	req.Header.Set("lat", input.Latitude)
	req.Header.Set("lon", input.Longitude)
	req.Header.Set("origin", "https://blinkit.com")
	req.Header.Set("platform", "mobile_web")
	req.Header.Set("referer", "https://blinkit.com")
	req.Header.Set("user-agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result BlinkitResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var products []Product
	for _, snippet := range result.Response.Snippets {
		data := snippet.Data
		product := Product{
			ID:       data.Identity.ID,
			Name:     data.Name.Text,
			ImageURL: data.Image.URL,
		}
		products = append(products, product)
	}

	return products, nil
}

func saveToCSV(products []Product, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Name", "Image URL"})
	for _, p := range products {
		writer.Write([]string{
			p.ID,
			p.Name,
			p.ImageURL,
		})
	}
	return nil
}

func main() {
	fmt.Println(" Scraping BlinkIt Subcategories...")

	targets := []ScrapeTarget{
		{Latitude: "12.9266817", Longitude: "77.6690633", CategoryID: 1237, SubcategoryID: 316},
	}

	var allProducts []Product

	for _, target := range targets {
		fmt.Printf("📦 Fetching: category %d, subcategory %d, location (%s, %s)\n",
			target.CategoryID, target.SubcategoryID, target.Latitude, target.Longitude)

		products, err := fetchProducts(target, 0, 30)
		if err != nil {
			fmt.Println(" Error:", err)
			continue
		}
		allProducts = append(allProducts, products...)
	}

	if len(allProducts) == 0 {
		fmt.Println(" No products fetched.")
		return
	}

	err := saveToCSV(allProducts, "blinkit_products.csv")
	if err != nil {
		fmt.Println(" Failed to write CSV:", err)
	} else {
		fmt.Println(" Output saved to blinkit_products.csv")
	}
}