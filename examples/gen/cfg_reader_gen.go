// Code generated by proto_code_gen. DO NOT EDIT
package gen

import (
    
	"github.com/fish-tennis/proto_code_gen/examples/pb"
)


type ExampleR struct {
	v *pb.Example
}

func NewExampleR(src *pb.Example) *ExampleR {
	return &ExampleR{v:src}
}

func (r *ExampleR) Raw() *pb.Example {
	return r.v
}

func (r *ExampleR) GetInt32Field() int32 {
	return r.v.GetInt32Field()
}

func (r *ExampleR) GetStringField() string {
	return r.v.GetStringField()
}

func (r *ExampleR) GetFloatField() float32 {
	return r.v.GetFloatField()
}

func (r *ExampleR) GetInt32Slice() []int32 {
	src := r.v.GetInt32Slice()
    if src == nil {
        return nil
    }
    copySlice := make([]int32,len(src))
    copy(copySlice, src)
    return copySlice
}

func (r *ExampleR) GetStringSlice() []string {
	src := r.v.GetStringSlice()
    if src == nil {
        return nil
    }
    copySlice := make([]string,len(src))
    copy(copySlice, src)
    return copySlice
}


func (r *ExampleR) GetSingleChild() *ChildR {
	return NewChildR(r.v.GetSingleChild())
}

func (r *ExampleR) GetChildren() []*ChildR {
	src := r.v.GetChildren()
    if src == nil {
        return nil
    }
    sliceReader := make([]*ChildR,len(src))
    for i,v := range src {
        sliceReader[i] = NewChildR(v)
    }
    return sliceReader
}


type Example2R struct {
	v *pb.Example2
}

func NewExample2R(src *pb.Example2) *Example2R {
	return &Example2R{v:src}
}

func (r *Example2R) Raw() *pb.Example2 {
	return r.v
}

func (r *Example2R) GetInt32Field() int32 {
	return r.v.GetInt32Field()
}

func (r *Example2R) GetStringField() string {
	return r.v.GetStringField()
}


type ExampleWithoutTagR struct {
	v *pb.ExampleWithoutTag
}

func NewExampleWithoutTagR(src *pb.ExampleWithoutTag) *ExampleWithoutTagR {
	return &ExampleWithoutTagR{v:src}
}

func (r *ExampleWithoutTagR) Raw() *pb.ExampleWithoutTag {
	return r.v
}

func (r *ExampleWithoutTagR) GetInt32Field() int32 {
	return r.v.GetInt32Field()
}

func (r *ExampleWithoutTagR) GetStringField() string {
	return r.v.GetStringField()
}


type ChildR struct {
	v *pb.Child
}

func NewChildR(src *pb.Child) *ChildR {
	return &ChildR{v:src}
}

func (r *ChildR) Raw() *pb.Child {
	return r.v
}

func (r *ChildR) GetInt32Field() int32 {
	return r.v.GetInt32Field()
}

func (r *ChildR) GetStringField() string {
	return r.v.GetStringField()
}

