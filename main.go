package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func calculateAromaAmount(liquidAmount, aromaConcentration float64) float64 {
	return aromaConcentration / 100 * liquidAmount
}

func calculateNicAmount(desNic, actNic, liquidAmount float64) float64 {
	shotNic := liquidAmount * desNic
	if actNic == 0.0 {
		return 0
	}
	return shotNic / actNic
}

func calculateVGPG(liquidAmount, vg float64) (float64, float64) {
	return vg / 100 * liquidAmount, (100 - vg) / 100 * liquidAmount
}

func subtLiquid(liquidAmount *float64, aromaAmount, nicAmount float64) {
	*liquidAmount = *liquidAmount - aromaAmount - nicAmount
}

func resultPostHandler(c *gin.Context) {
	var liquidAmount, desNic, actNic, armConc, vgPer float64
	var err error
	if liquidAmount, err = strconv.ParseFloat(c.PostForm("liquidAmount"), 64); err != nil {
		liquidAmount = 100
	}
	if desNic, err = strconv.ParseFloat(c.PostForm("desNic"), 64); err != nil {
		desNic = 0
	}
	if actNic, err = strconv.ParseFloat(c.PostForm("actNic"), 64); err != nil {
		actNic = 0
	}
	if armConc, err = strconv.ParseFloat(c.PostForm("armConc"), 64); err != nil {
		armConc = 10
	}
	if vgPer, err = strconv.ParseFloat(c.PostForm("vgPer"), 64); err != nil {
		vgPer = 50
	}
	aromaAmount := calculateAromaAmount(liquidAmount, armConc)
	nicAmount := calculateNicAmount(desNic, actNic, liquidAmount)
	subtLiquid(&liquidAmount, aromaAmount, nicAmount)
	vgAmount, pgAmount := calculateVGPG(liquidAmount, vgPer)
	c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
		"nicotineAmount": nicAmount,
		"aromaAmount":    aromaAmount,
		"pg":             pgAmount,
		"vg":             vgAmount,
	})
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/cover", func(c *gin.Context) {
		var nicotineAmount, aromaAmount, pg, vg float64
		c.HTML(http.StatusOK, "cover.tmpl.html", gin.H{
			"nicotineAmount": nicotineAmount,
			"aromaAmount":    aromaAmount,
			"pg":             pg,
			"vg":             vg,
		})
	})

	router.POST("/result", resultPostHandler)

	router.Run(":" + port)
}
