package main

import (
	"bufio"
	"log"
	"os"
	"net/http"
	"io/ioutil"
	"fmt"
)

func init() {
	log.SetFlags(0)
}

// Other funcs, if any.
func sendGetRequest(url string) {

	response, err := http.Get(url)
	//defer response.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}

	//reader := strings.NewReader("body")
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Sprint(body)
	_, err = os.Stdout.Write(body)
	if err != nil {
		log.Fatalln(err)
	}

}


func main() {

	// Initialization, if any.

	scanner := bufio.NewScanner(os.Stdin)


	for scanner.Scan() {
		// Do something with strings here.
		log.Printf("String: %s\n", scanner.Text())
		go sendGetRequest(scanner.Text())


	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	// Other code, if any.

	total := 0

	// Other code, if any.

	var input string
	fmt.Scanln(&input)
	log.Printf("Total: %v", total)
}
