package main

import (
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files"       // swagger embed files

import _ "github.com/ndjordjevic/cloud-native-gosb/cmd/orders-gin-api/docs"

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

// @title Orders Gin API
// @version 1.0.0
// @description This is a sample server for gin order CRUD ops
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationUrl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationUrl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information
func main() {
	gin.ForceConsoleColor()
	r := gin.Default()

	r.Use(cors.Default())

	v1 := r.Group("/v1/orders")
	{
		v1.GET("/", getAllOrders)
		v1.GET("/:id", getOrderById)
		v1.POST("/", createOrder)
		v1.PUT("/:id", updateOrder)
		v1.DELETE("/:id", deleteOrder)
		v1.PATCH("/:id", patchOrder)
	}

	securityConfig := secure.DefaultConfig()
	securityConfig.SSLRedirect = false
	//r.Use(secure.New(securityConfig))

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	//_ = r.RunTLS(":8080", "cmd/orders-gin-api/domain.crt", "cmd/orders-gin-api/domain.key")
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

// deleteOrder godoc
// @Summary Delete order
// @Description Delete by order ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param  id path int true "Order ID" Format(int64)
// @Success 204 "Order is successfully deleted"
// @Router /orders/{id} [delete]
func deleteOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}

	delete(dataStore, id)

	c.String(http.StatusNoContent, "Order is successfully deleted")
}

// createOrder godoc
// @Summary Create new order
// @Description Create new order
// @Tags orders
// @Accept  json
// @Produce  json
// @Param account body Order true "New Order"
// @Success 200 "Order is successfully created"
// @Router /orders [post]
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

// getOrderById godoc
// @Summary Get order by id
// @Description Get order by id
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} Order
// @Router /orders/{id} [get]
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

// getAllOrders godoc
// @Summary Get all orders
// @Description Returns all orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} Order
// @Router /orders [get]
func getAllOrders(c *gin.Context) {
	values := make([]Order, 0, len(dataStore))

	for _, v := range dataStore {
		values = append(values, v)
	}
	c.JSON(http.StatusOK, values)
}
