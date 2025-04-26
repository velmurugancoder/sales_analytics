package products

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	global "sales_analytics/Global"
	"sales_analytics/dbconnection"
	helperpkg "sales_analytics/helper_pkg"
	"sales_analytics/sales_analyticsprocess/common"
	"sales_analytics/sales_analyticsprocess/model"
	"strconv"
	"strings"
)

func Get_Productsdetails(w http.ResponseWriter, r *http.Request) {

	log.Println(" Get_Productsdetails (+) ")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Indicator, NINDICATOR, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	if r.Method == http.MethodPost {

		lIndicator := r.Header.Get("Indicator")
		lNINDICATOR := r.Header.Get("NINDICATOR")

		var lLmit int
		var linputRec model.GetDetails
		var lGetRevenueDetailsRec model.GetRevenueDetails
		lGetRevenueDetailsRec.Status = global.SuccessCode

		lErr := json.NewDecoder(r.Body).Decode(&linputRec)
		if lErr != nil {
			helperpkg.LogError(lErr)
			fmt.Fprint(w, helperpkg.GetErrorString("GPD01", lErr.Error()))
			return
		}

		if strings.TrimSpace(linputRec.StartDate) == "" || strings.TrimSpace(linputRec.EndDate) == "" {
			helperpkg.LogError(fmt.Errorf("date field mandatory"))
			fmt.Fprint(w, helperpkg.GetErrorString("GRD02", "Date field mandatory"))
			return
		}

		if lNINDICATOR == "" {
			lLmit = 10
		} else {
			lLmit, lErr = strconv.Atoi(lNINDICATOR)
			if lErr != nil {
				helperpkg.LogError(lErr)
				lLmit = 10
			}
		}

		switch lIndicator {
		case "Overall":
			lGetRevenueDetailsRec.TopProduct, lErr = GetOverallProducts(linputRec, lLmit)

		case "Category":
			lGetRevenueDetailsRec.TopCategory, lErr = GetTopCategories(linputRec, lLmit)

		case "Region":
			lGetRevenueDetailsRec.TopRegion, lErr = GetTopRegions(linputRec, lLmit)
		}

		if lErr != nil {
			helperpkg.LogError(lErr)
			fmt.Fprint(w, helperpkg.GetErrorString("GPD02", lErr.Error()))
			return
		}

		lData, lErr := json.Marshal(lGetRevenueDetailsRec)
		if lErr != nil {
			helperpkg.LogError(lErr)
			fmt.Fprint(w, helperpkg.GetErrorString("GPD03", lErr.Error()))
			return
		}

		fmt.Fprint(w, string(lData))

	}
	log.Println(" Get_Productsdetails (-) ")

}

func GetOverallProducts(pInputRec model.GetDetails, lLimit int) ([]model.TopProduct, error) {

	var lTopProdectArr []model.TopProduct

	lStartDate := common.GetDate(pInputRec.StartDate)
	lEndDate := common.GetDate(pInputRec.EndDate)

	lErr := dbconnection.G_Db_instance.Mysql_gormdb.Table("order_items").
		Select("nvl(products.name,''), nvl(SUM(order_items.quantity_sold ),'') as total_quantity").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.date_of_sale BETWEEN ? AND ?", lStartDate, lEndDate).
		Group("products.name").
		Order("total_quantity DESC").
		Limit(lLimit).
		Scan(&lTopProdectArr).Error

	if lErr != nil {
		helperpkg.LogError(lErr)
		return lTopProdectArr, lErr
	}

	return lTopProdectArr, nil
}

func GetTopCategories(pInputRec model.GetDetails, pLimit int) ([]model.TopCategory, error) {

	log.Println(" GetTopCategories (+) ")

	var lTopCategoryArr []model.TopCategory
	lStartDate := common.GetDate(pInputRec.StartDate)
	lEndDate := common.GetDate(pInputRec.EndDate)

	lErr := dbconnection.G_Db_instance.Mysql_gormdb.Table("order_items").
		Select("nvl(products.category,'') as category, nvl(SUM(order_items.quantity_sold),'') as total_quantity").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.date_of_sale BETWEEN ? AND ?", lStartDate, lEndDate).
		Group("products.category").
		Order("total_quantity DESC").
		Limit(pLimit).
		Scan(&lTopCategoryArr).Error

	if lErr != nil {
		helperpkg.LogError(lErr)
		return lTopCategoryArr, lErr
	}

	log.Println(" GetTopCategories (-) ")

	return lTopCategoryArr, nil
}

func GetTopRegions(pInputRec model.GetDetails, pLimit int) ([]model.TopRegion, error) {
	log.Println(" GetTopRegions (+) ")

	var lTopRegionArr []model.TopRegion
	lStartDate := common.GetDate(pInputRec.StartDate)
	lEndDate := common.GetDate(pInputRec.EndDate)

	lErr := dbconnection.G_Db_instance.Mysql_gormdb.Table("order_items").
		Select("nvl(orders.region,'') as region, nvl(SUM(order_items.quantity_sold),'') as total_quantity").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.date_of_sale BETWEEN ? AND ?", lStartDate, lEndDate).
		Group("orders.region").
		Order("total_quantity DESC").
		Limit(pLimit).
		Scan(&lTopRegionArr).Error

	if lErr != nil {
		helperpkg.LogError(lErr)
		return lTopRegionArr, lErr
	}

	log.Println(" GetTopRegions (-) ")
	return lTopRegionArr, lErr

}
