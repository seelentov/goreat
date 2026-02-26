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

type DefinitionField struct {
	Field

	Type          FieldType
	ContainerType ContainerType
}

func NewDefinitionField(name string, fType FieldType, cType ContainerType) *DefinitionField {
	return &DefinitionField{
		Field: Field{
			Name: name,
		},
		ContainerType: cType,
		Type:          fType,
	}
}

type FieldValueInfo struct {
	FieldType     FieldType
	ContainerType ContainerType
}

type FieldValueData struct {
	FieldValueInfo

	Value any
}
