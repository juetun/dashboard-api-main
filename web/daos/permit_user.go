// Package daos /**
package daos

type DaoPermitUser interface {
	UpdateDataByUserHIds(data map[string]interface{}, userHIds ...string) (err error)
}
