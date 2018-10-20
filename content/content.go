package content

import (
	"bufio"
	"os"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
)

// AddFeatures добавляет характеристики товаров
func AddFeatures(names []string) {
	for _, name := range names {
		id, err := database.NewFeature(name)
		if err != nil {
			logger.Warn(err, "Не удалось добавить новую характеристику!")
			return
		}
		logger.Info("Добавление характеристики ", id, " прошло успешно.")
	}
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func AddFeaturesFromFile(fileName string) {
	features := make([]string, 0)
	file, err := os.Open(fileName)
	if err != nil {
		logger.Warn(err, "Не удалось открыть файл!")
		return
	}
	r := bufio.NewReader(file)
	s, e := Readln(r)
	for e == nil {
		features = append(features, s)
		s, e = Readln(r)
	}
	file.Close()
	for _, feature := range features {
		logger.Info(feature)
	}
}
