package dbconnection

import (
	"database/sql"
	"fmt"
	"log"
	helperpkg "sales_analytics/helper_pkg"
	"sales_analytics/tomlreader"
	"strconv"

	"gorm.io/gorm"
)

type DatabaseDetails struct {
	User     string
	Port     int
	Server   string
	Password string
	Database string
	DBType   string
	DB       string
}

type Db_Details struct {
	Mysql DatabaseDetails
}

type Db_instance struct {
	Mysql_sqldb  *sql.DB
	Mysql_gormdb *gorm.DB
}

// Variable holds the instances of the database connection.
var G_Db_instance Db_instance

// Struct for to hold the connection pool configuration
type connectionpoolconfig struct {
	OpenConnCount    int
	IdleConnCount    int
	MaxIdleConnCount int
}

/*
Method will read the database detail from the toml
Ex : Userdetail : "root", port: 3306 etc..,
*/
func Db_detailReading() Db_Details {
	log.Println("Db_detailReading (+) ")

	var lDbDetailsRec Db_Details

	lDbConfig := tomlreader.ReadTomlFile("./toml/dbconfig.toml")
	lDbDetailsRec.Mysql.User = fmt.Sprintf("%v", lDbConfig.(map[string]interface{})["Db_User"])
	lDbDetailsRec.Mysql.Port, _ = strconv.Atoi(fmt.Sprintf("%v", lDbConfig.(map[string]interface{})["Db_Port"]))
	lDbDetailsRec.Mysql.Server = fmt.Sprintf("%v", lDbConfig.(map[string]interface{})["Db_Server"])
	lDbDetailsRec.Mysql.Password = fmt.Sprintf("%v", lDbConfig.(map[string]interface{})["Db_Password"])
	lDbDetailsRec.Mysql.Database = fmt.Sprintf("%v", lDbConfig.(map[string]interface{})["Db_Database"])
	lDbDetailsRec.Mysql.DBType = fmt.Sprintf("%v", lDbConfig.(map[string]interface{})["Db_Type"])
	lDbDetailsRec.Mysql.DB = fmt.Sprintf("%v", lDbConfig.(map[string]interface{})["Db_Name"])

	log.Println("Db_detailReading (-) ")
	return lDbDetailsRec
}

func connectionpoolConfig() connectionpoolconfig {
	log.Println("connectionpoolConfig (+) ")

	var lConpoolconfig connectionpoolconfig
	var lErr error

	// reading a connection pool details from the toml
	lDbConnectionpool := tomlreader.ReadTomlFile("./toml/dbconfig.toml")
	lSetMaxOpenConns := fmt.Sprintf("%v", lDbConnectionpool.(map[string]interface{})["SetMaxOpenConnsdb"])
	lSetMaxIdleConnsdb := fmt.Sprintf("%v", lDbConnectionpool.(map[string]interface{})["SetMaxIdleConnsdb"])
	lSetConnMaxIdleTime := fmt.Sprintf("%v", lDbConnectionpool.(map[string]interface{})["SetConnMaxIdleTimedb"])

	// If the details not properly readen from the toml file this will handle the issue
	if lSetMaxOpenConns == "" {
		lSetMaxOpenConns = "10"
	}

	if lSetMaxIdleConnsdb == "" {
		lSetMaxIdleConnsdb = "5"
	}

	if lSetConnMaxIdleTime == "" {
		lSetConnMaxIdleTime = "5"
	}

	lConpoolconfig.OpenConnCount, lErr = strconv.Atoi(lSetMaxOpenConns)
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lConpoolconfig
	}

	lConpoolconfig.IdleConnCount, lErr = strconv.Atoi(lSetMaxIdleConnsdb)
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lConpoolconfig
	}

	lConpoolconfig.IdleConnCount, lErr = strconv.Atoi(lSetConnMaxIdleTime)
	if lErr != nil {
		helperpkg.LogError(lErr)
		return lConpoolconfig
	}

	log.Println("connectionpoolConfig (-) ")
	return lConpoolconfig
}
