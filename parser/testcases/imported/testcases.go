package imported

type (
	SizeCache       = int32
	WeakFields      = map[int32]Message
	UnknownFields   = unknownFieldsA // TODO: switch to unknownFieldsB
	unknownFieldsA  = []byte
	unknownFieldsB  = *[]byte
	ExtensionFields = map[int32]ExtensionField
)

type Message interface {
	Test()
}

type ExtensionField struct {
	Field int
}
