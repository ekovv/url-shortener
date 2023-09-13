package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		body, _ := c.GetRawData() // получаем тело запроса в виде []byte
		str := string(body)       // преобразовываем []byte в string
		fmt.Println(str)          // выводим полученную строку в консоль

		c.Status(200)
	})

	r.Run(":8080")
}
