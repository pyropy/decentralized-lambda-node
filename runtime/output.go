package runtime

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (r *LambdaRuntime) writeOutput(result *http.Response, meta *ResultMetadata) error {
	headersJson, err := json.MarshalIndent(meta, "", "")
	if err != nil {
		log.Fatal("Error serializing metadata to json")
	}

	err = ioutil.WriteFile("/outputs/metadata.json", headersJson, 0644)
	if err != nil {
		log.Fatal("Error writing metadata json to file")
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatal("Error reading response body")
	}

	err = ioutil.WriteFile("/outputs/body", body, 0644)
	if err != nil {
		log.Fatal("Error writing response body")
	}

	return nil
}
