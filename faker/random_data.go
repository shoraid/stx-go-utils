package faker

import (
	"math/rand"
	"time"

	"github.com/shoraid/stx-go-utils/genericutil"
)

func PickRandom(elements ...any) any {
	index := RandInt(0, len(elements)-1)
	return elements[index]
}

func RandBool() bool {
	return rand.Intn(2) == 1
}

func RandBoolPtr() *bool {
	return genericutil.Ptr(RandBool())
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func RandIntPtr(min, max int) *int {
	return genericutil.Ptr(RandInt(min, max))
}

func RandSentence(wordCount int) string {
	words := []string{
		"lorem", "ipsum", "dolor", "sit", "amet",
		"consectetur", "adipiscing", "elit",
		"sed", "do", "eiusmod", "tempor", "incididunt",
		"ut", "labore", "et", "dolore", "magna", "aliqua",
		"Ut", "enim", "ad", "minim", "veniam",
		"quis", "nostrud", "exercitation", "ullamco", "laboris",
		"nisi", "aliquip", "ex", "ea", "commodo", "consequat",
		"Duis", "aute", "irure", "in", "reprehenderit",
		"voluptate", "velit", "esse", "cillum", "eu",
		"fugiat", "nulla", "pariatur",
		"Excepteur", "sint", "occaecat", "cupidatat", "non", "proident",
		"sunt", "culpa", "qui", "officia", "deserunt",
		"mollit", "anim", "id", "est", "laborum",
	}
	result := ""
	for range wordCount {
		result += words[rand.Intn(len(words))] + " "
	}
	return result[:len(result)-1]
}

func RandSentencePtr(wordCount int) *string {
	return genericutil.Ptr(RandSentence(wordCount))
}

func RandString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seed := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(seed)

	result := make([]byte, length)

	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result)
}

func RandStringPtr(length int) *string {
	return genericutil.Ptr(RandString(length))
}

func RandTime(start, end time.Time) time.Time {
	if start.After(end) {
		start, end = end, start
	}

	duration := rand.Int63n(end.Unix() - start.Unix())
	return time.Unix(start.Unix()+duration, 0)
}

func RandTimePtr(start, end time.Time) *time.Time {
	t := RandTime(start, end)
	return &t
}

func RandURL() string {
	domains := []string{"example.com", "test.com", "dummy.net", "sample.org"}
	return "https://" + RandString(8) + "." + domains[rand.Intn(len(domains))]
}

func RandURLPtr() *string {
	return genericutil.Ptr(RandURL())
}
