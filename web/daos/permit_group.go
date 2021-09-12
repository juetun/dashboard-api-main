/**
* @Author:changjiang
* @Description:
* @File:permit_group
* @Version: 1.0.0
* @Date 2021/9/12 11:40 上午
 */
package daos

type DaoPermitGroup interface {

	DeleteUserGroupByUserId(userId ...string) (err error)

	DeleteUserGroupPermitByGroupId(ids ...string) (err error)

	DeleteUserGroupPermit(pathType string, menuId ...int) (err error)
}
