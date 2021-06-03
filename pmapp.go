/*
 * Copyright: Oleg Borodin <onborodin@gmail.com>
 */

package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "sync"
    "time"
    "errors"
    "encoding/json"
    "io/ioutil"
    "bytes"

    "app/pmconfig"
    "app/pmlog"

    "github.com/gorilla/websocket"
    "github.com/xeipuuv/gojsonschema"
    "github.com/gin-gonic/gin"
)

func main() {
    var err error
    app := NewApp()

    err = app.AppStart()
    if err != nil {
        pmlog.LogError("application error:", err)
        os.Exit(1)
    }
}

const (
    loopPeriod      time.Duration   = 1000  // msec
    appAlivePeriod  int64           = 60    // sec
)

type Application struct {
    config      *pmconfig.Config
    context     context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
}

func NewApp() *Application {
    var app Application
    app.context, app.cancel = context.WithCancel(context.Background())
    app.config  = pmconfig.NewConfig()
    return &app
}

func (this *Application) AppStart() error {
    var err error

    err = this.config.Setup()
    if err != nil {
        return err
    }

    err = this.startWeb()
    if err != nil {
        return err
    }

    err = this.startLoop()
    if err != nil {
        return err
    }
    return err
}
//
// startLoop()
//
func (this *Application) startLoop() error {
    var err error
    loopFunc := func() {
        pmlog.LogInfo("application loop started")
        timer := time.NewTicker(loopPeriod * time.Millisecond)
        for nextTime := range timer.C {
            switch {
                case nextTime.Unix() % appAlivePeriod == 0:
                    pmlog.LogInfo("application is alive")
            }
        }
    }
    loopFunc()
    return err
}
//
// startWeb()
//
func (this *Application) startWeb() error {
    var err error
    
    gin.DisableConsoleColor()
    gin.SetMode(gin.ReleaseMode)

    router := gin.New()
    logFormatter := func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s %s %s %s %s %d %d %s\n",
            param.TimeStamp.Format(time.RFC3339),
            param.ClientIP,
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.BodySize,
            param.Latency,
        )
    }
    router.Use(gin.LoggerWithFormatter(logFormatter))
    router.Use(gin.Recovery())

    err = os.MkdirAll(filepath.Dir(this.config.AccessLogPath), 0755)
    if err != nil {
        return err
    }
    accessLogFile, err := os.OpenFile(this.config.AccessLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
    if err != nil {
      return err
    } 
    gin.DefaultWriter = io.MultiWriter(accessLogFile, os.Stdout)

    helloCont := NewController()

    // Classic Controller-Model design 
    router.GET("/hello", helloCont.Hello)

    // JSON RPC 2.0 over HTTP
    router.POST("/fc", webFuncResolver)

    // JSON RPC 2.0 over Web Socket
    router.GET("/ws", func(ctx *gin.Context) {
        wsResolveHandler(ctx.Writer, ctx.Request)
    })

    // Subscription over Web Socket
    router.GET("/sc", func(ctx *gin.Context) {
        wsSubsrHandler(ctx.Writer, ctx.Request)
    })
    
    router.NoRoute(func(ctx *gin.Context) {
        SendError(ctx, emptyRequestId, errorMethodNotFound, errors.New("route not found"))
    })

    go router.Run(this.config.GetListenParam())    
    return err
}

func wsResolveHandler(writer http.ResponseWriter, request *http.Request) {
    var wsUpgrader = websocket.Upgrader{
        ReadBufferSize:  1024*8,
        WriteBufferSize: 1024*8,
    }
    conn, err := wsUpgrader.Upgrade(writer, request, nil)
    defer conn.Close()
    
    if err != nil {
        pmlog.LogError("failed to set websocket upgrade:", err)
        return
    }

    for {
        _, reqPayload, err := conn.ReadMessage()
        if err != nil {
            pmlog.LogDebug("ws read error:", err)
            break
        }
        // Function resolver
        resPayload, err := wsFuncResolver(reqPayload)

        err = conn.WriteMessage(websocket.TextMessage, resPayload)
        if err != nil {
            pmlog.LogDebug("ws write error:", err)
            break
        }
    }
}

func wsSubsrHandler(writer http.ResponseWriter, request *http.Request) {
    var wsUpgrader = websocket.Upgrader{
        ReadBufferSize:  1024*8,
        WriteBufferSize: 1024*8,
    }
    conn, err := wsUpgrader.Upgrade(writer, request, nil)
    defer conn.Close()
    
    if err != nil {
        pmlog.LogError("failed to set websocket upgrade:", err)
        return
    }

    for {
        type Status struct {
            Timestamp  time.Time
        }
        status := Status{ Timestamp: time.Now() }
        resPayload, _ := wsMakeResult(emptyRequestId, status)

        err = conn.WriteMessage(websocket.TextMessage, resPayload)
        if err != nil {
            pmlog.LogDebug("ws write error:", err)
            break
        }
    }
}


const (
    mimeApplicationJson string  = "application/json"
    jsonrpcId           string  = "2.0"
    funcAddName         string  = "add"
    emptyRequestId      string  = ""
    
    errorDefaultError   int     = 1
    errorParseError     int     = -32700 	//Invalid JSON was received by the server.
    errorInvalidRequest int     = -32600 	//The JSON sent is not a valid Request object.
    errorMethodNotFound int     = -32601 	//The method does not exist / is not available.
    errorInvalidParams  int     = -32602 	//Invalid method parameter(s).
    errorInternalError  int     = -32603 	 //Internal JSON-RPC error.
)

type BaseRequest struct {
    JSONRPC     string          `json:"jsonrpc"`
    Id          string          `json:"id"`
    Method    string            `json:"method"`
}

type  ResponseError struct  {
    Code    int                 `json:"code,omitempty"`
    Message string              `json:"message,omitempty"`
} 

type Response struct {
    JSONRPC     string          `json:"jsonrpc"`
    Id          string          `json:"id,omitempty"`
    Error       interface{}     `json:"error,omitempty"`
    Result      interface{}     `json:"result,omitempty"`
}

func wsFuncResolver(requestBody []byte) ([]byte, error) {
    var err error
    var request BaseRequest

    err = json.Unmarshal(requestBody, &request)
    if err != nil {
        return wsMakeError(request.Id, errorParseError, err)
    }

    switch request.Method {
        case funcAddName:
            return wsAdd(requestBody)
    }
    return wsMakeError(request.Id, errorMethodNotFound, errors.New("method not found"))
}


//
// wFuncResolver()
//
func webFuncResolver(ctx *gin.Context) {
    var err error
    var request BaseRequest

    var requestBytes []byte
    if ctx.Request.Body != nil {
        requestBytes, _ = ioutil.ReadAll(ctx.Request.Body)
        ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBytes))
    }

    err = ctx.ShouldBind(&request)
    if err != nil {
        responseBytes, _ := wsMakeError(request.Id, errorParseError, err)
        ctx.Data(http.StatusOK, mimeApplicationJson, responseBytes)
        return
    }
    
    switch request.Method {
        case funcAddName:
            responseBytes, _ := wsAdd(requestBytes)
            ctx.Data(http.StatusOK, mimeApplicationJson, responseBytes)
            return
    }
    response, _ := wsMakeError(request.Id, errorMethodNotFound, err)
    ctx.Data(http.StatusOK, mimeApplicationJson, response)
    return
}


const (
    addReqSchemaV4 string = `
        {
            "$schema": "http://json-schema.org/draft-04/schema",
            "type": "object",
            "required": [
                "method",
                "params"
            ],
            "properties": {
                "method": {
                    "type": "string"
                },
                "params": {
                    "type": "object",
                    "required": [
                        "first",
                        "second"
                    ],
                    "properties": {
                        "first": {
                            "type": "integer"
                        },
                        "second": {
                            "type": "integer"
                        }
                    }
                }
            }
        }`

    addReqSchemaV4min string = `{"$schema":"http://json-schema.org/draft-04/schema","type":"object","required":["method","params"],"properties":{"method":{"type":"string"},"params":{"type":"object","required":["first","second"],"properties":{"first":{"type":"integer"},"second":{"type":"integer"}}}}}`
)

type AddRequest struct {
    BaseRequest
    Params struct {
        First   int64     `json:"first"`
        Second  int64     `json:"second"`
    } `json:"params"`
}
type AddResult = int64

func wsAdd(requestBytes []byte) ([]byte, error) {
    var err error
    var request AddRequest

    // Map request body to structure
    err = json.Unmarshal(requestBytes, &request)
    if err != nil {
        return wsMakeError(request.Id, errorParseError, err)
    }

    // Validate request schema 
    err = jsonValidator(addReqSchemaV4min, string(requestBytes))
    if err != nil {
        return wsMakeError(request.Id, errorParseError, err)
    }
    
    // Function body
    var result AddResult = request.Params.First + request.Params.Second
    // Make RPC response body
    return wsMakeResult(request.Id, result)
}

func jsonValidator(schema string, document string) error {
    var err error
    schemaLoader := gojsonschema.NewStringLoader(schema)
    documentLoader := gojsonschema.NewStringLoader(document)
    result, err := gojsonschema.Validate(schemaLoader, documentLoader)
    if err != nil {
        return err
    }
    if result != nil && !result.Valid() {
        var errorDescriptions string = "validation error:"
        for _, err := range result.Errors() {
            errorDescriptions = errorDescriptions + " " + err.Context().String() + ":" + err.Description()
        }
        return errors.New(errorDescriptions)
    }
    return err
}

func wsMakeError(reqId string, errorCode int, funcErr error) ([]byte, error) {
    var err error
    if funcErr == nil {
        funcErr = errors.New("undefined")
    }
    responseError := ResponseError{
        Code:       errorCode,
        Message:    fmt.Sprintf("%s", funcErr),
    }
    response := Response{
        JSONRPC:    jsonrpcId,
        Id:         reqId,
        Error:      responseError,
    }
    resBytes, err := json.Marshal(response)
    return resBytes, err
}

func wsMakeResult(reqId string, result interface{}) ([]byte, error) {
    var err error
    response := Response{
        JSONRPC:    jsonrpcId,
        Id:         reqId,
        Result:     result,
    }
    resBytes, err := json.Marshal(response)
    return resBytes, err
}

func wsMakeOk() ([]byte, error) {
    var err error
    response := Response{
        JSONRPC:    jsonrpcId,
    }
    resBytes, err := json.Marshal(response)
    return resBytes, err
}

//
// Classical Controller-Model 
//

//
// Controller
//
type Controller struct{
}

func (this *Controller) Hello(ctx *gin.Context) {
    model := NewModel()
    result, err := model.Hello()
    if err != nil {
        SendError(ctx, emptyRequestId, errorParseError, err)
    }
    SendResult(ctx, emptyRequestId, &result)
}

func NewController() *Controller {
    return &Controller{
    }
}
//
// Model
//
type Model struct {
}

type Hello struct {
    Message string     `db:"message" json:"message"`
}

func (this *Model) Hello() (*Hello, error) {
    var err error
    result := Hello{ Message: "Hello, World" }
    return &result, err
}

func NewModel() *Model {
    var model Model
    return &model
}
//
// Tools
//
func SendError(ctx *gin.Context, requestId string, errorCode int, err error) {
    if err == nil {
        err = errors.New("undefined")
    }
    responseError := ResponseError{
        Code:       errorCode,
        Message:    fmt.Sprintf("%s", err),
    }
    response := Response{
        JSONRPC:    jsonrpcId,
        Id:         requestId,
        Error:      responseError,
    }
    ctx.JSON(http.StatusOK, response)
}

func SendResult(ctx *gin.Context, requestId string, result interface{}) {
    response := Response{
        JSONRPC:    jsonrpcId,
        Id:         requestId,
        Result: result,
    }
    ctx.JSON(http.StatusOK, &response)
}

func SendOk(ctx *gin.Context) {
    response := Response{
        JSONRPC:    jsonrpcId,
    }
    ctx.JSON(http.StatusOK, response)
}

//func SendMessage(ctx *gin.Context, message string) {
//    response := Response{
//        Error: false,
//        Message: fmt.Sprintf("%s", message),
//    }
//    ctx.IndentedJSON(http.StatusOK, response)
//}

//func AbortContext(ctx *gin.Context, code int, err error) {
//    if err == nil {
//        err = errors.New("undefined")
//    }
//    response := Response{
//        Error: true,
//        Message: fmt.Sprintf("%s", err),
//    }
//    ctx.JSON(code, response)
//    ctx.Abort()
//}
//EOF
