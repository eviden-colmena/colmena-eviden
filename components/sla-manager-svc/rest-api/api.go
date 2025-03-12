/*
  COLMENA-DESCRIPTION-SERVICE
  Copyright © 2024 EVIDEN
  
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  
  http://www.apache.org/licenses/LICENSE-2.0
  
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  
  This work has been implemented within the context of COLMENA project.
*/
/*
Copyright © 2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This work has been implemented within the context of COLMENA project.
*/
package restapi

import (
	"colmena/sla-management-svc/app/assessment"
	"colmena/sla-management-svc/app/assessment/monitor"
	"colmena/sla-management-svc/app/common/logs"
	"colmena/sla-management-svc/app/model"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
	"time"
)

// path used in logs
const pathLOG string = "SLA > REST-API > "

// App is a main application "object", to be built by main and testmain
type App struct {
	Router      *gin.Engine
	Repository  model.IRepository
	Monitor     monitor.MonitoringAdapter
	Port        string
	SslEnabled  bool
	SslCertPath string
	SslKeyPath  string
	externalIDs bool
	validator   model.Validator
}

func New(config assessment.Config, repository model.IRepository, validator model.Validator, monitor monitor.MonitoringAdapter) (App, error) {
	a := App{
		Repository: repository,
		Monitor:    monitor,
		/*Port:        config.GetString(portPropertyName),
		SslEnabled:  config.GetBool(enableSslPropertyName),
		SslCertPath: config.GetString(sslCertPathPropertyName),
		SslKeyPath:  config.GetString(sslKeyPathPropertyName),
		externalIDs: config.GetBool(utils.ExternalIDsPropertyName),*/
		validator: validator,
	}

	//a.initialize(repository)

	return a, nil
}

/*
InitializeRESTAPI initialization function
*/
func (a *App) InitializeRESTAPI() {
	logs.GetLogger().Info(pathLOG + "[InitializeRESTAPI] Initializing REST API Server ...")

	// router
	a.Router = gin.Default()

	// CORS https://github.com/gin-contrib/cors
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	a.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "DELETE", "PATCH", "OPTIONS", "GET"}, // "POST, GET, DELETE, PUT, OPTIONS"
		AllowHeaders:     []string{"Origin"},                                           // "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		/*AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},*/
		MaxAge: 12 * time.Hour,
	}))

	// ping
	a.Router.GET("/api/v1/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	api := a.Router.Group("/api")
	{
		// public methods
		public := api.Group("/v1")
		{
			// default - configuration - status
			public.GET("/", responseNotImplementedFunc)

			// sla
			public.POST("/sla", a.CreateSLAv2)
			public.GET("/sla", a.GetSLAs)
			public.GET("/sla/:id", a.GetSLA)
			public.DELETE("/sla/:id", a.DeleteSLA)

			// query metrics
			// api/v1/query?metric=<METRIC>&path=<PATH>
			public.GET("/query", a.Query)

			// force violation
			public.POST("/sla/violation/:fid", responseNotImplementedFunc)
		}
	}

	// swagger
	// a.Router.Static("/swaggerui", "./rest-api/swaggerui/")
	// TLS
	//r.RunTLS(":443", "resources/sec/server.crt", "resources/sec/server.key")

	/////////////////////////////////////////////////////////////////
	// start server: 8333 (default)

	srv := &http.Server{
		Addr:    ":8080", //":" + strconv.Itoa(cfg.Config.E2COPort), //":8080",
		Handler: a.Router,
	}

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		// HTTPS
		/*
			if err := srv.ListenAndServeTLS("resources/sec/server.crt", "resources/sec/server.key"); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Error(pathLOG + "[InitializeRESTAPIv2] ListenAndServeTLS Error: ", err)
			}
		*/
		///*
		logs.GetLogger().Info(pathLOG + "[InitializeRESTAPI] Running and listening on port 8080 ...")
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logs.GetLogger().Error(pathLOG+"[InitializeRESTAPIv2] ListenAndServe Error: ", err)
		}
		//*/
	}()

	/////////////////////////////////////////////////////////////////
	// stop server:
	/*
		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-quit

		log.Info("..............................................................................")

		// close DB
		//db.CloseConnection()

		log.Info(pathLOG + "[InitializeRESTAPI] Shutting down server ...")

		// The context is used to inform the server it has 5 seconds to finish the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("Shutdown error: " + err.Error())
		}

		time.Sleep(1 * time.Second)
		log.Info(pathLOG + "[InitializeRESTAPI] Terminated")
	*/
}

///////////////////////////////////////////////////////////////////////////////

/*
responseNotImplementedFunc Default Function for not implemented calls
*/
func responseNotImplementedFunc(c *gin.Context) {
	logs.GetLogger().Warn(pathLOG + "Function not implemented")

	c.JSON(200, gin.H{
		"Resp":    "ok",
		"Method":  "responseNotImplementedFunc",
		"Message": "Function not implemented"})
}

/*
responseError response function
*/
func responseError(c *gin.Context, method string, message string) {
	logs.GetLogger().Error(pathLOG + "[" + method + "] " + message)

	c.JSON(500, gin.H{
		"Resp":    "error",
		"Method":  method,
		"Message": message})
}

/*
responseOk response function
*/
func responseOk(c *gin.Context, method string, message string, code int, obj interface{}) {
	logs.GetLogger().Info(pathLOG + "[" + method + "] " + message)

	c.JSON(code, gin.H{
		"Resp":     "ok",
		"Method":   method,
		"Message":  message,
		"Response": obj})
}

// create
func create(c *gin.Context, m string, decode func() error, create func() (model.Identity, error)) {
	errDec := decode()
	if errDec != nil {
		responseError(c, m, "Error decoding input: "+errDec.Error())
	}
	/* check errors */
	created, err := create()
	if err != nil {
		responseError(c, m, "Error creating object: "+err.Error())
	} else {
		responseOk(c, m, "Object created", http.StatusCreated, created)
	}
}

// getAll
func getAll(c *gin.Context, m string, f func() (interface{}, error)) {
	list, err := f()
	if err != nil {
		responseError(c, m, "Error getting objects: "+err.Error())
	} else {
		responseOk(c, m, "Objects found", http.StatusOK, list)
	}
}

// get
func get(c *gin.Context, m string, f func(string) (interface{}, error)) {
	id := c.Param("id")

	res, err := f(id)
	if err != nil {
		responseError(c, m, "Error getting object: "+err.Error())
	} else {
		responseOk(c, m, "Object found", http.StatusOK, res)
	}
}

// delete
func delete(c *gin.Context, m string, f func(string) error) {
	id := c.Param("id")

	err := f(id)
	if err != nil {
		responseError(c, m, "Error deleting object: "+err.Error())
	} else {
		responseOk(c, m, "Object deleted", http.StatusOK, nil)
	}
}

// violation
func violation(c *gin.Context, m string, f func(string) error) {
	fid := c.Param("fid")

	err := f(fid)
	if err != nil {
		responseError(c, m, "Error creating violation object: "+err.Error())
	} else {
		responseOk(c, m, "Violation created", http.StatusOK, nil)
	}
}

// violationsList
func violationsList(c *gin.Context, m string, f func() error) {
	err := f()
	if err != nil {
		responseError(c, m, "Error creating violations list object: "+err.Error())
	} else {
		responseOk(c, m, "Violations list created", http.StatusOK, nil)
	}
}

///////////////////////////////////////////////////////////////////////////////
// API METHODS:

/*
GetMetric gets metric's value: "api/v1/query?metric=<METRIC>&path=<PATH>"
Example:

	curl http://localhost:8080/api/v1/query?metric=colmena_metric1&path=/tests/planta01/habitacion01
	curl http://localhost:8080/api/v1/query?metric=colmena_metric1&path=~/tests/planta01/.*
*/
func (a *App) Query(c *gin.Context) {

	metric := c.Query("metric")
	path := c.DefaultQuery("path", "")

	if metric == "" {
		responseError(c, "Query", "Metric not defined: Query format: 'api/v1/query?metric=METRIC&path=PATH'")
	}

	res, err := a.Monitor.Query(metric, path) // return model.MetricValue
	if err != nil {
		responseError(c, "Query", "Error getting metric: "+err.Error())
	} else {
		responseOk(c, "Query", "Metric retrieved", http.StatusOK, res)
	}
}

/*
CreateSLA creates a SLA passed by REST params
*/
func (a *App) CreateSLA(c *gin.Context) {
	var qos model.SLA

	create(c, "CreateSLA",
		func() error {
			return c.ShouldBindJSON(&qos)
		},
		func() (model.Identity, error) {
			return a.Repository.CreateSLA(&qos)
		})
}

/*
CreateSLAv2 creates multiple SLAs passed by REST params
*/
func (a *App) CreateSLAv2(c *gin.Context) {
	slas, err := model.InputSLAModelToSLAModel(c)
	if err != nil {
		responseError(c, "CreateSLA", "Error decoding input: "+err.Error())
	}

	anyError := false
	var resSlas []model.SLA
	var resError []error

	for _, sla := range slas {
		m, e := a.Repository.CreateSLA(&sla)

		if e != nil {
			anyError = true
			resError = append(resError, e)
		} else {
			resSlas = append(resSlas, *m)
		}
	}

	if anyError {
		str1 := fmt.Sprintf("%#v", resError)
		str2 := fmt.Sprintf("%#v", resSlas)
		responseError(c, "CreateSLA", "Error creating SLA(s): [slas = "+str2+"]; [errors = "+str1+"]")
	} else {
		responseOk(c, "CreateSLA", "SLA(s) created", http.StatusOK, nil)
	}

}

/*
GetSLAs return all SLAs in db
*/
func (a *App) GetSLAs(c *gin.Context) {
	getAll(c, "GetSLAs", func() (interface{}, error) {
		return a.Repository.GetSLAs()
	})
}

/*
GetSLA gets a QoS Definition by REST ID
*/
func (a *App) GetSLA(c *gin.Context) {
	get(c, "GetSLA", func(id string) (interface{}, error) {
		return a.Repository.GetSLA(id)
	})
}

/*
DeleteSLA deletes a SLA
*/
func (a *App) DeleteSLA(c *gin.Context) {
	delete(c, "DeleteSLA", func(id string) error {
		return a.Repository.DeleteSLA(id)
	})
}
