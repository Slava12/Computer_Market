package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
)

func makePairs(w http.ResponseWriter, r *http.Request) {
	units, err := database.GetAllUnits()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список товаров!")
		return
	}
	for i := 1; i < len(units); i++ {
		for j := i + 1; j < i+6; j++ {
			if j == len(units)-1 {
				continue
			}
			count := 1
			id, err := database.NewPair(i, j, count)
			if err != nil {
				logger.Warn(err, "Не удалось добавить новую пару!")
				return
			}
			logger.Info("Добавление пары ", id, " прошло успешно.")
		}
	}
}
