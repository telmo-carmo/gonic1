package dal

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type Bonus struct {
	Ename string `json:"ename"`
	Job   string `json:"job"`
	Sal   int    `json:"sal"`
	Comm  int    `json:"comm"`
}

type Dept struct {
	Deptno    int    `json:"deptno"`
	Dname     string `json:"dname"`
	Loc       string `json:"loc"`
	VersionNo int    `json:"versionno"`
	Active    int    `json:"active"`
}

var DB *gorm.DB

// InitDatabase does db inital stuff
func InitDatabase() {
	dbname := "gonic1.db"
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database " + dbname)
	}
	log.Printf("[INFO] using db: %s", dbname)
	db.AutoMigrate(&Book{})

	/*
		// Migrate the schema
		db.AutoMigrate(&Product{})

		// Create
		db.Create(&Product{Code: "D42", Price: 100})
		//db.Create(&Book{Title: "Dune", Author: "John Adams"})

		// Read
		var product Product
		db.First(&product, 1)                 // find product with integer primary key
		db.First(&product, "code = ?", "D42") // find product with code D42

		// Update - update product's price to 200
		db.Model(&product).Update("Price", 200)
		log.Println(product)

		// Update - update multiple fields
		db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
		db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

		// Delete - delete product
		db.Delete(&product, 1)
	*/
	DB = db
}
