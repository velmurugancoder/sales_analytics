package readfile

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

	for lIdx, lRows := range lRecords {

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

		lCustomerId, lErr := CheckAlreadyPresent("customers", `customer_id = '`+lCustomerRec.Customer_id+`'`)
		if lErr != nil {
			helperpkg.LogError(lErr)
			return lErr
		}

		if lCustomerId == 0 {
			lErr = dbconnection.G_Db_instance.Mysql_gormdb.Table("customers").Create(&lCustomerRec).Error
			if lErr != nil {
				helperpkg.LogError(lErr)
				return lErr
			}

			lCustomerId = lCustomerRec.ID

		}

		lProductsRec.Product_id = lRows[1]
		lProductsRec.Name = lRows[3]
		lProductsRec.Category = lRows[4]

		lproductId, lErr := CheckAlreadyPresent("products", `product_id = '`+lProductsRec.Product_id+`'`)
		if lErr != nil {
			helperpkg.LogError(lErr)
			return lErr
		}

		if lproductId == 0 {
			lErr = dbconnection.G_Db_instance.Mysql_gormdb.Table("products").Create(&lProductsRec).Error
			if lErr != nil {
				helperpkg.LogError(lErr)
				return lErr
			}

			lproductId = lProductsRec.ID

		}

		lOrdersRec.Order_id = lRows[0]
		lOrdersRec.Customer_id = lCustomerId
		lOrdersRec.Region = lRows[5]
		lOrdersRec.Date_of_sale = lRows[6]
		lOrdersRec.Payment_method = lRows[11]
		lOrdersRec.Shipping_cost = lRows[10]

		lOrderId, lErr := CheckAlreadyPresent("orders", `order_id = '`+lOrdersRec.Order_id+`'`)
		if lErr != nil {
			helperpkg.LogError(lErr)
			return lErr
		}

		if lOrderId == 0 {
			lErr = dbconnection.G_Db_instance.Mysql_gormdb.Table("orders").Create(&lOrdersRec).Error
			if lErr != nil {
				helperpkg.LogError(lErr)
				return lErr
			}
		} else {
			continue
		}

		lOrder_itemsRec.Order_id = lOrdersRec.ID
		lOrder_itemsRec.Product_id = lproductId
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

		lErr = dbconnection.G_Db_instance.Mysql_gormdb.Table("order_items").Create(&lOrder_itemsRec).Error
		if lErr != nil {
			helperpkg.LogError(lErr)
			return lErr
		}
	}
	return nil
}

func CheckAlreadyPresent(pTableName, pWherecon string) (uint, error) {
	log.Println(" CheckAlreadyPresent (+) ")

	var lId uint

	lErr := dbconnection.G_Db_instance.Mysql_gormdb.Table(pTableName).Select("id").Where(pWherecon).Scan(&lId).Error
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lId, lErr
	}

	log.Println(" CheckAlreadyPresent (-) ")
	return lId, lErr
}
