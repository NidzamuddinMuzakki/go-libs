package payload

import (
	"go-cimb-lib/log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	ResponseHeader struct {
		ResponseCode    string `example:"2000100" json:"response_code"`
		ResponseMessage string `example:"Success" json:"response_message"`
	}
	Response struct {
		ResponseHeader
		Meta Meta        `json:"meta"`
		Data interface{} `json:"data,omitempty"`
	}
	ResponsePagination struct {
		Response
		Pagination Pagination `json:"pagination"`
	}
	Pagination struct {
		PageNum       int   `json:"page_num"`
		RecordPerPage int   `json:"record_per_page"`
		TotalPage     int   `json:"total_page"`
		TotalRecord   int64 `json:"total_record"`
	}
	Meta struct {
		DeviceID    string `example:"EA7583CD-A667-48BC-B806-42ECB2B48606" json:"device_id,omitempty"`
		TraceID     string `example:"97125121-ea32-4ee0-8706-5b7375e83e94" json:"trace_id"`
		DebugParam  string `example:"" json:"debug_param,omitempty"`
		Description string `example:"Success" json:"description,omitempty"`
	}
)

func GenerateCommonResponse(res *Response) {
	switch res.ResponseCode {
	case "200":
		res.ResponseCode = "200"
		res.ResponseMessage = "Success"
		res.Meta.Description = "Success"
	case "400":
		res.ResponseCode = "400"
		res.ResponseMessage = "Bad Request"
	case "401":
		res.ResponseCode = "401"
		res.ResponseMessage = "Unathorize"
	case "500":
		res.ResponseCode = "500"
		res.ResponseMessage = "Internal Server Err"
	}

}

func Json(c *gin.Context, res *Response, err *error, logging log.ILogging) {
	var (
		url = c.Request.Host + c.Request.URL.Path
	)
	errData := *err
	if errData != nil {
		res.Meta.DebugParam = errData.Error()
	}
	res.Meta.TraceID = c.GetHeader("Request-ID")
	GenerateCommonResponse(res)
	logging.Http(res.Meta.TraceID, "Controller Info", url, c.Request.Method, c.Request.Header, c.Request.Body, &res)
	code, _ := strconv.Atoi(res.ResponseCode)
	c.JSON(code, res)
	logging.Info(res.Meta.TraceID, "Closing", c.Request.URL.Path)
}
