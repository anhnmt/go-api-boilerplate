// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package gormgen

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	sessionentity "github.com/anhnmt/go-api-boilerplate/internal/service/session/entity"
)

func newSession(db *gorm.DB, opts ...gen.DOOption) session {
	_session := session{}

	_session.sessionDo.UseDB(db, opts...)
	_session.sessionDo.UseModel(&sessionentity.Session{})

	tableName := _session.sessionDo.TableName()
	_session.ALL = field.NewAsterisk(tableName)
	_session.ID = field.NewString(tableName, "id")
	_session.CreatedAt = field.NewTime(tableName, "created_at")
	_session.UpdatedAt = field.NewTime(tableName, "updated_at")
	_session.DeviceID = field.NewString(tableName, "device_id")
	_session.Token = field.NewString(tableName, "token")
	_session.IsRevoked = field.NewBool(tableName, "is_revoked")
	_session.ExpiredAt = field.NewTime(tableName, "expired_at")

	_session.fillFieldMap()

	return _session
}

type session struct {
	sessionDo

	ALL       field.Asterisk
	ID        field.String
	CreatedAt field.Time
	UpdatedAt field.Time
	DeviceID  field.String
	Token     field.String
	IsRevoked field.Bool
	ExpiredAt field.Time

	fieldMap map[string]field.Expr
}

func (s session) Table(newTableName string) *session {
	s.sessionDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s session) As(alias string) *session {
	s.sessionDo.DO = *(s.sessionDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *session) updateTableName(table string) *session {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewString(table, "id")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")
	s.DeviceID = field.NewString(table, "device_id")
	s.Token = field.NewString(table, "token")
	s.IsRevoked = field.NewBool(table, "is_revoked")
	s.ExpiredAt = field.NewTime(table, "expired_at")

	s.fillFieldMap()

	return s
}

func (s *session) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *session) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 7)
	s.fieldMap["id"] = s.ID
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
	s.fieldMap["device_id"] = s.DeviceID
	s.fieldMap["token"] = s.Token
	s.fieldMap["is_revoked"] = s.IsRevoked
	s.fieldMap["expired_at"] = s.ExpiredAt
}

func (s session) clone(db *gorm.DB) session {
	s.sessionDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s session) replaceDB(db *gorm.DB) session {
	s.sessionDo.ReplaceDB(db)
	return s
}

type sessionDo struct{ gen.DO }

type ISessionDo interface {
	gen.SubQuery
	Debug() ISessionDo
	WithContext(ctx context.Context) ISessionDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ISessionDo
	WriteDB() ISessionDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ISessionDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ISessionDo
	Not(conds ...gen.Condition) ISessionDo
	Or(conds ...gen.Condition) ISessionDo
	Select(conds ...field.Expr) ISessionDo
	Where(conds ...gen.Condition) ISessionDo
	Order(conds ...field.Expr) ISessionDo
	Distinct(cols ...field.Expr) ISessionDo
	Omit(cols ...field.Expr) ISessionDo
	Join(table schema.Tabler, on ...field.Expr) ISessionDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ISessionDo
	RightJoin(table schema.Tabler, on ...field.Expr) ISessionDo
	Group(cols ...field.Expr) ISessionDo
	Having(conds ...gen.Condition) ISessionDo
	Limit(limit int) ISessionDo
	Offset(offset int) ISessionDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ISessionDo
	Unscoped() ISessionDo
	Create(values ...*sessionentity.Session) error
	CreateInBatches(values []*sessionentity.Session, batchSize int) error
	Save(values ...*sessionentity.Session) error
	First() (*sessionentity.Session, error)
	Take() (*sessionentity.Session, error)
	Last() (*sessionentity.Session, error)
	Find() ([]*sessionentity.Session, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*sessionentity.Session, err error)
	FindInBatches(result *[]*sessionentity.Session, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*sessionentity.Session) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ISessionDo
	Assign(attrs ...field.AssignExpr) ISessionDo
	Joins(fields ...field.RelationField) ISessionDo
	Preload(fields ...field.RelationField) ISessionDo
	FirstOrInit() (*sessionentity.Session, error)
	FirstOrCreate() (*sessionentity.Session, error)
	FindByPage(offset int, limit int) (result []*sessionentity.Session, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ISessionDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	FindByIdsIn(ids []string) (result []string, err error)
}

// SELECT id from @@table WHERE id IN @ids;
func (s sessionDo) FindByIdsIn(ids []string) (result []string, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, ids)
	generateSQL.WriteString("SELECT id from sessions WHERE id IN ?; ")

	var executeSQL *gorm.DB
	executeSQL = s.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (s sessionDo) Debug() ISessionDo {
	return s.withDO(s.DO.Debug())
}

func (s sessionDo) WithContext(ctx context.Context) ISessionDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s sessionDo) ReadDB() ISessionDo {
	return s.Clauses(dbresolver.Read)
}

func (s sessionDo) WriteDB() ISessionDo {
	return s.Clauses(dbresolver.Write)
}

func (s sessionDo) Session(config *gorm.Session) ISessionDo {
	return s.withDO(s.DO.Session(config))
}

func (s sessionDo) Clauses(conds ...clause.Expression) ISessionDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s sessionDo) Returning(value interface{}, columns ...string) ISessionDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s sessionDo) Not(conds ...gen.Condition) ISessionDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s sessionDo) Or(conds ...gen.Condition) ISessionDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s sessionDo) Select(conds ...field.Expr) ISessionDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s sessionDo) Where(conds ...gen.Condition) ISessionDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s sessionDo) Order(conds ...field.Expr) ISessionDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s sessionDo) Distinct(cols ...field.Expr) ISessionDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s sessionDo) Omit(cols ...field.Expr) ISessionDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s sessionDo) Join(table schema.Tabler, on ...field.Expr) ISessionDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s sessionDo) LeftJoin(table schema.Tabler, on ...field.Expr) ISessionDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s sessionDo) RightJoin(table schema.Tabler, on ...field.Expr) ISessionDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s sessionDo) Group(cols ...field.Expr) ISessionDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s sessionDo) Having(conds ...gen.Condition) ISessionDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s sessionDo) Limit(limit int) ISessionDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s sessionDo) Offset(offset int) ISessionDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s sessionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ISessionDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s sessionDo) Unscoped() ISessionDo {
	return s.withDO(s.DO.Unscoped())
}

func (s sessionDo) Create(values ...*sessionentity.Session) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s sessionDo) CreateInBatches(values []*sessionentity.Session, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s sessionDo) Save(values ...*sessionentity.Session) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s sessionDo) First() (*sessionentity.Session, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*sessionentity.Session), nil
	}
}

func (s sessionDo) Take() (*sessionentity.Session, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*sessionentity.Session), nil
	}
}

func (s sessionDo) Last() (*sessionentity.Session, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*sessionentity.Session), nil
	}
}

func (s sessionDo) Find() ([]*sessionentity.Session, error) {
	result, err := s.DO.Find()
	return result.([]*sessionentity.Session), err
}

func (s sessionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*sessionentity.Session, err error) {
	buf := make([]*sessionentity.Session, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s sessionDo) FindInBatches(result *[]*sessionentity.Session, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s sessionDo) Attrs(attrs ...field.AssignExpr) ISessionDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s sessionDo) Assign(attrs ...field.AssignExpr) ISessionDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s sessionDo) Joins(fields ...field.RelationField) ISessionDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s sessionDo) Preload(fields ...field.RelationField) ISessionDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s sessionDo) FirstOrInit() (*sessionentity.Session, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*sessionentity.Session), nil
	}
}

func (s sessionDo) FirstOrCreate() (*sessionentity.Session, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*sessionentity.Session), nil
	}
}

func (s sessionDo) FindByPage(offset int, limit int) (result []*sessionentity.Session, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s sessionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s sessionDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s sessionDo) Delete(models ...*sessionentity.Session) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *sessionDo) withDO(do gen.Dao) *sessionDo {
	s.DO = *do.(*gen.DO)
	return s
}