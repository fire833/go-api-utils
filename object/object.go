package vtobject

import "google.golang.org/protobuf/reflect/protoreflect"

// Object is the primary interface for data objects that are managed and served by applications. Each Object, assuming
// it has the correct comment tags, will have an ObjectManager that is generated and imported into the application
// to allow for HTTP and Swagger support for the object. The manager will also have metrics for that object that can
// be updated and collected at runtime, including counts for each operation completed by that manager (ie number of GETs,
// PUTs, etc.)
type Object interface {
	// ProtoMessage allows Objects to be marshalled to protobuf automatically. This allows for an additional
	// MIME type to be supported by the API for free, essentially.
	protoreflect.ProtoMessage
}

type ObjectList interface {
	protoreflect.ProtoMessage

	GetItems() []Object
}
