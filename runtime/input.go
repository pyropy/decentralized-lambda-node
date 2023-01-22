package runtime

import (
	"bufio"
	"log"
	"net/http"
	"os"
)

func (r *LambdaRuntime) parseInput() (*http.Request, error) {
	file, err := os.Open("/inputs.json")
	if err != nil {
		log.Println("Error opening input json")
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	req, err := http.NewRequest("POST", "/", reader)
	if err != nil {
		log.Println("Error creating request")
		return nil, err
	}

	return req, nil
}
