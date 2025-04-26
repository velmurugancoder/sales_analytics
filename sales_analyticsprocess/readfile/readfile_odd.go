package readfile

/*
import (
	"encoding/csv"
	"log"
	"os"
	"sales_analytics/dbconnection"
	helperpkg "sales_analytics/helper_pkg"
	"sales_analytics/sales_analyticsprocess/model"
	"strconv"
)

func CsvFile_Reader(pFilepath string) error {

	//Opening the file
	lFile, lErr := os.Open(pFilepath)
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lErr
	}

	//Closing the file
	defer lFile.Close()

	// Create a new CSV reader
	lReader := csv.NewReader(lFile)

	//Reading all the rows in the file
	lRecords, lErr := lReader.ReadAll()
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lErr
	}

	lCustomerMap := make(map[string]bool)
	var lCustomerArr []model.Customerdetails

	lProductMap := make(map[string]bool)
	var lProductArr []model.Products

	var lOrdersArr []model.Orders

	var lOderItemsArr []model.Order_items

	for lIdx, lRows := range lRecords {
		log.Println("lRows", lRows)

		if lIdx == 0 {
			continue
		}

		var lCustomerRec model.Customerdetails
		var lProductsRec model.Products
		var lOrdersRec model.Orders
		var lOrder_itemsRec model.Order_items

		lCustomerRec.Customer_id = lRows[2]
		lCustomerRec.Customer_name = lRows[12]
		lCustomerRec.Customer_email = lRows[13]
		lCustomerRec.Customer_address = lRows[14]

		if !lCustomerMap[lCustomerRec.Customer_id] {
			lCustomerMap[lCustomerRec.Customer_id] = true
			lCustomerArr = append(lCustomerArr, lCustomerRec)
		}

		lProductsRec.Product_id = lRows[1]
		lProductsRec.Name = lRows[3]
		lProductsRec.Category = lRows[4]

		if !lProductMap[lProductsRec.Product_id] {
			lProductMap[lProductsRec.Product_id] = true
			lProductArr = append(lProductArr, lProductsRec)
		}

		lOrdersRec.Order_id = lRows[0]
		lOrdersRec.Customer_id = lRows[2]
		lOrdersRec.Region = lRows[5]
		lOrdersRec.Date_of_sale = lRows[6]
		lOrdersRec.Payment_method = lRows[11]
		lOrdersRec.Shipping_cost = lRows[10]
		lOrdersArr = append(lOrdersArr, lOrdersRec)

		lOrder_itemsRec.Order_id = lRows[0]
		lOrder_itemsRec.Product_id = lRows[1]
		lOrder_itemsRec.Quantity_sold, lErr = strconv.Atoi(lRows[7])
		if lErr != nil {
			helperpkg.LogError(lErr)
			return lErr
		}
		lOrder_itemsRec.Unit_price, lErr = strconv.ParseFloat(lRows[8], 64)
		if lErr != nil {
			helperpkg.LogError(lErr)
			return lErr
		}
		lOrder_itemsRec.Discount, lErr = strconv.ParseFloat(lRows[9], 64)
		if lErr != nil {
			helperpkg.LogError(lErr)
			return lErr
		}

		lOderItemsArr = append(lOderItemsArr, lOrder_itemsRec)

	}

	lErr = dbconnection.G_Db_instance.Mysql_gormdb.Table("customers").Create(&lCustomerArr).Error
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lErr
	}

	lErr = dbconnection.G_Db_instance.Mysql_gormdb.Table("products").Create(&lProductArr).Error
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lErr
	}

	lErr = dbconnection.G_Db_instance.Mysql_gormdb.Table("orders").Create(&lOrdersArr).Error
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lErr
	}

	lErr = dbconnection.G_Db_instance.Mysql_gormdb.Table("order_items").Create(&lOderItemsArr).Error
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lErr
	}

	return nil
}

func insertCustomerdetails(pCustomerRec model.Customerdetails) error {
	log.Println("insertCustomerdetails (+) ")

	lErr := dbconnection.G_Db_instance.Mysql_gormdb.Table("").Create(&pCustomerRec).Error
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lErr
	}

	log.Println("insertCustomerdetails (-) ")
}

*/
