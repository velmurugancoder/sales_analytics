package revenue

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
	"strings"
)

func Get_revenuedetails(w http.ResponseWriter, r *http.Request) {

	log.Println("Get_revenuedetails (+) ")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Indicator, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	if r.Method == http.MethodPost {

		lIndicator := r.Header.Get("Indicator")
		var linputRec model.GetDetails
		var lGetRevenueDetailsRec model.GetRevenueDetails
		lGetRevenueDetailsRec.Status = global.SuccessCode

		lErr := json.NewDecoder(r.Body).Decode(&linputRec)
		if lErr != nil {
			helperpkg.LogError(lErr)
			fmt.Fprint(w, helperpkg.GetErrorString("Invalid input GRD01", lErr.Error()))
			return
		}

		if strings.TrimSpace(linputRec.StartDate) == "" || strings.TrimSpace(linputRec.EndDate) == "" {
			helperpkg.LogError(fmt.Errorf("date field mandatory"))
			fmt.Fprint(w, helperpkg.GetErrorString("GRD02", "Date field mandatory"))
			return
		}

		switch lIndicator {
		case "Date_range":
			lGetRevenueDetailsRec.Total_revenue, lErr = GetTotalRevenue(linputRec)

		case "Product":
			lGetRevenueDetailsRec.TotProdRevenue, lErr = GetRevenuebyProd(linputRec)

		case "Category":
			lGetRevenueDetailsRec.TotalcatRevenue, lErr = GetRevenuebyCat(linputRec)

		case "Region":
			lGetRevenueDetailsRec.TotalRevenue_byreg, lErr = GetRevenuebyregion(linputRec)
		}

		if lErr != nil {
			helperpkg.LogError(lErr)
			fmt.Fprint(w, helperpkg.GetErrorString("GRD03", lErr.Error()))
			return
		}

		lData, lErr := json.Marshal(lGetRevenueDetailsRec)
		if lErr != nil {
			helperpkg.LogError(lErr)
			fmt.Fprint(w, helperpkg.GetErrorString("GRD04", lErr.Error()))
			return
		}

		fmt.Fprint(w, string(lData))

	}

	log.Println("Get_revenuedetails (-) ")

}

func GetTotalRevenue(pInputRec model.GetDetails) (string, error) {
	var ltotalRevenue string

	lStartDate := common.GetDate(pInputRec.StartDate)
	lEndDate := common.GetDate(pInputRec.EndDate)

	lErr := dbconnection.G_Db_instance.Mysql_gormdb.Debug().Table("orders AS o").
		Select("nvl(SUM((oi.quantity_sold * oi.unit_price * (1 - oi.discount)) + o.shipping_cost),'') AS total_revenue").
		Joins("JOIN order_items AS oi ON o.id = oi.order_id").
		Where("o.date_of_sale BETWEEN ? AND ?", lStartDate, lEndDate).
		Scan(&ltotalRevenue).Error
	if lErr != nil {
		helperpkg.LogError(lErr)
		return ltotalRevenue, lErr
	}

	return ltotalRevenue, lErr
}

func GetRevenuebyProd(pInputRec model.GetDetails) ([]model.ProductRevenue, error) {

	log.Println(" GetRevenuebyProd (+) ")

	var ltotalRevenue_byprod []model.ProductRevenue

	lStartDate := common.GetDate(pInputRec.StartDate)
	lEndDate := common.GetDate(pInputRec.EndDate)

	lErr := dbconnection.G_Db_instance.Mysql_gormdb.Debug().Table("order_items").
		Select("products.name, nvl(SUM(order_items.quantity_sold  * order_items.unit_price * (1 - order_items.discount)),'') as total_revenue").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN orders ON order_items.order_id  = orders.id").
		Where("orders.date_of_sale BETWEEN ? AND ?", lStartDate, lEndDate).
		Group("order_items.product_id, products.name").
		Scan(&ltotalRevenue_byprod).Error

	if lErr != nil {
		helperpkg.LogError(lErr)
		return ltotalRevenue_byprod, lErr
	}

	log.Println(" GetRevenuebyProd (-) ")
	return ltotalRevenue_byprod, nil
}

func GetRevenuebyCat(pInputRec model.GetDetails) ([]model.CategoryRevenue, error) {

	log.Println("GetRevenuebyCat(+)")

	var ltotalRevenue_bycat []model.CategoryRevenue

	lStartDate := common.GetDate(pInputRec.StartDate)
	lEndDate := common.GetDate(pInputRec.EndDate)

	lErr := dbconnection.G_Db_instance.Mysql_gormdb.Debug().Table("order_items").
		Select("products.category as category, nvl(SUM(order_items.quantity_sold  * order_items.unit_price * (1 - order_items.discount)),'') as total_revenue").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.date_of_sale BETWEEN ? AND ?", lStartDate, lEndDate).
		Group("products.category").
		Scan(&ltotalRevenue_bycat).Error

	if lErr != nil {
		helperpkg.LogError(lErr)
		return ltotalRevenue_bycat, lErr
	}

	log.Println("GetRevenuebyCat(-)")
	return ltotalRevenue_bycat, lErr
}

func GetRevenuebyregion(pInputRec model.GetDetails) ([]model.RegionRevenue, error) {

	log.Println(" GetRevenuebyregion (+) ")

	var ltotalRevenueByregion []model.RegionRevenue

	lStartDate := common.GetDate(pInputRec.StartDate)
	lEndDate := common.GetDate(pInputRec.EndDate)

	lErr := dbconnection.G_Db_instance.Mysql_gormdb.Table("order_items").
		Select("nvl(orders.region,''), nvl(SUM(order_items.quantity_sold  * order_items.unit_price * (1 - order_items.discount)),'') as total_revenue").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.date_of_sale BETWEEN ? AND ?", lStartDate, lEndDate).
		Group("orders.region").
		Scan(&ltotalRevenueByregion).Error

	if lErr != nil {
		helperpkg.LogError(lErr)
		return ltotalRevenueByregion, lErr
	}

	log.Println(" GetRevenuebyregion (-) ")
	return ltotalRevenueByregion, lErr
}
