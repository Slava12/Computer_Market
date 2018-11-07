package content

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
)

// AddRAM добавляет оперативную память из указанного адреса
func AddRAM(url string, filesFolder string) {
	nodes, err := getURLNodes(url)
	if err != nil {
		logger.Warn(err, "Не удалось распарсить html-страницу!")
		return
	}
	featuresNames := [7]string{"Производитель", "Серия", "Модель", "Объем модуля памяти", "Тип памяти", "Частота функционирования", "Количество модулей в комплекте"}
	featuresValues := make([]string, 7)
	featuresValues[0] = strings.TrimSpace(nodes.Find("#tdsa2943").Text())
	featuresValues[1] = strings.TrimSpace(nodes.Find("#tdsa5562").Text())
	featuresValues[2] = strings.TrimSpace(nodes.Find("#tdsa2944").Text())
	temp := strings.Split(strings.TrimSpace(nodes.Find("#tdsa3360").Text()), " ")
	featuresValues[3] = temp[0] + " " + temp[1]
	tempString := strings.TrimSpace(nodes.Find("#tdsa2510").Text())
	if strings.Contains(tempString, "LV") || strings.Contains(tempString, "ECC") { // Убрать низковольтажные и некоторые серверные
		return
	}
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa2510").Text()), " ")
	featuresValues[4] = temp[2]
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa1475").Text()), " ")
	if len(temp) < 3 { // если частота написана с ошибками
		return
	}
	featuresValues[5] = temp[1] + " " + temp[2]
	featuresValues[6] = strings.TrimSpace(nodes.Find("#tdsa4273").Text())
	result := database.Unit{}
	result.Features = make([]database.FeatureUnit, 7)
	for i := range result.Features {
		result.Features[i].Name = featuresNames[i]
		result.Features[i].Value = featuresValues[i]
	}
	if result.Features[1].Value != "" {
		result.Name = result.Features[0].Value + " " + result.Features[1].Value + " " + result.Features[2].Value
	} else {
		result.Name = result.Features[0].Value + " " + result.Features[2].Value
	}
	result.CategoryID = 4
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
	tempString, _ = nodes.Find("#gallery-image-2254").Attr("href")
	if tempString != "" {
		result.Pictures = append(result.Pictures, tempString)
	}
	tempString, _ = nodes.Find("#gallery-image-2245").Attr("href")
	if tempString != "" {
		result.Pictures = append(result.Pictures, tempString)
	}
	tempString, _ = nodes.Find("#gallery-image-2258").Attr("href")
	if tempString != "" {
		result.Pictures = append(result.Pictures, tempString)
	}
	tempString, _ = nodes.Find("#gallery-image-2248").Attr("href")
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

// AddRAMFromURL добавляет оперативную память из указанного URL
func AddRAMFromURL(url string, filesFolder string) {
	nodes, err := getURLNodes(url)
	if err != nil {
		logger.Warn(err, "Не удалось распарсить html-страницу!")
		return
	}
	nodes.Find(".t").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		if !strings.Contains(href, "Registered") || !strings.Contains(href, "ECC") { // Убрать серверные
			url = "https://www.nix.ru/" + href
			AddRAM(url, filesFolder)
		}
	})
}

// AddRAMFromFile добавляет оперативную память из указанного файла
func AddRAMFromFile(fileName string, filesFolder string) {
	urls, err := makeStringsArrayFromFile(fileName)
	if err != nil {
		logger.Info(err, "Не удалось открыть файл!")
		return
	}
	for i := 1; i < len(urls); i++ {
		AddRAMFromURL(urls[i], filesFolder)
	}
}
