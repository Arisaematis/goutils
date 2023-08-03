package redis

import (
	"fmt"
	"math/rand"
	"testing"
)

type SquareJob struct {
	Student Student
}

type Student struct {
	Id     int
	Name   string
	Age    int
	Scores []int
}

func (s *SquareJob) Execute() error {
	var st Student
	fmt.Printf("the result is student[Id:%v Name:%v Aget:%v Scores:%v]\n", st.Id, st.Name, st.Age, st.Scores)
	return nil
}

func createStudents() []Student {
	names := []string{"Tom", "Kate", "Lucy", "Jim", "Jack", "King", "Lee", "Mask"}
	students := make([]Student, 10)
	rnd := func(start, end int) int { return rand.Intn(end-start) + start }
	for i := 0; i < 10; i++ {
		students[i] = Student{
			Id:     i + 1,
			Name:   names[rand.Intn(len(names))],
			Age:    rnd(15, 26),
			Scores: []int{rnd(60, 100), rnd(60, 100), rnd(60, 100)},
		}
	}
	return students
}
func TestRedisMq(t *testing.T) {

}
