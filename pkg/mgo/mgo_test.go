package mgo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"k8s.io/klog/v2"
	"math/rand"
	"testing"
	"time"
)

// Get 查询条件
type Get struct {
	TimeQuantum TimeQuantum `json:"-" search:"-,omitempty" bson:"time"`
	Name        string      `json:"name" search:"name,vague,omitempty"`
	Id          string      `json:"id" search:"id,omitempty"`
	Ids         []string    `json:"ids" search:"id,in,omitempty"`
}

type TimeQuantum struct {
	StartTime time.Time `json:"startTime" search:"startTime"`
	EndTime   time.Time `json:"endTime" search:"endTime"`
}

func (timeQuantum TimeQuantum) SearchIn() (string, bson.M) {
	return "startTime", bson.M{
		"$gte": timeQuantum.StartTime,
		"$lte": timeQuantum.EndTime,
	}
}

var (
	TestCollection = NewMgo("test")
)

type Test struct {
	Id        int       `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	StartTime time.Time `json:"startTime" bson:"startTime"`
	Value     int       `json:"value" bson:"value"`
}

type Client interface {
	Insert(ctx context.Context, test []Test) error
	Get(ctx context.Context, get Get) ([]bson.M, error)
}

func NewClient() Client {
	return &clientImpl{}
}

type clientImpl struct {
}

func (c clientImpl) Insert(ctx context.Context, test []Test) error {
	var data []interface{}
	for _, t := range test {
		data = append(data, t)
	}
	_, err := TestCollection.InsertMany(ctx, data)
	return err
}

func (c clientImpl) Get(ctx context.Context, get Get) (result []bson.M, err error) {
	err = TestCollection.Match(get).Aggregation(ctx, &result)
	return
}

func createTest() []Test {
	names := []string{"Tom", "Kate", "Lucy", "Jim", "Jack", "King", "Lee", "Mask"}
	test := make([]Test, 10)
	rnd := func(start, end int) int { return rand.Intn(end-start) + start }
	for i := 0; i < 10; i++ {
		test[i] = Test{
			Id:        i + 1,
			Name:      names[rand.Intn(len(names))],
			Value:     rnd(15, 26),
			StartTime: time.Now(),
		}
	}
	return test
}

func TestInsert(t *testing.T) {
	test := createTest()
	err := NewClient().Insert(context.Background(), test)
	if err != nil {
		klog.Error(err)
	}
}

func TestGet(t *testing.T) {
	get, err := NewClient().Get(context.Background(), Get{Name: "Tom"})
	if err != nil {
		klog.Error(err)
	}
	klog.Info(get)
}
