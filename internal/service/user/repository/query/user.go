package query

// Dynamic SQL
type User interface {
	// SELECT id from @@table WHERE id IN @ids;
	FindByIdsIn(ids ...string) ([]string, error)
}
