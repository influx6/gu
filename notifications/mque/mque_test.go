package mque_test

import (
	"testing"

	"github.com/gu-io/gu/notifications/mque"
)

// failedMark is the Unicode codepoint for an X mark.
const failedMark = "\u2717"

// succeedMark is the Unicode codepoint for a check mark.
const succeedMark = "\u2713"

// TestQueue validates the behaviour of que api.
func TestQueue(t *testing.T) {
	t.Logf("Should be able to use a argument selective queue")
	{

		t.Logf("\tWhen giving a mque.Que and a string only allowed constraint")
		{

			passed := make(chan int)
			failed := make(chan int)

			q := mque.New()

			q.Q(func(letter string) {
				go func() { passed <- 1 }()
			})

			q.Q(func(item int) {
				go func() { failed <- 1 }()
			})

			q.Run("letter")

			select {
			case <-passed:
				t.Logf("\t%s\tShould have received a string", succeedMark)
			case <-failed:
				t.Errorf("\t%s\tShould have received a string", failedMark)
			}

		}

		t.Logf("\tWhen giving a mque.Que and a int only allowed constraint")
		{
			passed := make(chan int)
			failed := make(chan int)

			q := mque.New()

			q.Q(func(letter string) {
				go func() { failed <- 1 }()
			})

			q.Q(func(item int) {
				go func() { passed <- 1 }()
			})

			q.Run(30)

			select {
			case <-failed:
				t.Fatalf("\t%s\tShould have received a integer", failedMark)
			case <-passed:
				t.Logf("\t%s\tShould have received a integer", succeedMark)
			}
		}
	}
}

func BenchmarkQueue(b *testing.B) {
	b.ResetTimer()
	defer b.ReportAllocs()

	q := mque.New()

	q.Q(func(item int) {})
	q.Q(func(item int) {})

	for i := 0; i < b.N; i++ {
		q.Run(i)
	}
}

func BenchmarkQueueWithMultitypes(b *testing.B) {
	b.ResetTimer()
	defer b.ReportAllocs()

	q := mque.New()

	q.Q(func(item string) {})
	q.Q(func(item int) {})

	for i := 0; i < b.N; i++ {
		q.Run(i)
	}
}

// TestQueueEnd validates the behaviour of que subscriber End call.
func TestQueueEnd(t *testing.T) {
	t.Logf("Should be able to use a argument selective queue")
	{

		t.Logf("\tWhen needing to unsubscribe from a queue")
		{

			var count int
			var ended bool

			q := mque.New()

			q.Q(func(item int) {
				count++
			})

			sub := q.Q(func(item int) {
				count++
			}, func() {
				ended = true
			})

			q.Run(20)
			sub.End()

			q.Run(40)

			if count <= 2 || count >= 4 {
				t.Fatalf("\t%s\tShould have received only two events", failedMark)
			}
			t.Logf("\t%s\tShould have received only two events", succeedMark)

			if !ended {
				t.Fatalf("\t%s\tShould have ended subscriber", failedMark)
			}
			t.Logf("\t%s\tShould have ended subscriber", succeedMark)

		}
	}
}
