package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-api/common"
	"io/ioutil"
	"net/http"
)

type reqData struct {
	Aid string `json:"aid"`
}

func GetVideoInfo(ctx *gin.Context) {
	url := "https://api.bilibili.com/x/web-interface/view?aid=170001"

	res, err := http.Get(url)
	if err != nil {
		common.Abort(common.ErrBind, err, ctx)
		return
	}

	rbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.Abort(common.ErrBind, err, ctx)
		return
	}

	var successData json.RawMessage
	d := json.NewDecoder(bytes.NewBuffer(rbody))
	err = d.Decode(&successData)

	common.SuccessReturn(successData, ctx)
}
