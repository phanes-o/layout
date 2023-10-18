package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/phanes-o/lib/otel/trace"
	"phanes/errors"
)

type ErrorHandler interface {
	Unexportable(c *gin.Context)
	Exportable(c *gin.Context)
}

type httpErrorHandler struct {
	err error
}

func NewHttpErrorHandler(e error) ErrorHandler  {
	return &httpErrorHandler{
		err: e,
	}
}


func (h *httpErrorHandler) Unexportable(c *gin.Context) {
	var (
		traceID = trace.TraceIDFromContext(c.Request.Context())
	)
	// Request params error return 400
	var errs validator.ValidationErrors
	if errors.As(h.err, &errs) {
		c.JSON(400, gin.H{
			"trace_id": traceID,
			"code":     errors.BadRequest,
			"msg":      RemoveTopStruct(errs.Translate(translate)),
		})
	}else {
	// Some can't exportable error return 500
		c.JSON(500, gin.H{
			"trace_id": traceID,
			"code":     500,
			"msg":      "Server Internal Error",
		})
	}
}

func (h *httpErrorHandler) Exportable(c *gin.Context) {
	var (
		err = h.err
		traceID = trace.TraceIDFromContext(c.Request.Context())
	)

	errCode := errors.GetType(err)
	if errCode == 1000 {
		c.JSON(http.StatusUnauthorized, nil)
	//  Common errors handle
	} else if errCode > 1000 && errCode < 2000 {
		c.JSON(http.StatusOK, gin.H{
			"trace_id": traceID,
			"code":     errCode,
			"msg":      err.Error(),
		})
	// Customer error handle here
	} else if errCode >= 2000 && errCode < 3000 {
		c.JSON(http.StatusOK, gin.H{
			"trace_id": traceID,
			"code":     errCode,
			"msg":      err.Error(),
		})
	}
}



