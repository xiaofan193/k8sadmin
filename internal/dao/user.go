package dao

import (
	"context"
	"errors"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"github.com/xiaofan193/k8sadmin/internal/cache"
	"github.com/xiaofan193/k8sadmin/internal/database"
	"github.com/xiaofan193/k8sadmin/internal/model"
)

var _ UserDao = (*userDao)(nil)

// UserDao defining the dao interface
type UserDao interface {
	Create(ctx context.Context, table *model.User) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.User, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.User) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.User) error
}

type userDao struct {
	db    *gorm.DB
	cache cache.UserCache     // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewUserDao creating the dao interface
func NewUserDao(db *gorm.DB, xCache cache.UserCache) UserDao {
	if xCache == nil {
		return &userDao{db: db}
	}
	return &userDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *userDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new user, insert the record and the id value is written back to the table
func (d *userDao) Create(ctx context.Context, table *model.User) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a user by id
func (d *userDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a user by id, support partial update
func (d *userDao) UpdateByID(ctx context.Context, table *model.User) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *userDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.User) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Name != "" {
		update["name"] = table.Name
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a user by id
func (d *userDao) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	// no cache
	if d.cache == nil {
		record := &model.User{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.User{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				if errors.Is(err, database.ErrRecordNotFound) {
					// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, id, table, cache.UserExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.User)
		if !ok {
			return nil, database.ErrRecordNotFound
		}
		return table, nil
	}

	if d.cache.IsPlaceholderErr(err) {
		return nil, database.ErrRecordNotFound
	}

	return nil, err
}

// GetByColumns get a paginated list of users by custom conditions.
// For more details, please refer to https://go-sponge.com/component/custom-page-query.html
func (d *userDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.User, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.UserColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.User{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.User{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *userDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.User) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *userDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *userDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.User) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
