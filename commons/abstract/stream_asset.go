package abstract

type StreamInfo struct {
	Name    string
	Comment string
	Schema  map[int]string
}

// Partial constructor for StreamInfo
func GetStreamInfoByName(streamName string) (StreamInfo, error) {
	streamInfo := StreamInfo{}
	streamInfo.Name = streamName
	streamInfo.Comment = ""
	streamInfo.Schema = make(map[int]string)
	return streamInfo, nil
}

func (si *StreamInfo) BuildAsset() (*Asset, error) {
	return NewStreamBuilder().
		SetName(si.Name).
		SetDescription(si.Comment).
		SetSchema(si.Schema).
		Build()
}

type streamBuilder struct{ asset Asset }

func NewStreamBuilder() *streamBuilder {
	builder := &streamBuilder{}
	builder.asset.Type = _Stream
	return builder
}

func (b *streamBuilder) SetName(name string) *streamBuilder {
	b.asset.Name = name
	return b
}

func (b *streamBuilder) SetDescription(description string) *streamBuilder {
	b.asset.Description = description
	return b
}

func (b *streamBuilder) SetTags(tags []string) *streamBuilder {
	if b.asset.Tags == nil {
		b.asset.Tags = []string{}
	}
	b.asset.Tags = append(b.asset.Tags, tags...)
	return b
}

func (b *streamBuilder) SetSchema(schema map[int]string) *streamBuilder {
	if b.asset.Labels == nil {
		b.asset.Labels = make(map[string]interface{})
	}
	b.asset.Labels[L_SCHEMA] = schema
	return b
}

func (b *streamBuilder) Build() (*Asset, error) {
	if err := b.asset.Validate(); err != nil {
		return nil, err
	}
	return &b.asset, nil
}
