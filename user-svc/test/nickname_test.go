package test

import (
	"fmt"
	"testing"
	"time"

	"math/rand/v2"
)

var adjectives = []string{
	"Swift", "Silent", "Crazy", "Angry", "Happy", "Sneaky", "Brave", "Loyal", "Wild", "Lazy",
}

var nouns = []string{
	"Tiger", "Panda", "Wolf", "Falcon", "Ninja", "Coder", "Wizard", "Ghost", "Samurai", "Knight",
}

func TestNicknameGenerator(t *testing.T) {
	nickname := generateNickname()
	fmt.Println("Generated Nickname:", nickname)
}

func generateNickname() string {
	rand.Int64N((time.Now().Unix()))
	adj := adjectives[rand.IntN(len(adjectives))]
	noun := nouns[rand.IntN(len(nouns))]
	number := rand.IntN(999) + 1 // Random number between 1-999

	return fmt.Sprintf("%s%s%d", adj, noun, number)
}
