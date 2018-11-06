package content

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/files"
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

func saveImageFromURL(url string, path string) {
	response, err := http.Get(url)
	if err != nil {
		logger.Warn(err, "Не удалось обратиться по указанному адресу!")
		return
	}

	defer response.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		logger.Warn(err, "Не удалось создать файл!")
		return
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		logger.Warn(err, "Не удалось скопировать изображение!")
		return
	}
	file.Close()
}

func updateUnitPictures(unit database.Unit, id int, filesFolder string) {
	filePath := filesFolder + strconv.Itoa(id) + "/"
	files.CreateDirectory(filePath)
	pictures := make([]string, len(unit.Pictures))
	copy(pictures, unit.Pictures)
	unit.Pictures = make([]string, len(pictures))
	for i, pic := range pictures {
		segments := strings.Split(pic, "/")
		segments = strings.Split(segments[len(segments)-1], ".")
		fileExt := segments[len(segments)-1]
		name := strconv.Itoa(i) + "." + fileExt
		path := filePath + name
		saveImageFromURL(pic, path)
		unit.Pictures[i] = path
	}
	errUpdate := database.UpdateUnit(id, unit.Name, unit.CategoryID, unit.Quantity, unit.Price, unit.Discount, unit.Popularity, unit.Features, unit.Pictures)
	if errUpdate != nil {
		logger.Warn(errUpdate, "Не удалось обновить информацию о товаре ", id, "!")
		return
	}
	logger.Info("Обновление информации о товаре ", id, " прошло успешно.")
}
