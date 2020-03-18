package common

import (
	"github.com/juetun/study/app-dashboard/lib/utilsold/hashid"
	"github.com/speps/go-hashids"
)

var (
	ZHashId *hashids.HashID
)

func PluginsHashId() (err error) {

	hd := new(hashid.HashIdParams)
	salt := hd.SetHashIdSalt("i must add a salt what is only for me")
	hdLength := hd.SetHashIdLength(8)
	zHashId, err := hd.HashIdInit(hdLength, salt)
	if err != nil {
		return
	}
	ZHashId = zHashId
	return
}
