package salesanalyticsprocess

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	global "sales_analytics/Global"
	helperpkg "sales_analytics/helper_pkg"
	"sales_analytics/sales_analyticsprocess/model"
	"sales_analytics/sales_analyticsprocess/readfile"
	"sales_analytics/tomlreader"
)

func Uploadfiledetails(w http.ResponseWriter, r *http.Request) {
	var lRespRec model.Response
	lRespRec.Status = global.SuccessCode
	lRespRec.Status = "File uploaded"

	lErr := ReadFile_updatedata()
	if lErr != nil {
		helperpkg.LogError(lErr)
		fmt.Fprint(w, helperpkg.GetErrorString("SYD01", lErr.Error()))
		return
	}

	lData, lErr := json.Marshal(lRespRec)
	if lErr != nil {
		helperpkg.LogError(lErr)
		fmt.Fprint(w, helperpkg.GetErrorString("SYD01", lErr.Error()))
		return
	}

	fmt.Fprint(w, string(lData))
}

func ReadFile_updatedata() error {
	log.Println("ReadFile_updatedata (+) ")

	lpathConfig := tomlreader.ReadTomlFile("./toml/filereadconfig.toml")
	lFilepath := fmt.Sprintf("%v", lpathConfig.(map[string]interface{})["FileReadyPath"])

	lErr := readfile.CsvFile_Reader(lFilepath)
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lErr
	}

	log.Println("ReadFile_updatedata (-) ")
	return nil

}
