package validators

import "reflect"

// mockFieldLevel is a mock implementation of the validator.FieldLevel interface
type mockFieldLevel struct {
	value           reflect.Value
	fieldName       string
	structFieldName string
	param           string
	tag             string
}

// Implement Top method
func (m *mockFieldLevel) Top() reflect.Value {
	return m.value
}

// Implement Parent method
func (m *mockFieldLevel) Parent() reflect.Value {
	return m.value
}

// Implement Field method
func (m *mockFieldLevel) Field() reflect.Value {
	return m.value
}

// Implement FieldName method
func (m *mockFieldLevel) FieldName() string {
	return m.fieldName
}

// Implement StructFieldName method
func (m *mockFieldLevel) StructFieldName() string {
	return m.structFieldName
}

// Implement Param method
func (m *mockFieldLevel) Param() string {
	return m.param
}

// Implement GetTag method
func (m *mockFieldLevel) GetTag() string {
	return m.tag
}

// Implement ExtractType method
func (m *mockFieldLevel) ExtractType(field reflect.Value) (reflect.Value, reflect.Kind, bool) {
	return field, field.Kind(), true
}

// Implement GetStructFieldOK method
func (m *mockFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	return m.value, m.value.Kind(), true
}

// Implement GetStructFieldOKAdvanced method
func (m *mockFieldLevel) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	return val, val.Kind(), true
}

// Implement GetStructFieldOK2 method
func (m *mockFieldLevel) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return m.value, m.value.Kind(), true, true
}

// Implement GetStructFieldOKAdvanced2 method
func (m *mockFieldLevel) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return val, val.Kind(), true, true
}
