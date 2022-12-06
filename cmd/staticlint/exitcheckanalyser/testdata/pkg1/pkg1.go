package main

import (
	"fmt"
	"os"
)

func main() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	fmt.Println("testing_text")
	os.Exit(100) // wa_nt "os.Exit declaration"
}
