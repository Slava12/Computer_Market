package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
)

func createPairsFromAllUnits(w http.ResponseWriter, r *http.Request) {
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

func createPair(id1 int, id2 int) {
	if id1 == id2 {
		return
	}
	pair, _ := database.GetPairByUnitsID(id1, id2)
	if pair.ID == 0 {
		pair, _ = database.GetPairByUnitsID(id2, id1)
		if pair.ID == 0 {
			if id1 < id2 {
				id, err := database.NewPair(id1, id2, 1)
				if err != nil {
					logger.Warn(err, "Не удалось добавить новую пару!")
					return
				}
				logger.Info("Добавление пары ", id, " прошло успешно.")
			} else {
				id, err := database.NewPair(id2, id1, 1)
				if err != nil {
					logger.Warn(err, "Не удалось добавить новую пару!")
					return
				}
				logger.Info("Добавление пары ", id, " прошло успешно.")
			}
			return
		}
	}
	err := database.UpdatePairCount(pair.ID, pair.Count+1)
	if err != nil {
		logger.Warn(err, "Не удалось обновить пару ", pair.ID, "!")
		return
	}
	logger.Info("Обновление пары ", pair.ID, " прошло успешно.")
}

func createPairs(units []int) {
	for i := 0; i < len(units)-1; i++ {
		for j := i + 1; j < len(units); j++ {
			createPair(units[i], units[j])
		}
	}
}
