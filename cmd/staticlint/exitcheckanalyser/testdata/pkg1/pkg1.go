package main

import (
	"fmt"
	"os"
)

func main() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	fmt.Println("testing_text")
	os.Exit(100) // want "os.Exit declaration"
}
