package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sales_analytics/dbconnection"
	helperpkg "sales_analytics/helper_pkg"
	salesanalyticsprocess "sales_analytics/sales_analyticsprocess"
	"sales_analytics/sales_analyticsprocess/products"
	"sales_analytics/sales_analyticsprocess/revenue"
	"sales_analytics/tomlreader"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	log.Println("Server started Port NUmber... 23434")

	lFile, lErr := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if lErr != nil {
		log.Fatalf("error opening file: %v", lErr)
	}
	defer lFile.Close()

	log.SetOutput(lFile)

	// Database connection process
	lErr = dbconnection.BuildConnection()
	if lErr != nil {
		log.Fatal(lErr)
	}
	defer dbconnection.G_Db_instance.Mysql_sqldb.Close()

	lAutorunConfig := tomlreader.ReadTomlFile("./toml/serviceconfig.toml")
	lisAutorundaily := fmt.Sprintf("%v", lAutorunConfig.(map[string]interface{})["AutoRun"])
	if lisAutorundaily == "Y" {
		go DaiyFileUploadProcess()
	}

	lRouter := mux.NewRouter()
	lRouter.HandleFunc("/Uploadfiledetails", salesanalyticsprocess.Uploadfiledetails).Methods(http.MethodGet)
	lRouter.HandleFunc("/based_onRevenue", revenue.Get_revenuedetails).Methods(http.MethodPost)
	lRouter.HandleFunc("/Get_Productsdetails", products.Get_Productsdetails).Methods(http.MethodPost)

	lSrv := &http.Server{
		Handler: lRouter,
		Addr:    ":23434",
	}

	log.Fatal(lSrv.ListenAndServe())

}

func DaiyFileUploadProcess() {
	for {

		lNow := time.Now()

		lTimeConfig := tomlreader.ReadTomlFile("./toml/serviceconfig.toml")
		lHour := fmt.Sprintf("%v", lTimeConfig.(map[string]interface{})["hour"])
		lminute := fmt.Sprintf("%v", lTimeConfig.(map[string]interface{})["minute"])

		lHour_int, lErr := strconv.Atoi(lHour)
		if lErr != nil {
			lHour_int = 8
		}
		lminute_int, lErr := strconv.Atoi(lminute)
		if lErr != nil {
			lHour_int = 0
		}

		if lNow.Hour() == lHour_int && lNow.Minute() == lminute_int {

			lErr := salesanalyticsprocess.ReadFile_updatedata()
			if lErr != nil {
				helperpkg.LogError(lErr)
			}

			time.Sleep(61 * time.Second)

		} else {
			time.Sleep(30 * time.Second)
		}

	}
}
