package query

// Dynamic SQL
type UserQuery interface {
	// SELECT id from @@table WHERE id IN @ids;
	FindByIdsIn(ids ...string) ([]string, error)
}
