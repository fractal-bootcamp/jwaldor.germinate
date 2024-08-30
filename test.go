package another

import (
	"fmt"
	"os"
)

func another() {
	err := os.MkdirAll("myDir", 0744)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.Chdir("myDir")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile("myFile.txt", []byte("This is text content written in myDir"), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
