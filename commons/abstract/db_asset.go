package abstract

// types used in parsing

// DBInfo ... Name and description for a database
type DBInfo struct {
	Name    string
	Comment string
}

// ColumnInfo ... Type and description for a table column
type ColumnInfo struct {
	Type    string
	Comment string
}

// TableInfo ... Name, schema and description for a table
type TableInfo struct {
	Name    string
	Schema  map[string]ColumnInfo
	Comment string
}

// Partial constructor for DBInfo
func GetDBInfoByName(dbName string) (DBInfo, error) {
	dbInfo := DBInfo{}
	dbInfo.Name = dbName
	dbInfo.Comment = ""
	return dbInfo, nil
}

// Partial constructor for TableInfo
func GetTableInfoByName(tableName string) (TableInfo, error) {

	tableInfo := TableInfo{}
	tableInfo.Name = tableName
	tableInfo.Comment = ""

	return tableInfo, nil
}

// structs shall be unexported by default
type databaseBuilder struct{ asset Asset }

// NewDatabaseBuilder ... builder for a database asset type
func NewDatabaseBuilder() *databaseBuilder {
	builder := &databaseBuilder{}
	builder.asset.Type = _Database
	return builder
}

func (b *databaseBuilder) SetName(name string) *databaseBuilder {
	b.asset.Name = name
	return b
}

func (b *databaseBuilder) SetDescription(description string) *databaseBuilder {
	b.asset.Description = description
	return b
}

func (b *databaseBuilder) SetTags(tags []string) *databaseBuilder {
	if b.asset.Tags == nil {
		b.asset.Tags = []string{}
	}
	b.asset.Tags = append(b.asset.Tags, tags...)
	return b
}

func (b *databaseBuilder) Build() (*Asset, error) {
	if err := b.asset.Validate(); err != nil {
		return nil, err
	}
	return &b.asset, nil
}

func (db *DBInfo) BuildAsset() (*Asset, error) {
	return NewDatabaseBuilder().
		SetName(db.Name).
		SetDescription(db.Comment).
		Build()
}

type tableBuilder struct{ asset Asset }

// NewTableBuilder ... table builder
func NewTableBuilder() *tableBuilder {
	builder := &tableBuilder{}
	builder.asset.Type = _Table
	return builder
}

func (b *tableBuilder) SetName(name string) *tableBuilder {
	b.asset.Name = name
	return b
}

func (b *tableBuilder) SetDescription(description string) *tableBuilder {
	b.asset.Description = description
	return b
}

func (b *tableBuilder) SetSchema(schema map[string]ColumnInfo) *tableBuilder {
	if b.asset.Labels == nil {
		b.asset.Labels = make(map[string]interface{})
	}
	b.asset.Labels[L_SCHEMA] = schema
	return b
}

func (b *tableBuilder) SetTags(tags []string) *tableBuilder {
	if b.asset.Tags == nil {
		b.asset.Tags = []string{}
	}
	b.asset.Tags = append(b.asset.Tags, tags...)
	return b
}

func (b *tableBuilder) Build() (*Asset, error) {
	if err := b.asset.Validate(); err != nil {
		return nil, err
	}
	return &b.asset, nil
}

func (tb *TableInfo) BuildAsset() (*Asset, error) {
	return NewTableBuilder().
		SetName(tb.Name).
		SetDescription(tb.Comment).
		SetSchema(tb.Schema).
		Build()
}
