package helpers

import "bytes"

/*
  StringBuilder struct.
	  Usage:
		builder := NewStringBuilder()
		builder.Append("a").Append("b").Append("c")
		s := builder.String()
		println(s)
	  print:
		abc
*/

// StringBuilder ...
type StringBuilder struct {
	buffer bytes.Buffer
}

// NewStringBuilder ...
func NewStringBuilder() *StringBuilder {
	var builder StringBuilder
	return &builder
}

// Append ...
func (builder *StringBuilder) Append(s string) *StringBuilder {
	builder.buffer.WriteString(s)
	return builder
}

// AppendStrings ...
func (builder *StringBuilder) AppendStrings(ss ...string) *StringBuilder {
	for i := range ss {
		builder.buffer.WriteString(ss[i])
	}
	return builder
}

// Clear ...
func (builder *StringBuilder) Clear() *StringBuilder {
	var buffer bytes.Buffer
	builder.buffer = buffer
	return builder
}

// ToString ...
func (builder *StringBuilder) ToString() string {
	return builder.buffer.String()
}
