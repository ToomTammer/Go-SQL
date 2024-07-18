package main

import "database/sql"

// func createProduct(product *Product) error {
//     _, err := db.Exec(`INSERT INTO products(name, price) VALUES($1, $2);`, product.Name, product.Price)
//     return err
// }

func createProduct(db *sql.DB, p *Product) (int, error) {
    var id int
    err := db.QueryRow(`INSERT INTO products(name, price, category, quantity) VALUES($1, $2, $3, $4) RETURNING id;`, p.Name, p.Price, p.Category, p.Quantity).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func createProductReturnModel(db *sql.DB, p *Product) error {
	_, err := db.Exec("INSERT INTO products (name, price, category) VALUES ($1, $2, $3)", p.Name, p.Price, p.Category)
  
	return err
}

func getProduct(db *sql.DB, id int) (Product, error) {
	var p Product
	row := db.QueryRow(`SELECT id, name, price, category, quantity FROM products WHERE id = $1;`, id)
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Quantity)
	if err != nil {
	  return Product{}, err
	}
	return p, nil
}

func getProducts() ([]Product, error) {
  rows, err := db.Query("SELECT id, name, price, category, quantity FROM products")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var products []Product
  for rows.Next() {
    var p Product
    err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Quantity)
    if err != nil {
      return nil, err
    }
    products = append(products, p)
  }

  // Check for errors from iterating over rows
  if err = rows.Err(); err != nil {
    return nil, err
  }

  return products, nil
}

func updateProduct(db *sql.DB, p *Product) error {
	// _ is used to ignore a value that you do not need.
	_, err := db.Exec(`UPDATE products SET name = $1, price = $2, category = $3, quantity = $4 WHERE id = $5;`, p.Name, p.Price, p.Category, p.Quantity, p.ID)
	return err
}

func deleteProduct(id int) error {
	_, err := db.Exec(`DELETE FROM products WHERE id = $1;`, id)
	return err
}

// Note

// Exec = เป็น method ที่ใช้สำหรับการ execute SQL โดยไม่มีการ return rows กลับคืนมา (เช่น INSERT, UPDATE, DELETE)
// QueryRow = เป็น method ที่ใช้สำหรับ query SQL เพื่อดึงข้อมูลกลับมา (เป็นข้อมูลตัวเดียว) ปกติจะใช้กับตระกูลของ SELECT

// เพิ่มผ่าน JOIN
func getProductsAndSuppliers(db *sql.DB) ([]ProductWithSupplier, error) {
  // SQL JOIN query
  query := `
      SELECT
          p.id AS product_id,
          p.name AS product_name,
          p.price,
          s.name AS supplier_name
      FROM
          products p
      INNER JOIN suppliers s
          ON p.supplier_id = s.id;`

  rows, err := db.Query(query)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var products []ProductWithSupplier
  for rows.Next() {
    var p ProductWithSupplier
    err := rows.Scan(&p.ProductID, &p.ProductName, &p.Price, &p.SupplierName)
    if err != nil {
      return nil, err
    }
    products = append(products, p)
  }

  if err = rows.Err(); err != nil {
    return nil, err
  }

  return products, nil
}