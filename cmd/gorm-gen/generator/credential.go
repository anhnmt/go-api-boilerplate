package generator

// Dynamic SQL
type Credential interface {
	// SELECT id from @@table WHERE id IN @ids;
	FindByIdsIn(ids ...string) ([]string, error)
}
