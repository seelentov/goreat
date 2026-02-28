package entities

type FieldType int

const (
	FieldTypeString FieldType = iota
	FieldTypeInt
	FieldTypeFloat
	FieldTypeBool
	FieldTypeDateTime
	FieldTypeSource
)

type ContainerType int

const (
	ContainerTypeSingle ContainerType = iota
	ContainerTypeList
)

type Field struct {
	Name string
}

type FieldInfo struct {
	FieldType     FieldType
	ContainerType ContainerType
}

type DefinitionField struct {
	Field
	FieldInfo
}

func NewDefinitionField(name string, fType FieldType, cType ContainerType) *DefinitionField {
	return &DefinitionField{
		Field: Field{
			Name: name,
		},
		FieldInfo: FieldInfo{
			ContainerType: cType,
			FieldType:     fType,
		},
	}
}
