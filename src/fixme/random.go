package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

func generateText(min, max int) string {
	resp, err := http.Get(fmt.Sprintf("http://www.randomtext.me/api/lorem/p-1/%d-%d", min, max))
	if err != nil {
		return "Err:" + err.Error()
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	var t struct {
		TextOut string `json:"text_out"`
	}
	d.Decode(&t)
	return t.TextOut[4 : len(t.TextOut)-5]
}

func GenRandomProblem(id int) Problem {
	e := []string{EffectNo, EffectYes, EffectUnknown}
	return Problem{
		Id:          fmt.Sprintf("%d", id),
		Title:       generateText(5, 7),
		Description: generateText(10, 20),
		Effected:    e[rand.Int31n(3)],
	}
}
