package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

var apiURL string = "https://api.openai.com/v1/audio/transcriptions"

func main() {

	fmt.Println(os.Args)

	if len(os.Args) != 3 {
		fmt.Println("Please enter audio filename and key")
		fmt.Println("Example:  go run . \"sample-0.mp3\" \"open-ai-key\"")
		os.Exit(-1)
	}

	filename := os.Args[1]
	key := os.Args[2]

	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", file.Name())
	io.Copy(part, file)

	writer.WriteField("model", "whisper-1")

	writer.Close()

	r, _ := http.NewRequest("POST", apiURL, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.Header.Add("Authorization", "Bearer "+key)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return
	}

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", string(responseBody))

}
