package content

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
)

// AddVideocard добавляет видеокарту из указанного адреса
func AddVideocard(url string, filesFolder string) {
	nodes, err := getURLNodes(url)
	if err != nil {
		logger.Warn(err, "Не удалось распарсить html-страницу!")
		return
	}
	featuresNames := [8]string{"Производитель", "Модель", "Описание", "Длина", "Частота GPU", "Видеопамять", "Тип видеопамяти", "Интерфейс"}
	featuresValues := make([]string, 8)
	featuresValues[0] = strings.TrimSpace(nodes.Find("#tdsa2943").Text())
	temp := strings.Split(strings.TrimSpace(nodes.Find("#tdsa2944").Text()), " ")
	for i := range temp {
		if i == len(temp)-4 {
			featuresValues[1] += temp[i]
			break
		}
		featuresValues[1] += temp[i] + " "
	}
	featuresValues[2] = strings.TrimSpace(nodes.Find("#tdsa981").Text())
	featuresValues[3] = strings.TrimSpace(nodes.Find("#tdsa955").Text())
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa4191").Text()), " ")
	featuresValues[4] = temp[0] + " " + temp[1]
	featuresValues[5] = strings.TrimSpace(nodes.Find("#tdsa689").Text())
	featuresValues[6] = strings.TrimSpace(nodes.Find("#tdsa4187").Text())
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa567").Text()), " ")
	featuresValues[7] = temp[0] + " " + temp[1] + " " + temp[2]
	result := database.Unit{}
	result.Features = make([]database.FeatureUnit, 8)
	for i := range result.Features {
		result.Features[i].Name = featuresNames[i]
		result.Features[i].Value = featuresValues[i]
	}
	result.Name = featuresValues[0] + " " + featuresValues[1]
	result.CategoryID = 3
	result.Quantity = 5
	temp = strings.Split(strings.TrimSpace(nodes.Find(".price").First().Text()), " ")
	temp = strings.Split(temp[1], string(160))
	price := ""
	for i := 0; i < len(temp)-1; i++ {
		price += temp[i]
	}
	result.Price, err = strconv.Atoi(price)
	if err != nil {
		logger.Info(err, "Не удалось преобразовать строку в число")
		return
	}
	result.Discount = 0
	tempString, _ := nodes.Find("#gallery-image-2254").Attr("href")
	if tempString != "" {
		result.Pictures = append(result.Pictures, tempString)
	}
	tempString, _ = nodes.Find("#gallery-image-15675").Attr("href")
	if tempString != "" {
		result.Pictures = append(result.Pictures, tempString)
	}
	tempString, _ = nodes.Find("#gallery-image-2245").Attr("href")
	if tempString != "" {
		result.Pictures = append(result.Pictures, tempString)
	}
	if len(result.Pictures) == 0 {
		return
	}
	id, errAdd := database.NewUnit(result.Name, result.CategoryID, result.Quantity, result.Price, result.Discount, result.Popularity, result.Features, result.Pictures)
	if errAdd != nil {
		logger.Warn(errAdd, "Не удалось добавить новый товар!")
		return
	}
	logger.Info("Добавление товара ", id, " прошло успешно.")
	updateUnitPictures(result, id, filesFolder)
}

// AddVideocardsFromURL добавляет видеокарты из указанного URL
func AddVideocardsFromURL(url string, filesFolder string) {
	nodes, err := getURLNodes(url)
	if err != nil {
		logger.Warn(err, "Не удалось распарсить html-страницу!")
		return
	}
	nodes.Find(".t").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		if !strings.Contains(href, "BRIDGE") { // Убрать мосты для видеокарт
			url = "https://www.nix.ru/" + href
			AddVideocard(url, filesFolder)
		}
	})
}

// AddVideocardsFromFile добавляет видеокарты из указанного файла
func AddVideocardsFromFile(fileName string, filesFolder string) {
	urls, err := makeStringsArrayFromFile(fileName)
	if err != nil {
		logger.Info(err, "Не удалось открыть файл!")
		return
	}
	for i := 1; i < len(urls); i++ {
		AddVideocardsFromURL(urls[i], filesFolder)
	}
}
