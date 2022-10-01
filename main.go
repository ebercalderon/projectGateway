package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	request "github.com/ebercalderon/projectGateway/request_redirect"
	"github.com/ebercalderon/projectGateway/types"
)

func getAPI(c *gin.Context) {
	msg := "Saludos! Bienvenidos al APIGateway"
	successful := true

	APIData := types.APIResponse{Message: &msg, Successful: &successful, Data: nil}
	out, err := json.Marshal(APIData)
	if err != nil {
		panic(err)
	}
	data := string(out)
	APIRes := types.APIData{
		Data: &data,
	}

	c.JSON(http.StatusOK, APIRes)
}

func getSummary(c *gin.Context) {
	// fechas, err := strconv.ParseInt(c.Param("fecha"), 10, 64)
	// if err != nil {
	// 	msg := "Formato de fecha erroneo. Debe ser Unix Epoch en milisegundos."
	// 	successful := false
	// 	c.JSON(http.StatusBadRequest, types.APIResponse{
	// 		Message:    &msg,
	// 		Successful: &successful,
	// 		Data:       nil,
	// 	})
	// 	return
	// }
	c.JSON(http.StatusOK, request.RequestGetAnalysis(c.Param("fecha"), os.Getenv("ERPANALYSIS_URL")+"api/analytics/summary"))
}

func postRegistro(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	APIRes := request.RedirectRequest(body, os.Getenv("ERPREGISTRATION_URL")+"api/empleados", "POST")
	c.JSON(http.StatusOK, APIRes)
}

func postGraphQL(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	APIRes := request.RedirectRequest(body, os.Getenv("ERPBACK_URL")+"graphql", "POST")
	c.JSON(http.StatusOK, APIRes)
}

func postRegistroConfirmacion(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	url := os.Getenv("ERPREGISTRATION_URL") + "api/confirmacion/" + c.Param("token")
	jsonData := fmt.Sprintf(`{"password": "%s"}`, strings.Split(string(body), "=")[1])
	request, reqErr := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if reqErr != nil {
		fmt.Printf("The HTTP request failed with error %s\n", reqErr)
		panic(reqErr)
	}
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	defer response.Body.Close()

	resContent, _ := ioutil.ReadAll(response.Body)
	c.Data(http.StatusOK, "text/html; charset=utf-8", resContent)
}

func getRegistroConfirmacion(c *gin.Context) {
	token := c.Param("token")
	url := os.Getenv("ERPREGISTRATION_URL") + "api/confirmacion/" + token
	request, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		fmt.Printf("The HTTP request failed with error %s\n", reqErr)
		panic(reqErr)

	}
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	defer response.Body.Close()

	resContent, _ := ioutil.ReadAll(response.Body)
	c.Data(http.StatusOK, "text/html; charset=utf-8", resContent)
}

func postRecommendation(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	url := os.Getenv("ERPRECOMMENDER_URL")
	request, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if reqErr != nil {
		fmt.Printf("The HTTP request failed with error %s\n", reqErr)
		panic(reqErr)
	}
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	defer response.Body.Close()

	resContent, _ := ioutil.ReadAll(response.Body)
	c.JSON(http.StatusOK, string(resContent))
}

func main() {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatalf("Some error occured. Err: %s", errEnv)
	}

	router := gin.Default()
	router.GET("/", getAPI)
	router.GET("/api", getAPI)
	router.GET("/api/analytics/summary/:fecha", getSummary)
	router.POST("/api/registro", postRegistro)
	router.GET("/api/registro/confirmacion/:token", getRegistroConfirmacion)
	router.POST("/api/registro/confirmacion/:token", postRegistroConfirmacion)
	router.POST("/api/graphql", postGraphQL)
	router.POST("/api/recommender", postRecommendation)

	router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("GATEWAY_PORT")))

	log.Println("¡API Gateway iniciado!")
}
