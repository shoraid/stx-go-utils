package faker

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/shoraid/stx-go-utils/genericutil"
)

// PickRandom returns a random element from the provided list.
func PickRandom(elements ...any) any {
	index := RandInt(0, len(elements)-1)
	return elements[index]
}

// RandBool returns a random boolean value (true or false).
func RandBool() bool {
	return rand.Intn(2) == 1
}

// RandBoolPtr returns a pointer to a random boolean value (true or false).
func RandBoolPtr() *bool {
	return genericutil.Ptr(RandBool())
}

// RandEmail generates a random email address in the format "username123@domain".
func RandEmail() string {
	usernames := []string{
		"john", "jane", "alex", "mike", "sara",
		"emma", "lisa", "david", "kevin", "nina",
		"peter", "sophia", "mark", "olivia", "jack",
		"lucas", "mia", "ryan", "chloe", "daniel",
		"zoe", "adam", "ella", "sam", "grace",
		"noah", "ava", "liam", "isabella", "ethan",
	}

	domains := []string{
		"gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "icloud.com",
		"example.com", "test.com", "dummy.net", "sample.org", "mail.com",
	}

	username := usernames[rand.Intn(len(usernames))]

	usernameWithDigits := username + fmt.Sprintf("%03d", RandInt(0, 999))

	domain := domains[rand.Intn(len(domains))]

	return usernameWithDigits + "@" + domain
}

// RandEmailPtr returns a pointer to a randomly generated email address.
func RandEmailPtr() *string {
	return genericutil.Ptr(RandEmail())
}

// RandInt returns a random integer within the range [min, max].
func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandIntPtr returns a pointer to a random integer within the range [min, max].
func RandIntPtr(min, max int) *int {
	return genericutil.Ptr(RandInt(min, max))
}

// RandSentence generates a random sentence consisting of `wordCount` words.
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

// RandSentencePtr returns a pointer to a random sentence consisting of `wordCount` words.
func RandSentencePtr(wordCount int) *string {
	return genericutil.Ptr(RandSentence(wordCount))
}

// RandString generates a random alphanumeric string with a given length.
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

// RandStringPtr returns a pointer to a random alphanumeric string with a given length.
func RandStringPtr(length int) *string {
	return genericutil.Ptr(RandString(length))
}

// RandTime generates a random time between `start` and `end`.
func RandTime(start, end time.Time) time.Time {
	if start.After(end) {
		start, end = end, start
	}
	duration := rand.Int63n(end.Unix() - start.Unix())
	return time.Unix(start.Unix()+duration, 0)
}

// RandTimePtr returns a pointer to a random time between `start` and `end`.
func RandTimePtr(start, end time.Time) *time.Time {
	t := RandTime(start, end)
	return &t
}

// RandURL generates a random URL with a random domain.
func RandURL() string {
	domains := []string{"example.com", "test.com", "dummy.net", "sample.org"}
	return "https://" + RandString(8) + "." + domains[rand.Intn(len(domains))]
}

// RandURLPtr returns a pointer to a random URL with a random domain.
func RandURLPtr() *string {
	return genericutil.Ptr(RandURL())
}

// UUID generates a random UUID v7.
func UUID() string {
	return uuid.Must(uuid.NewV7()).String()
}

// UUIDPtr returns a pointer to a random UUID v7.
func UUIDPtr() *string {
	return genericutil.Ptr(uuid.Must(uuid.NewV7()).String())
}
