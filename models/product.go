package models

type Product struct {
	ID    		  int     `json:"id"`
	Name  		  string  `json:"name"`
	Price 		  int     `json:"price"`
	Stock 		  int     `json:"stock"`
	Category_id   int     `json:"category_id"`
	Category_name string  `json:"category_name"`
}

type Productsy struct {
	ID    		  int     `json:"id"`
	Name  		  string  `json:"name"`
	Price 		  int     `json:"price"`
	Stock 		  int     `json:"stock"`
	Category_id   int     `json:"-"`
	Category_name string  `json:"category_name"`
}