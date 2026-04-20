package gordle

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestGameAsk(t *testing.T) {
	tt := map[string]struct {
		input string
		want  []rune
	}{
		"5 characters in english": {
			input: "HELLO",
			want:  []rune("HELLO"),
		},
		"5 characters in arabic": {
			input: "مرحبا",
			want:  []rune("مرحبا"),
		},
		"5 characters in japanese": {
			input: "こんにちは",
			want:  []rune("こんにちは"),
		},
		"3 characters in japanese": {
			input: "こんに\nこんにちは",
			want:  []rune("こんにちは"),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			g := New(strings.NewReader(tc.input), string(tc.want), 0)

			got := g.ask()
			if !slices.Equal(got, tc.want) {
				t.Errorf("got = %v, want %v", string(got), string(tc.want))
			}
		})
	}
}

func TestGameValidateGuess(t *testing.T) {
	tt := map[string]struct {
		word      []rune
		wantError error
	}{
		"normal": {
			word:      []rune("GUESS"),
			wantError: nil,
		},
		"too long": {
			word:      []rune("POCKET"),
			wantError: errInvalidWordLength,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			g := New(nil, string(""), 0)

			err := g.validateGuess(tc.word)
			if err == nil && tc.wantError == nil {
				return
			}
			if !errors.Is(err, errInvalidWordLength) {
				t.Errorf("%c, expected %q, got %q", tc.word, tc.wantError, err)
			}
		})
	}
}

func TestComputeFeedback(t *testing.T) {
	tt := map[string]struct {
		guess            string
		solution         string
		expectedFeedback feedback
	}{
		"normal": {
			guess:    "hello",
			solution: "hello",
			expectedFeedback: feedback{
				CORRECT_POSITION,
				CORRECT_POSITION,
				CORRECT_POSITION,
				CORRECT_POSITION,
				CORRECT_POSITION,
			},
		},
		"double character": {
			guess:    "creep",
			solution: "speed",
			expectedFeedback: feedback{
				ABSENT_CHARACTER,
				ABSENT_CHARACTER,
				CORRECT_POSITION,
				CORRECT_POSITION,
				WRONG_POSITION,
			},
		},
		"two identical, but not in the right position": {
			guess:    "hlleo",
			solution: "hello",
			expectedFeedback: feedback{
				CORRECT_POSITION,
				WRONG_POSITION,
				CORRECT_POSITION,
				WRONG_POSITION,
				CORRECT_POSITION},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			fb := computeFeedback([]rune(tc.guess), []rune(tc.solution))
			if !tc.expectedFeedback.Equal(fb) {
				t.Errorf(
					"guess: %q, got the wrong feedback, wanted %v, got %v",
					tc.guess, tc.expectedFeedback, fb)
			}
		})
	}
}
