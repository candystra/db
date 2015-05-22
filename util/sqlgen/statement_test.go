package sqlgen

import (
	"regexp"
	"strings"
	"testing"
)

var (
	reInvisible = regexp.MustCompile(`[\t\n\r]`)
	reSpace     = regexp.MustCompile(`\s+`)
)

func trim(a string) string {
	a = reInvisible.ReplaceAllString(strings.TrimSpace(a), " ")
	a = reSpace.ReplaceAllString(strings.TrimSpace(a), " ")
	return a
}

func TestTruncateTable(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Truncate,
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `TRUNCATE TABLE "table_name"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestDropTable(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  DropTable,
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `DROP TABLE "table_name"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestDropDatabase(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:     DropDatabase,
		Database: &Database{Name: "table_name"},
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `DROP DATABASE "table_name"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestCount(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Count,
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT COUNT(1) AS _t FROM "table_name"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestCountRelation(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Count,
		Table: TableWithName("information_schema.tables"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT COUNT(1) AS _t FROM "information_schema"."tables"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestCountWhere(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Count,
		Table: TableWithName("table_name"),
		Where: NewWhere(
			&ColumnValue{Column: Column{Name: "a"}, Operator: "=", Value: NewValue(RawValue("7"))},
		),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT COUNT(1) AS _t FROM "table_name" WHERE ("a" = 7)`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectStarFrom(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Select,
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT * FROM "table_name"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectStarFromAlias(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Select,
		Table: TableWithName("table.name AS foo"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT * FROM "table"."name" AS "foo"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectStarFromRawWhere(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Select,
		Table: TableWithName("table.name AS foo"),
		Where: NewWhere(
			&Raw{Value: "foo.id = bar.foo_id"},
		),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT * FROM "table"."name" AS "foo" WHERE (foo.id = bar.foo_id)`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}

	stmt = Statement{
		Type:  Select,
		Table: TableWithName("table.name AS foo"),
		Where: NewWhere(
			&Raw{Value: "foo.id = bar.foo_id"},
			&Raw{Value: "baz.id = exp.baz_id"},
		),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT * FROM "table"."name" AS "foo" WHERE (foo.id = bar.foo_id AND baz.id = exp.baz_id)`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectStarFromMany(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Select,
		Table: TableWithName("first.table AS foo, second.table as BAR, third.table aS baz"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT * FROM "first"."table" AS "foo", "second"."table" AS "BAR", "third"."table" AS "baz"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectArtistNameFrom(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Select,
		Table: TableWithName("artist"),
		Columns: NewColumns(
			Column{Name: "artist.name"},
		),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "artist"."name" FROM "artist"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectRawFrom(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Select,
		Table: TableWithName(`artist`),
		Columns: NewColumns(
			Column{Name: `artist.name`},
			Column{Name: Raw{Value: `CONCAT(artist.name, " ", artist.last_name)`}},
		),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "artist"."name", CONCAT(artist.name, " ", artist.last_name) FROM "artist"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectFieldsFrom(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectFieldsFromWithLimitOffset(t *testing.T) {
	var s, e string
	var stmt Statement

	// LIMIT only.
	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		Limit: 42,
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" LIMIT 42`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}

	// OFFSET only.
	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		Offset: 17,
		Table:  TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" OFFSET 17`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}

	// LIMIT AND OFFSET.
	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		Limit:  42,
		Offset: 17,
		Table:  TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" LIMIT 42 OFFSET 17`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestStatementGroupBy(t *testing.T) {
	var s, e string
	var stmt Statement

	// Simple GROUP BY
	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		GroupBy: NewGroupBy(
			Column{Name: "foo"},
		),
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" GROUP BY "foo"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}

	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		GroupBy: NewGroupBy(
			Column{Name: "foo"},
			Column{Name: "bar"},
		),
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" GROUP BY "foo", "bar"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectFieldsFromWithOrderBy(t *testing.T) {
	var s, e string
	var stmt Statement

	// Simple ORDER BY
	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		OrderBy: NewOrderBy(
			NewSortColumns(
				SortColumn{Column: Column{Name: "foo"}},
			),
		),
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" ORDER BY "foo"`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}

	// ORDER BY field ASC
	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		OrderBy: NewOrderBy(
			NewSortColumns(
				SortColumn{Column: Column{Name: "foo"}, Order: Ascendent},
			),
		),
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" ORDER BY "foo" ASC`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}

	// ORDER BY field DESC
	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		OrderBy: NewOrderBy(
			NewSortColumns(
				SortColumn{Column: Column{Name: "foo"}, Order: Descendent},
			),
		),
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" ORDER BY "foo" DESC`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}

	// ORDER BY many fields
	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		OrderBy: NewOrderBy(
			NewSortColumns(
				SortColumn{Column: Column{Name: "foo"}, Order: Descendent},
				SortColumn{Column: Column{Name: "bar"}, Order: Ascendent},
				SortColumn{Column: Column{Name: "baz"}, Order: Descendent},
			),
		),
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" ORDER BY "foo" DESC, "bar" ASC, "baz" DESC`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}

	// ORDER BY function
	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		OrderBy: NewOrderBy(
			NewSortColumns(
				SortColumn{Column: Column{Name: Raw{Value: "FOO()"}}, Order: Descendent},
				SortColumn{Column: Column{Name: Raw{Value: "BAR()"}}, Order: Ascendent},
			),
		),
		Table: TableWithName("table_name"),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" ORDER BY FOO() DESC, BAR() ASC`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectFieldsFromWhere(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		Table: TableWithName("table_name"),
		Where: NewWhere(
			&ColumnValue{Column: Column{Name: "baz"}, Operator: "=", Value: NewValue(99)},
		),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" WHERE ("baz" = '99')`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestSelectFieldsFromWhereLimitOffset(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type: Select,
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		Table: TableWithName("table_name"),
		Where: NewWhere(
			&ColumnValue{Column: Column{Name: "baz"}, Operator: "=", Value: NewValue(99)},
		),
		Limit:  10,
		Offset: 23,
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `SELECT "foo", "bar", "baz" FROM "table_name" WHERE ("baz" = '99') LIMIT 10 OFFSET 23`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestDelete(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Delete,
		Table: TableWithName("table_name"),
		Where: NewWhere(
			&ColumnValue{Column: Column{Name: "baz"}, Operator: "=", Value: NewValue(99)},
		),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `DELETE FROM "table_name" WHERE ("baz" = '99')`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestUpdate(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Update,
		Table: TableWithName("table_name"),
		ColumnValues: NewColumnValues(
			ColumnValue{Column: Column{Name: "foo"}, Operator: "=", Value: NewValue(76)},
		),
		Where: NewWhere(
			&ColumnValue{Column: Column{Name: "baz"}, Operator: "=", Value: NewValue(99)},
		),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `UPDATE "table_name" SET "foo" = '76' WHERE ("baz" = '99')`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}

	stmt = Statement{
		Type:  Update,
		Table: TableWithName("table_name"),
		ColumnValues: NewColumnValues(
			ColumnValue{Column: Column{Name: "foo"}, Operator: "=", Value: NewValue(76)},
			ColumnValue{Column: Column{Name: "bar"}, Operator: "=", Value: NewValue(Raw{Value: "88"})},
		),
		Where: NewWhere(
			&ColumnValue{Column: Column{Name: "baz"}, Operator: "=", Value: NewValue(99)},
		),
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `UPDATE "table_name" SET "foo" = '76', "bar" = 88 WHERE ("baz" = '99')`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestInsert(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Insert,
		Table: TableWithName("table_name"),
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		Values: Values{
			Value{V: "1"},
			Value{V: 2},
			Value{V: Raw{Value: "3"}},
		},
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `INSERT INTO "table_name" ("foo", "bar", "baz") VALUES ('1', '2', 3)`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func TestInsertExtra(t *testing.T) {
	var s, e string
	var stmt Statement

	stmt = Statement{
		Type:  Insert,
		Table: TableWithName("table_name"),
		Extra: "RETURNING id",
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		Values: Values{
			Value{V: "1"},
			Value{V: 2},
			Value{V: Raw{Value: "3"}},
		},
	}

	s = trim(stmt.Compile(defaultTemplate))
	e = `INSERT INTO "table_name" ("foo", "bar", "baz") VALUES ('1', '2', 3) RETURNING id`

	if s != e {
		t.Fatalf("Got: %s, Expecting: %s", s, e)
	}
}

func BenchmarkStatementSimpleQuery(b *testing.B) {
	stmt := Statement{
		Type:  Count,
		Table: TableWithName("table_name"),
		Where: NewWhere(
			&ColumnValue{Column: Column{Name: "a"}, Operator: "=", Value: NewValue(Raw{Value: "7"})},
		),
	}

	for i := 0; i < b.N; i++ {
		_ = stmt.Compile(defaultTemplate)
	}
}

func BenchmarkStatementSimpleQueryHash(b *testing.B) {
	stmt := Statement{
		Type:  Count,
		Table: TableWithName("table_name"),
		Where: NewWhere(
			&ColumnValue{Column: Column{Name: "a"}, Operator: "=", Value: NewValue(Raw{Value: "7"})},
		),
	}

	for i := 0; i < b.N; i++ {
		_ = stmt.Hash()
	}
}

func BenchmarkStatementSimpleQueryNoCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stmt := Statement{
			Type:  Count,
			Table: TableWithName("table_name"),
			Where: NewWhere(
				&ColumnValue{Column: Column{Name: "a"}, Operator: "=", Value: NewValue(Raw{Value: "7"})},
			),
		}
		_ = stmt.Compile(defaultTemplate)
	}
}

func BenchmarkStatementComplexQuery(b *testing.B) {
	stmt := Statement{
		Type:  Insert,
		Table: TableWithName("table_name"),
		Columns: NewColumns(
			Column{Name: "foo"},
			Column{Name: "bar"},
			Column{Name: "baz"},
		),
		Values: Values{
			Value{V: "1"},
			Value{V: 2},
			Value{V: Raw{Value: "3"}},
		},
	}

	for i := 0; i < b.N; i++ {
		_ = stmt.Compile(defaultTemplate)
	}
}

func BenchmarkStatementComplexQueryNoCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stmt := Statement{
			Type:  Insert,
			Table: TableWithName("table_name"),
			Columns: NewColumns(
				Column{Name: "foo"},
				Column{Name: "bar"},
				Column{Name: "baz"},
			),
			Values: Values{
				Value{V: "1"},
				Value{V: 2},
				Value{V: Raw{Value: "3"}},
			},
		}
		_ = stmt.Compile(defaultTemplate)
	}
}
