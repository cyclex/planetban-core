package pkg

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func TokenGenerator(num int) string {
	b := make([]byte, num)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func IsInt(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}

	return false
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ContainsInt(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ExtractUsernames(text string) []string {
	// Define the regular expression pattern
	pattern := `@([a-zA-Z0-9_]+)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find all matches in the text
	matches := re.FindAllString(text, -1)

	return matches
}

func ExtractSentencesAfterWord(text, keyword string) []string {
	// Split the text into sentences
	sentences := strings.Split(text, ".")

	var result []string

	// Iterate through each sentence
	for _, sentence := range sentences {
		// Find the position of the keyword in the sentence
		index := strings.Index(sentence, keyword)
		if index != -1 {
			// If the keyword is found, extract the part of the sentence after the keyword
			result = append(result, strings.TrimSpace(sentence[index+len(keyword):]))
		}
	}

	return result
}

func IsTimeInBetween(startEpoch, endEpoch int64) bool {
	// Convert epoch values to time.Time
	startTime := time.Unix(startEpoch, 0)
	endTime := time.Unix(endEpoch, 0)

	// Get the current time
	currentTime := time.Now()

	return currentTime.After(startTime) && currentTime.Before(endTime)
}

func IsLetter(s string) bool {
	var regex, _ = regexp.Compile(`^[a-zA-z\s]+$`)
	return regex.MatchString(s)
}

func ShortUUID(uuid string) string {
	// Hash the UUID using MD5
	hash := md5.Sum([]byte(uuid))

	// Convert the hash to base64 encoding
	base64Str := base64.StdEncoding.EncodeToString(hash[:])

	// Take the first 6 characters from the base64 string
	return base64Str[:6]
}
