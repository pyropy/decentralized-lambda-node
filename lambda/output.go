package lambda

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

func writeOutput(result io.Reader) error {
	b, err := ioutil.ReadAll(result)
	if err != nil {
		return err
	}

	err = os.WriteFile("/outputs/output.json", b, 0644)
	if err != nil {
		log.Fatal("Error writing response body")
		return err
	}

	return nil
}
