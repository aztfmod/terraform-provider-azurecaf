package azurecaf

import (
	"strings"
)

type NameBuilder struct {
	MaxLength int
	Separator string
	content   []NameSegment
}

type NameSegment struct {
	Value   string
	Include bool
}

func NewNameBuilder(maxLength int, separator string) *NameBuilder {
	return &NameBuilder{
		MaxLength: maxLength,
		Separator: separator,
		content:   []NameSegment{},
	}
}

func (b NameBuilder) include(segment string) bool {
	trimmedLength := b.GetTrimmedLength()
	segmentLength := len(segment)
	if trimmedLength == 0 {
		return segmentLength <= b.MaxLength
	}
	return trimmedLength+segmentLength+len(b.Separator) <= b.MaxLength
}

func (b *NameBuilder) Append(segment string) {
	b.content = append(b.content, NameSegment{Value: segment, Include: b.include(segment)})
}

func (b *NameBuilder) Prepend(segment string) {
	b.content = append([]NameSegment{{Value: segment, Include: b.include(segment)}}, b.content...)
}

func (b NameBuilder) GetLength() int {
	return len(b.GetName())
}

func (b NameBuilder) GetTrimmedLength() int {
	return len(b.GetTrimmedName())
}

func (b NameBuilder) GetName() string {
	return strings.Join(b.getAllContent(), b.Separator)
}

func (b NameBuilder) GetTrimmedName() string {
	return strings.Join(b.getIncludedContent(), b.Separator)
}

func (b NameBuilder) getIncludedContent() []string {
	values := make([]string, 0, len(b.content))
	for _, segment := range b.content {
		if segment.Include {
			values = append(values, segment.Value)
		}
	}
	return values
}

func (b NameBuilder) getAllContent() []string {
	values := make([]string, 0, len(b.content))
	for _, segment := range b.content {
		values = append(values, segment.Value)
	}
	return values
}
