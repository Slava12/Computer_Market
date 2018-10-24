package content

import (
	"bufio"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
	"golang.org/x/net/html/charset"
)

// AddFeature добавляет характеристику товара
func AddFeature(name string) {
	id, err := database.NewFeature(name)
	if err != nil {
		logger.Warn(err, "Не удалось добавить новую характеристику!")
		return
	}
	logger.Info("Добавление характеристики ", id, " прошло успешно.")
}

func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func makeStringsArrayFromFile(fileName string) ([]string, error) {
	values := make([]string, 0)
	file, err := os.Open(fileName)
	if err != nil {
		return values, err
	}
	r := bufio.NewReader(file)
	s, e := readln(r)
	for e == nil {
		values = append(values, s)
		s, e = readln(r)
	}
	file.Close()
	return values, nil
}

// AddFeaturesFromFile добавляет характеристики из файла
func AddFeaturesFromFile(fileName string) {
	features, err := makeStringsArrayFromFile(fileName)
	if err != nil {
		logger.Warn(err, "Не удалось открыть файл!")
		return
	}
	for i := 1; i < len(features); i++ {
		AddFeature(features[i])
	}
}

func getURLNodes(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return &goquery.Document{}, err
	}

	defer resp.Body.Close()
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return &goquery.Document{}, err
	}
	doc, err := goquery.NewDocumentFromReader(utf8)
	if err != nil {
		return &goquery.Document{}, err
	}
	return doc, nil
}

// AddMotherboard добавляет материнскую плату из указанного адреса
func AddMotherboard(url string) {
	nodes, err := getURLNodes(url)
	if err != nil {
		logger.Warn(err, "Не удалось распарсить html-страницу!")
		return
	}
	if strings.TrimSpace(nodes.Find("#tdsa3260").Text()) == "Сервер" {
		return
	}
	featuresNames := [11]string{"Производитель", "Модель", "Описание", "Гнездо процессора", "Чипсет мат. платы", "Формат платы", "Тип поддерживаемой памяти", "Количество разъемов памяти", "Версия PCI Express", "Количество разъёмов PCI Express", "Max объем оперативной памяти"}
	featuresValues := make([]string, 11)
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
	featuresValues[3] = strings.TrimSpace(nodes.Find("#tdsa1307").Text())
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa3362").Text()), " ")
	for i := range temp {
		if i == len(temp)-3 {
			featuresValues[4] += temp[i]
			break
		}
		featuresValues[4] += temp[i] + " "
	}
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa643").Text()), " ")
	featuresValues[5] = temp[0]
	if strings.TrimSpace(nodes.Find("#tds7068").Text()) != "" {
		temp = strings.Split(strings.TrimSpace(nodes.Find("#tds7068").Text()), " ")
		featuresValues[6] = temp[2]
		temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa7068").Text()), " ")
		featuresValues[7] = temp[0]
	}
	if strings.TrimSpace(nodes.Find("#tds3148").Text()) != "" {
		temp = strings.Split(strings.TrimSpace(nodes.Find("#tds3148").Text()), " ")
		featuresValues[6] = temp[2]
		temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa3148").Text()), " ")
		featuresValues[7] = temp[0]
	}
	if strings.TrimSpace(nodes.Find("#tdsa2158").Text()) == "" || !strings.Contains(nodes.Find("#tdsa2158").Text(), "PCI") { // В материнской плате нет PCI-E 16x или не указана версия
		return
	}
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa2158").Text()), " ")
	featuresValues[8] = temp[4]
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa2158").Text()), " ")
	featuresValues[9] = temp[0]
	temp = strings.Split(strings.TrimSpace(nodes.Find("#tdsa894").Text()), " ")
	featuresValues[10] = temp[0]
	result := database.Unit{}
	result.Features = make([]database.FeatureUnit, 11)
	for i := range result.Features {
		result.Features[i].Name = featuresNames[i]
		result.Features[i].Value = featuresValues[i]
	}
	result.Name = featuresValues[0] + " " + featuresValues[1]
	result.CategoryID = 2
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
	result.Pictures = append(result.Pictures, tempString)
	tempString, _ = nodes.Find("#gallery-image-15675").Attr("href")
	result.Pictures = append(result.Pictures, tempString)
	tempString, _ = nodes.Find("#gallery-image-2245").Attr("href")
	result.Pictures = append(result.Pictures, tempString)
	id, errAdd := database.NewUnit(result.Name, result.CategoryID, result.Quantity, result.Price, result.Discount, result.Features, result.Pictures)
	if errAdd != nil {
		logger.Warn(errAdd, "Не удалось добавить новый товар!")
		return
	}
	logger.Info("Добавление товара ", id, " прошло успешно.")
}

// AddMotherboardsFromURL добавляет материнские платы из указанного URL
func AddMotherboardsFromURL(url string) {
	nodes, err := getURLNodes(url)
	if err != nil {
		logger.Warn(err, "Не удалось распарсить html-страницу!")
		return
	}
	nodes.Find(".t").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		if !strings.Contains(href, "onboard") && strings.Contains(href, "PCI-E") { // Убрать с предустановленным процессором и без поддержки PCI-Express
			url = "https://www.nix.ru/" + href
			AddMotherboard(url)
		}
	})
}

// AddMotherboardsFromFile добавляет материнские платы из указанного файла
func AddMotherboardsFromFile(fileName string) {
	urls, err := makeStringsArrayFromFile(fileName)
	if err != nil {
		logger.Info(err, "Не удалось открыть файл!")
		return
	}
	for i := 1; i < len(urls); i++ {
		AddMotherboardsFromURL(urls[i])
	}
}
