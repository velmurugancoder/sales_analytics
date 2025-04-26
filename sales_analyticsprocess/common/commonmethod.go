package common

import (
	helperpkg "sales_analytics/helper_pkg"
	"time"
)

func GetDate(pdateStr string) time.Time {

	lDate, lErr := time.Parse("2006-01-02", pdateStr)
	if lErr != nil {
		helperpkg.LogError(lErr)
	}

	return lDate
}
