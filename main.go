package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	//// _ การ import เพื่อให้ lib ตัวอื่นสามารใช้ได้ เช่น เพื่อให้ "database/sql" ใช้
	_ "github.com/lib/pq"
)

const (
  host     = "localhost"  // or the Docker service name if running in another container
  port     = 5432         // default PostgreSQL port
  user     = "myuser"     // as defined in docker-compose.yml
  password = "mypassword" // as defined in docker-compose.yml
  dbname   = "mydatabase" // as defined in docker-compose.yml
)

var db *sql.DB

type Product struct {
	ID       int `json:"id"`
	Name     string `json:"name"`
	Price    int `json:"price"`
	Category string `json:"category"`
	Quantity int `json:"quantity"`
}

type ProductWithSupplier struct {
  ProductID        int `json:"productId"`
  ProductName      string `json:"productName"`
  Price            int `json:"price"`
  SupplierName     string `json:"supplierName"`
}

func SetupDatabase() *sql.DB {
  // Connection string
  connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
                          "password=%s dbname=%s sslmode=disable",
                          host, port, user, password, dbname)
  // Open a connection
  db, err := sql.Open("postgres", connectionString)
  if err != nil {
    log.Fatal(err)
  }

  // Check the connection
  if err = db.Ping(); err != nil {
    log.Fatal(err)
  }

  return db
}

func main() {
  app := fiber.New()
  db = SetupDatabase()
  // defer คือ คำสั่งที่จะทำงานก่อนคำสั่งสุดท้ายใน program (หรือ function นั้นๆ) โดยปกติมันจะใช้สำหรับคำสั่ง cleanup เพื่อปิด process ให้ครบก่อนที่จะหยุดทำงาน, ป้องกันการทำซ้ำ
  defer db.Close() 

  app.Get("/products", GetProducts) 
  app.Get("/product/:id", GetProduct) 
  app.Post("/product", CreateProduct) 
  app.Put("/product/:id", UpdateProduct) 
  app.Delete("/product/:id", DeleteProduct)

  app.Listen(":8080")

  // fmt.Println("Successfully connected!")

  // // var sid int
  // // sid, err = createProduct(&Product{Name: "Go product2", Price: 444, Category: "1", Quantity: 2})
  // // sid, err = createProduct(db, "Go product4", 555, "aloha", 2)
  // // pro, err := getProducts(db)
  // pro, err := getProductsAndSuppliers(db)
  // if err != nil {
	// 	log.Fatal(err)
	// }
	// // fmt.Println("Create connected!", sid)
	// fmt.Println("getProductsndSuppliers connected!", pro)

}

func GetProducts(c *fiber.Ctx) error {
  // Retrieve product from database
  products, err := getProducts()

  if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

  return c.JSON(products)
}

func GetProduct(c *fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

  // Retrieve product from database
  product, err := getProduct(db ,id)

  if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

  return c.JSON(product)
}

func CreateProduct(c *fiber.Ctx) error {
  p := new(Product)
  if err := c.BodyParser(p); err != nil {
    return err
  }
  
  // Insert product into database
  err := createProductReturnModel(db, p)
  if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

  return c.JSON(p)
}

func UpdateProduct(c *fiber.Ctx) error {
  id := c.Params("id")
  p := new(Product)
  if err := c.BodyParser(p); err != nil {
    return err
  }
  p.ID, _ = strconv.Atoi(id)

  // Update product in the database
  err := updateProduct(db, p)
  if err != nil {
    return err
  }

  return c.JSON(p)
}

func DeleteProduct(c *fiber.Ctx) error {
  id, err  := strconv.Atoi(c.Params("id"))
  if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

  // Delete product from database
  err = deleteProduct(id)
  
  if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

  return c.SendStatus(fiber.StatusNoContent)
}


