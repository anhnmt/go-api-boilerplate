package generator

// Dynamic SQL
type Session interface {
	// SELECT id from @@table WHERE id IN @ids;
	FindByIdsIn(ids ...string) ([]string, error)
}
