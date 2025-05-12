package main

import (
	"fmt"
	"prog/examples"
)

func main() {
	fmt.Print("Выберите что вы хотите проверить:\n1. Создание карты\n\n")
	var number int 
	fmt.Scan(&number)
	switch number {
		case 1:
			examples.ExamplMaps()
	}
}