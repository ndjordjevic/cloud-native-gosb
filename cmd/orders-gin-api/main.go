package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Order struct {
	ID         int     `json:"id"`
	Account    string  `json:"account"`
	Instrument string  `json:"instrument"`
	Volume     int     `json:"volume"`
	Price      float32 `json:"price"`
}

var dataStore = map[int]Order{
	1: {1, "acc", "BMW", 10, 5},
	2: {2, "acc", "Apple", 12, 6},
	3: {3, "acc", "Google", 7, 8},
}

func main() {
	r := gin.Default()

	v1 := r.Group("/v1/orders")
	{
		v1.GET("/", getAllOrders)
		v1.GET("/:id", getOrderById)
		v1.POST("/", createOrder)
		v1.DELETE("/:id", deleteOrder)
	}

	_ = r.Run()
}

func deleteOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}

	delete(dataStore, id)

	c.String(http.StatusOK, "Order is successfully deleted")
}

func createOrder(c *gin.Context) {
	var order Order

	if err := c.BindJSON(&order); err == nil {
		dataStore[order.ID] = order
		c.String(http.StatusCreated, "Order is successfully created")
	} else {
		c.String(http.StatusInternalServerError, "Order couldn't be created")
		log.Fatal(err)
	}
}

func getOrderById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	order, ok := dataStore[id]
	if ok {
		c.JSON(http.StatusOK, order)
	} else {
		c.String(http.StatusNotFound, "Couldn't be found!")
	}
}

func getAllOrders(c *gin.Context) {
	values := make([]Order, 0, len(dataStore))

	for _, v := range dataStore {
		values = append(values, v)
	}
	c.JSON(http.StatusOK, values)
}
