package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
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
		v1.PUT("/:id", updateOrder)
		v1.DELETE("/:id", deleteOrder)
		v1.PATCH("/:id", patchOrder)
	}

	_ = r.Run()
}

func patchOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}

	data, err := c.GetRawData()
	if err != nil {
		panic(err)
	}

	mapBody := make(map[string]string)
	if err = json.Unmarshal(data, &mapBody); err != nil {
		panic(err)
	}

	orderToPatch := dataStore[id]
	for k, v := range mapBody {
		switch k {
		case "volume":
			volume, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			reflect.ValueOf(&orderToPatch).Elem().FieldByName("Volume").SetInt(int64(volume))
			dataStore[id] = orderToPatch
		}
	}
}

func updateOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}

	var order Order
	if err := c.BindJSON(&order); err == nil {
		dataStore[id] = order
		c.String(http.StatusOK, "Order is successfully updated")
	} else {
		c.String(http.StatusInternalServerError, "Order couldn't be updated")
		log.Fatal(err)
	}
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
