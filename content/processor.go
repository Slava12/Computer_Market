package content

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
)

// AddProcessor добавляет процессор из указанного адреса
func AddProcessor(url string, filesFolder string) {
	nodes, err := getURLNodes(url)
	if err != nil {
		logger.Warn(err, "Не удалось распарсить html-страницу!")
		return
	}
	if strings.TrimSpace(nodes.Find("#tdsa3260").Text()) == "Сервер" {
		return
	}
	featuresNames := [13]string{"Производитель", "Модель", "Описание", "Частота работы процессора", "Гнездо процессора", "Ядро", "Количество ядер", "Кэш L1", "Кэш L2", "Кэш L3", "Видеоядро процессора", "Max объем оперативной памяти", "Техпроцесс"}
	featuresValues := make([]string, 13)
	featuresValues[0] = strings.TrimSpace(nodes.Find("#tdsa2943").Text())
	temp := strings.Split(strings.TrimSpace(nodes.Find("#tdsa2944").Text()), " ")
	for i := range temp {
		if i == len(temp)-4 {
			featuresValues[1] += temp[i]
			break
		}
		featuresValues[1] += temp[i] + " "
	}
	featuresValues[2] = strings.TrimSpace(nodes.Find("#tdsa934").Text())
	tempString := strings.TrimSpace(nodes.Find("#tdsa3808").Text())
	if strings.Contains(tempString, ",") {
		temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa3808").Text()), ",")
		featuresValues[3] = temp[0]
	} else {
		temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa3808").Text()), " ")
		featuresValues[3] = temp[0] + " " + temp[1]
	}
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa1307").Text()), " ")
	featuresValues[4] = temp[0] + " " + temp[1]
	tempString = strings.TrimSpace(nodes.Find("#tdsa1549").Text())
	if strings.Contains(tempString, "характеристики") {
		temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa1549").Text()), " ")
		for i := range temp {
			if i == len(temp)-4 {
				featuresValues[5] += temp[i]
				break
			}
			featuresValues[5] += temp[i] + " "
		}
	} else {
		featuresValues[5] = tempString
	}
	featuresValues[6] = strings.TrimSpace(nodes.Find("#tdsa2557").Text())
	featuresValues[7] = strings.TrimSpace(nodes.Find("#tdsa869").Text())
	featuresValues[8] = strings.TrimSpace(nodes.Find("#tdsa919").Text())
	featuresValues[9] = strings.TrimSpace(nodes.Find("#tdsa1946").Text())
	featuresValues[10] = strings.TrimSpace(nodes.Find("#tdsa5486").Text())
	featuresValues[11] = strings.TrimSpace(nodes.Find("#tdsa894").Text())
	featuresValues[12] = strings.TrimSpace(nodes.Find("#tdsa3735").Text())
	result := database.Unit{}
	result.Features = make([]database.FeatureUnit, 13)
	for i := range result.Features {
		result.Features[i].Name = featuresNames[i]
		result.Features[i].Value = featuresValues[i]
	}
	result.Name = featuresValues[0] + " " + featuresValues[1]
	result.CategoryID = 1
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

// AddProcessorsFromURL добавляет процессоры из указанного URL
func AddProcessorsFromURL(url string, filesFolder string) {
	nodes, err := getURLNodes(url)
	if err != nil {
		logger.Warn(err, "Не удалось распарсить html-страницу!")
		return
	}
	nodes.Find(".t").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		if !strings.Contains(href, "BOX") { // Убрать с типом Box
			url = "https://www.nix.ru/" + href
			AddProcessor(url, filesFolder)
		}
	})
}

// AddProcessorsFromFile добавляет материнские платы из указанного файла
func AddProcessorsFromFile(fileName string, filesFolder string) {
	urls, err := makeStringsArrayFromFile(fileName)
	if err != nil {
		logger.Info(err, "Не удалось открыть файл!")
		return
	}
	for i := 1; i < len(urls); i++ {
		AddProcessorsFromURL(urls[i], filesFolder)
	}
}
