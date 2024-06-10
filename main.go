package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	num := 1
	if len(os.Args) >= 2 {
		argNum, err := strconv.Atoi(os.Args[1])
		if err == nil && 1 <= argNum {
			num = argNum
		}
	}

	resp, err := http.Get("https://raw.githubusercontent.com/hermitdave/FrequencyWords/master/content/2018/ja/ja_full.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Status code", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	lines := strings.Split(string(body), "\n")

	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < num; i++ {
		randomIndex := rand.Intn(len(lines))
		randomIndex2 := rand.Intn(len(lines))
		randomWord := strings.Split(lines[randomIndex], " ")[0]
		randomWord2 := strings.Split(lines[randomIndex2], " ")[0]
		fmt.Println(randomWord + randomWord2)
	}
}
