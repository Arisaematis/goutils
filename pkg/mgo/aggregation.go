package mgo

import (
	"context"
	"encoding/json"
	"goutils/pkg/code"
	"k8s.io/klog/v2"

	"github.com/lstack-org/go-web-framework/pkg/req"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *Mgo) clone() *Mgo {
	clone := &Mgo{
		Collection: m.Collection,
		Pip:        make([]bson.D, 0),
	}
	for ix := range m.Pip {
		clone.Pip = append(clone.Pip, m.Pip[ix])
	}
	return clone
}

func (m *Mgo) MatchBson(d bson.D) *Mgo {
	clone := m.clone()
	if len(d) > 0 {
		clone.Pip = append(clone.Pip, bson.D{
			{Key: "$match", Value: d},
		})
	}
	return clone
}

func (m *Mgo) Match(args interface{}) *Mgo {
	filter, structToBsonErr := StructToBson(args, "search", "SearchIn")
	if structToBsonErr != nil {
		klog.Error("StructToBson failed because err: ", structToBsonErr)
	}
	clone := m.clone()
	if len(filter) > 0 {
		clone.Pip = append(clone.Pip, bson.D{
			{Key: "$match", Value: filter},
		})
	}
	return clone
}

func (m *Mgo) LookUp(lookUp ...bson.E) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{
			Key:   "$lookup",
			Value: lookUp,
		},
	})
	return clone
}

func (m *Mgo) Sort(sort bson.D) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$sort", Value: sort},
	})
	return clone
}

func (m *Mgo) LimitNum(limiter int) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$limit", Value: limiter},
	})
	return clone
}

func (m *Mgo) SkipNum(skipper int) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$skip", Value: skipper},
	})
	return clone
}

func NewPage(page req.Paging) (pip mongo.Pipeline) {
	if page.Page > 0 {
		pip = append(pip,
			bson.D{
				{Key: "$skip", Value: page.PageSize * (page.Page - 1)},
			})
	}
	if page.PageSize > 0 {
		pip = append(pip,
			bson.D{
				{Key: "$limit", Value: page.PageSize},
			})
	}
	return pip
}

func NewBaseQuery(page req.Paging) mongo.Pipeline {
	baseQuery := NewPage(page)
	baseQuery = append(baseQuery, NewSortForCreateTime())
	return baseQuery
}

func (m *Mgo) SetPage(page req.Paging) *Mgo {
	clone := m.clone()
	newPage := NewPage(page)
	clone.Pip = append(clone.Pip, newPage...)
	return clone
}

func (m *Mgo) Unwind(unwind string) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$unwind", Value: unwind},
	})
	return clone
}

func (m *Mgo) Group(group ...bson.E) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$group", Value: group},
	})
	return clone
}
func (m *Mgo) GroupM(group bson.M) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$group", Value: group},
	})
	return clone
}

func (m *Mgo) Project(group ...bson.E) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$project", Value: group},
	})
	return clone
}
func (m *Mgo) ProjectM(group bson.M) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$project", Value: group},
	})
	return clone
}

func (m *Mgo) AddFields(add ...bson.E) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$addFields", Value: add},
	})
	return clone
}

func (m *Mgo) ReplaceRoot(replaceRoot bson.M) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{Key: "$replaceRoot", Value: replaceRoot},
	})
	return clone
}

func NewCount() (pip mongo.Pipeline) {
	pip = append(pip, bson.D{
		{
			Key:   "$count",
			Value: "total",
		},
	})
	return
}

func (m *Mgo) Count() *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, NewCount()...)
	return clone
}

func NewSortForCreateTime() (bsonD bson.D) {
	bsonD = bson.D{
		{
			Key: "$sort",
			Value: bson.D{
				{
					Key:   "createTime",
					Value: -1,
				},
			},
		},
	}
	return
}

func (m *Mgo) SortForCreateTime() *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, NewSortForCreateTime())
	return clone
}

func (m *Mgo) Facet(bsonM bson.M) *Mgo {
	clone := m.clone()
	clone.Pip = append(clone.Pip, bson.D{
		{
			Key:   "$facet",
			Value: bsonM,
		},
	})
	return clone
}

func (m *Mgo) AggregationWithOpts(ctx context.Context, results interface{}, opts ...*options.AggregateOptions) error {
	aggregate, err := m.Aggregate(ctx, m.Pip, opts...)
	pipJson, _ := json.Marshal(m.Pip)
	klog.V(6).Info("mongo Aggregation bson is :", string(pipJson))
	if err != nil {
		return err
	}
	return aggregate.All(ctx, results)
}

func (m *Mgo) Aggregation(ctx context.Context, results interface{}) error {
	aggregate, err := m.Aggregate(ctx, m.Pip)
	pipJson, _ := json.Marshal(m.Pip)
	klog.V(6).Info("mongo Aggregation bson is :", string(pipJson))
	if err != nil {
		return err
	}
	return aggregate.All(ctx, results)
}

// AggregationFindById 用于详情查询
func (m *Mgo) AggregationFindById(ctx context.Context, result interface{}) error {
	pipJson, _ := json.Marshal(m.Pip)
	klog.V(6).Info("mongo Aggregation bson is :", string(pipJson))
	aggregate, err := m.Aggregate(ctx, m.Pip)
	if err != nil {
		return err
	}
	if aggregate.Next(ctx) {
		return aggregate.Decode(result)
	}
	// 异常抛出
	return code.ResourceIdNotFoundError
}

func (m *Mgo) AggregationList(ctx context.Context, page req.Paging, results interface{}) (total int, err error) {
	var listRes ListRes
	err = m.Facet(
		bson.M{
			"total": NewCount(),
			"items": NewBaseQuery(page),
		}).Aggregation(ctx, &listRes)
	if err != nil {
		return
	}
	err = listRes.Convert(results)
	if err != nil {
		return
	}
	return listRes.GetTotal(), nil
}

func (m *Mgo) Delete(ctx context.Context, args interface{}) error {
	deleteInfos, structToBsonErr := StructToBson(args, "delete", "DeleteIn")
	if structToBsonErr != nil {
		return structToBsonErr
	}
	klog.V(6).Info("delete mongo bson is : ", deleteInfos)
	_, err := m.DeleteMany(ctx, deleteInfos)
	return err
}
