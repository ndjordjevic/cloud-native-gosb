package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Order struct {
	ID         int `json:"id"`
	Account    string
	Instrument string
	Volume     int
	Price      float32
}

var dataStore = map[int]Order{
	1: {1, "acc", "BMW", 10, 5},
	2: {2, "acc", "Apple", 12, 6},
}

func main() {
	r := gin.Default()

	v1 := r.Group("/v1/order")
	{
		v1.GET("/", getAllOrders)
		v1.GET("/:id", getOrderById)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}

func getOrderById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	order, ok := dataStore[id]
	if ok {
		c.JSON(http.StatusOK, order)
	} else {
		c.JSON(http.StatusNotFound, nil)
	}
}

func getAllOrders(c *gin.Context) {
	values := make([]Order, 0, len(dataStore))

	for _, v := range dataStore {
		values = append(values, v)
	}
	c.JSON(http.StatusOK, values)
}
