package resource

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor"
	"github.com/qor/qor/roles"
	"github.com/qor/qor/utils"
)

func (res *Resource) findOneHandler(result interface{}, metaValues *MetaValues, context *qor.Context) error {
	if res.HasPermission(roles.Read, context) {
		var primaryKey string
		if metaValues == nil {
			primaryKey = context.ResourceID
		} else if id := metaValues.Get(res.PrimaryFieldName()); id != nil {
			primaryKey = utils.ToString(id.Value)
		}

		if primaryKey != "" {
			if metaValues != nil {
				if destroy := metaValues.Get("_destroy"); destroy != nil {
					if fmt.Sprint(destroy.Value) != "0" {
						context.GetDB().Delete(result, primaryKey)
						return ErrProcessorSkipLeft
					}
				}
			}
			return context.GetDB().First(result, primaryKey).Error
		}
		return errors.New("failed to find")
	} else {
		return roles.ErrPermissionDenied
	}
}

func (res *Resource) findManyHandler(result interface{}, context *qor.Context) error {
	if res.HasPermission(roles.Read, context) {
		return context.GetDB().Set("gorm:order_by_primary_key", "DESC").Find(result).Error
	} else {
		return roles.ErrPermissionDenied
	}
}

func (res *Resource) saveHandler(result interface{}, context *qor.Context) error {
	if (context.GetDB().NewScope(result).PrimaryKeyZero() &&
		res.HasPermission(roles.Create, context)) || // has create permission
		res.HasPermission(roles.Update, context) { // has update permission
		return context.GetDB().Save(result).Error
	} else {
		return roles.ErrPermissionDenied
	}
}

func (res *Resource) deleteHandler(result interface{}, context *qor.Context) error {
	if res.HasPermission(roles.Delete, context) {
		if !context.GetDB().First(result, context.ResourceID).RecordNotFound() {
			return context.GetDB().Delete(result).Error
		} else {
			return gorm.RecordNotFound
		}
	} else {
		return roles.ErrPermissionDenied
	}
}

func (res *Resource) CallFindOne(result interface{}, metaValues *MetaValues, context *qor.Context) error {
	return res.FindOneHandler(result, metaValues, context)
}

func (res *Resource) CallFindMany(result interface{}, context *qor.Context) error {
	return res.FindManyHandler(result, context)
}

func (res *Resource) CallSave(result interface{}, context *qor.Context) error {
	return res.SaveHandler(result, context)
}

func (res *Resource) CallDelete(result interface{}, context *qor.Context) error {
	return res.DeleteHandler(result, context)
}
