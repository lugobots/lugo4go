// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.19.4
// source: physics.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

//Vector represent one direction on a cartesian plan
type Vector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Coordinate X to define the vector direction.
	X float64 `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	// Coordinate Y to define the vector direction.
	Y float64 `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *Vector) Reset() {
	*x = Vector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vector) ProtoMessage() {}

func (x *Vector) ProtoReflect() protoreflect.Message {
	mi := &file_physics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vector.ProtoReflect.Descriptor instead.
func (*Vector) Descriptor() ([]byte, []int) {
	return file_physics_proto_rawDescGZIP(), []int{0}
}

func (x *Vector) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Vector) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

// Point represents one position on the cartesian plan of the game field.
// The coordinates start at the left bottom corner from the top view .
type Point struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Distance from the Y axis to right.
	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	// Distance from the X axis to up.
	Y int32 `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *Point) Reset() {
	*x = Point{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Point) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Point) ProtoMessage() {}

func (x *Point) ProtoReflect() protoreflect.Message {
	mi := &file_physics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Point.ProtoReflect.Descriptor instead.
func (*Point) Descriptor() ([]byte, []int) {
	return file_physics_proto_rawDescGZIP(), []int{1}
}

func (x *Point) GetX() int32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Point) GetY() int32 {
	if x != nil {
		return x.Y
	}
	return 0
}

// Velocity is a tuple with the direction (a vector) an a speed (float) values.
// It defines the velocity of an object.
type Velocity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Direction is a normalised vector that indicates the element direction
	Direction *Vector `protobuf:"bytes,1,opt,name=direction,proto3" json:"direction,omitempty"`
	// Speed of the element.
	Speed float64 `protobuf:"fixed64,2,opt,name=speed,proto3" json:"speed,omitempty"`
}

func (x *Velocity) Reset() {
	*x = Velocity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physics_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Velocity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Velocity) ProtoMessage() {}

func (x *Velocity) ProtoReflect() protoreflect.Message {
	mi := &file_physics_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Velocity.ProtoReflect.Descriptor instead.
func (*Velocity) Descriptor() ([]byte, []int) {
	return file_physics_proto_rawDescGZIP(), []int{2}
}

func (x *Velocity) GetDirection() *Vector {
	if x != nil {
		return x.Direction
	}
	return nil
}

func (x *Velocity) GetSpeed() float64 {
	if x != nil {
		return x.Speed
	}
	return 0
}

var File_physics_proto protoreflect.FileDescriptor

var file_physics_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x6c, 0x75, 0x67, 0x6f, 0x22, 0x24, 0x0a, 0x06, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12,
	0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a,
	0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x79, 0x22, 0x23, 0x0a, 0x05, 0x50,
	0x6f, 0x69, 0x6e, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x79,
	0x22, 0x4c, 0x0a, 0x08, 0x56, 0x65, 0x6c, 0x6f, 0x63, 0x69, 0x74, 0x79, 0x12, 0x2a, 0x0a, 0x09,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0c, 0x2e, 0x6c, 0x75, 0x67, 0x6f, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x09, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x70, 0x65, 0x65,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x42, 0x23,
	0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x75, 0x67,
	0x6f, 0x62, 0x6f, 0x74, 0x73, 0x2f, 0x6c, 0x75, 0x67, 0x6f, 0x34, 0x67, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_physics_proto_rawDescOnce sync.Once
	file_physics_proto_rawDescData = file_physics_proto_rawDesc
)

func file_physics_proto_rawDescGZIP() []byte {
	file_physics_proto_rawDescOnce.Do(func() {
		file_physics_proto_rawDescData = protoimpl.X.CompressGZIP(file_physics_proto_rawDescData)
	})
	return file_physics_proto_rawDescData
}

var file_physics_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_physics_proto_goTypes = []interface{}{
	(*Vector)(nil),   // 0: lugo.Vector
	(*Point)(nil),    // 1: lugo.Point
	(*Velocity)(nil), // 2: lugo.Velocity
}
var file_physics_proto_depIdxs = []int32{
	0, // 0: lugo.Velocity.direction:type_name -> lugo.Vector
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_physics_proto_init() }
func file_physics_proto_init() {
	if File_physics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_physics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vector); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_physics_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Point); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_physics_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Velocity); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_physics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_physics_proto_goTypes,
		DependencyIndexes: file_physics_proto_depIdxs,
		MessageInfos:      file_physics_proto_msgTypes,
	}.Build()
	File_physics_proto = out.File
	file_physics_proto_rawDesc = nil
	file_physics_proto_goTypes = nil
	file_physics_proto_depIdxs = nil
}
