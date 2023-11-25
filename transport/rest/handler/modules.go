package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github/yogabagas/join-app/config"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/pkg/storage/minio"
	"github/yogabagas/join-app/shared/util"
	"github/yogabagas/join-app/transport/rest/handler/response"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CreateModules handler
// @Summary Create New Modules
// @Description New Modules Registration
// @Tags Modules
// @Produce json
// @Param users body service.CreateModulesReq true "Request Create Modules"
// @Success 200 {object} response.JSONResponse().APIStatusCreated()
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/modules/create [POST]
func (h *HandlerImpl) CreateModules(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodPost {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	var req service.CreateModulesReq
	req.Name = r.FormValue("name")
	req.Description = r.FormValue("description")
	req.ModuleMaterial = service.ParseRequestModuleMaterial(r.FormValue("module_materials"))

	filename, err := h.ParseFileUpload(r, "file", "storage")
	req.File = filename

	if err != nil {
		res.SetError(response.ErrBadRequest).Send(w)
		return
	}

	userData := new(util.UserData)
	userData = userData.GetUserData(r)

	if err := h.Controller.ModulesController.CreateModules(r.Context(), req, userData); err != nil {
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(w)
		return
	}

	res.APIStatusCreated().Send(w)

}

func (h *HandlerImpl) ParseFileUpload(r *http.Request, keyFormFile string, folder string) (filename string, err error) {
	if err := r.ParseMultipartForm(1024); err != nil {
		return "", err
	}

	uploadedFile, handler, err := r.FormFile(keyFormFile)
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filename = handler.Filename
	fileLocation := filepath.Join(dir, folder, filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return "", err
	}

	// upload file to minio storage
	pathFile := "modules/files/" + filename
	_, err = minio.NewMinio().
		WithContext(r.Context()).
		SetBucket(config.GlobalCfg.Storage.Minio.Bucket).
		MakeBucket().
		UploadObject("modules/files/"+filename, fileLocation, util.GetFileContentType(targetFile))
	if err != nil {
		fmt.Println("upload to minio error:", err.Error())
	}

	// remove file on host
	err = os.Remove(fileLocation)
	if err != nil {
		fmt.Println(err)
	}

	return pathFile, err
}

// DownloadFromStorage handler
// @Summary Get File from storage
// @Description Download file from storage
// @Tags Modules
// @Produce json
// @Success 200 {object} response.JSONResponse().APIStatusCreated()
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /modules/download/file [GET]
func (h *HandlerImpl) DownloadFromStorage(w http.ResponseWriter, r *http.Request) {
	res := response.NewJSONResponse()

	objectName := r.URL.Query().Get("object_name")
	if objectName == "" {
		res.APIStatusBadRequest().SetMessage("param object_name is required").Send(w)
	}

	extenstion := ""
	objectNames := strings.Split(objectName, ".")
	if len(objectNames) > 1 {
		extenstion = objectNames[len(objectNames)-1]
	}

	bucketExist := minio.NewMinio().
		SetBucket(config.GlobalCfg.Storage.Minio.Bucket).
		BucketExist()
	if bucketExist {
		preSign := minio.NewMinio().
			WithContext(r.Context()).
			SetBucket(config.GlobalCfg.Storage.Minio.Bucket).
			SignUrl(objectName)

		resp, err := http.Get(fmt.Sprint(preSign))
		// resp, err := http.Get("https://golangcode.com/logo.svg")
		if err != nil {
			res.APIStatusErrorUnknown().Send(w)
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			res.APIStatusErrorUnknown().Send(w)
		}

		w.WriteHeader(http.StatusOK)
		//w.Header().Set("Content-Disposition", "attachment; filename="+filename[3])
		w.Header().Set("Content-Type", "image/"+extenstion)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data) //the memory take up 1.2~1.7G
	} else {
		res.APIStatusNotFound().SetMessage("Bucket not found").Send(w)
	}
}

// GetModules handler
// @Summary Get All Modules
// @Description New Modules Registration
// @Tags Modules
// @Produce json
// @Param users body service.CreateModulesReq true "Request Create Modules"
// @Success 200 {object} response.JSONResponse().APIStatusCreated()
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/modules [GET]
func (h *HandlerImpl) GetModulesWithPagination(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodGet {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	var req service.GetModulesWithPaginationReq

	if name := r.URL.Query().Get("name"); name != "" {
		req.Name = name
	}

	var limitToInt int
	if limit := r.URL.Query().Get("limit"); limit != "" {
		limitToInt, _ = strconv.Atoi(limit)
	}

	if limitToInt <= 0 {
		limitToInt = 10
	}
	req.Limit = limitToInt

	var pageToInt int
	if page := r.URL.Query().Get("page"); page != "" {
		pageToInt, _ = strconv.Atoi(page)
	}

	if pageToInt <= 0 {
		pageToInt = 1
	}
	req.Page = pageToInt

	resp, err := h.Controller.ModulesController.GetModulesWithPagination(r.Context(), req)
	if err != nil {
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(w)
		return
	}

	res.SetData(resp).Send(w)

}

// GetModules handler
// @Summary Get All Courses
// @Description New Resources Registration
// @Tags Resources
// @Produce json
// @Param users body service.CreateResourcesReq true "Request Create Resources"
// @Success 200 {object} response.JSONResponse().APIStatusCreated()
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/modules [GET]
func (h *HandlerImpl) UpdateCourses(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodPost {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	var req service.UpdateModulesReq
	req.Name = r.FormValue("name")
	req.Description = r.FormValue("description")
	req.ModuleMaterial = service.ParseRequestModuleMaterial(r.FormValue("module_materials"))

	_, handler, _ := r.FormFile("new_file")
	if handler.Filename != "" {
		filename, err := h.ParseFileUpload(r, "new_file", "storage")
		req.File = filename

		if err != nil {
			res.SetError(response.ErrBadRequest).Send(w)
			return
		}
	}

	userData := new(util.UserData)
	userData = userData.GetUserData(r)

	err := h.Controller.ModulesController.UpdateModules(r.Context(), req, userData)
	if err != nil {
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(w)
		return
	}

	res.APIStatusCreated().Send(w)

}

// GetModules handler
// @Summary Get All Courses
// @Description New Resources Registration
// @Tags Resources
// @Produce json
// @Param users body service.CreateResourcesReq true "Request Create Resources"
// @Success 200 {object} response.JSONResponse().APIStatusCreated()
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/modules [GET]
func (h *HandlerImpl) DeleteCourse(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodDelete {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}
	vars := mux.Vars(r)
	userData := new(util.UserData)
	userData = userData.GetUserData(r)

	err := h.Controller.ModulesController.DeleteModules(r.Context(), vars["uid"], *userData)
	if err != nil {
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(w)
		return
	}

	res.APIStatusCreated().Send(w)

}
