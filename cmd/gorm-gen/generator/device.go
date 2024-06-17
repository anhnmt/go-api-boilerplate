package generator

// Dynamic SQL
type Device interface {
	// SELECT id from @@table WHERE id IN @ids;
	FindByIdsIn(ids ...string) ([]string, error)
}
