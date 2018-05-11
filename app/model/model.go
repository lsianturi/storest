package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type CommonModel struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

//Product is
type Product struct {
	CommonModel
	Name      string `json:"name"`
	Detail    string `json:"detail"`
	Price     int64  `json:"price"`
	Uom       string `json:"uom"`
	Quality   string `json:"quality"`
	Inventory int    `json:"inventory"`
	ImagePath string `json:"image_path"`
}

// Customer structure
type Customer struct {
	CommonModel
	Email        string `gorm:"unique_index" json:"email"`
	Password     string `gorm:"-" json:"-"`
	Token        string `gorm:"-" json:"token"`
	PasswordHash string `gorm:"type:varchar(255); not null" json:"password_hash"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Role         string `gorm:"type:ENUM('CUSTOMER', 'ADMIN');default:'CUSTOMER'" json:"role"`
	Addresses    []Address `gorm:"ForeignKey:CustomerID" json:"addresses,omitempty"`
	Carts        []Cart `gorm:"ForeignKey:CustomerID" json:"carts,omitempty"`
}

// Address of customer
type Address struct {
	ID         int
	Address    string `json:"address"`
	CustomerID uint `json:"customer_id"`
}

//Cart structure
type Cart struct {
	CommonModel
	PaymentType     string      `gorm:"type:ENUM('COD', 'BT');default:'COD'" json:"payment_type"`
	PaymentStatus   bool        `json:"payment_status"`
	DeliveryStatus  bool        `json:"delivery_status"`
	Items           []CartItem  `gorm:"ForeignKey:CartID" json:"items,omitempty"`
	CustomerID      uint        `json:"customer_id"`
	DeliveryAddr    Address     `json:"address,omitempty"`
	DeliveryAddrID  int         `json:"delivery_address_id"`
}
//Paid is
func (c *Cart) Paid() {
	c.PaymentStatus = true
}
//Unpaid is
func (c *Cart) Unpaid() {
	c.PaymentStatus = false
}
//Delivered is
func (c *Cart) Delivered() {
	c.DeliveryStatus = true
}
//Undelivered is
func (c *Cart) Undelivered() {
	c.DeliveryStatus = false
}
// CartItem structure
type CartItem struct {
	CommonModel
	Quantity  int      `json:"qty"`
	Price     float64  `json:"price"`
	CartID    uint     `json:"cart_id"`
	Product   Product  `json:"product,omitempty"`
	ProductID uint     `json:"product_id"`
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Product{}, &Customer{}, &Address{}, &Cart{}, CartItem{})
	db.Model(&Address{}).AddForeignKey("customer_id", "customers(id)", "CASCADE", "CASCADE")
	db.Model(&Cart{}).AddForeignKey("customer_id", "customers(id)", "CASCADE", "CASCADE")
	db.Model(&CartItem{}).AddForeignKey("cart_id", "carts(id)", "CASCADE", "CASCADE")
	return db
}
