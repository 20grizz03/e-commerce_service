package hendlers

import "e-commerceProject/app/internal/models"

// интерфейс для описания логики действий продукта
type Catalog interface {
	GetAllProducts() ([]models.Item, error)
	GetProductByName(name string) (models.Item, error)
}

func (dc *DatabaseCatalog) GetAllProducts() ([]models.Item, error) {
	rows, err := dc.db.Query("SELECT id, name, price, quantity, category, info FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Item{}
	for rows.Next() {
		var product Product
		if err := rows.Scan(&models.Item.Id, &models.Item.Name, &product.Price, &product.Quantity, &product.Category, &product.Info); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
