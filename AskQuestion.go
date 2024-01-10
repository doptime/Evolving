package main

import (
	"encoding/json"
	"fmt"
)

func askQuestion(js string) (prompt string, err error) {
	var (
		myQuestion *MyQuestion
	)
	if err = json.Unmarshal([]byte(js), &myQuestion); err != nil {
		return "", err
	}
	//print the question to the console, and wait for the answer
	fmt.Println(" the product manager has a question: " + myQuestion.MyQuestion)
	//read the answer from the console
	fmt.Scanln(&prompt)
	return prompt, nil
}
