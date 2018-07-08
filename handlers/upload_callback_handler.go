package handlers

import (
	"github.com/akihokurino/crypto-assistant-server/applications"
	"net/http"
	"google.golang.org/appengine"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"google.golang.org/appengine/blobstore"
)

type UploadCallbackHandler interface {
	CallbackUserIcon(w http.ResponseWriter, r *http.Request)
}

type uploadCallbackHandler struct {
	uploadApplication applications.UploadApplication
}

func NewUploadCallbackHandler(uploadApplication applications.UploadApplication) UploadCallbackHandler {
	return &uploadCallbackHandler{
		uploadApplication: uploadApplication,
	}
}

func (h *uploadCallbackHandler) CallbackUserIcon(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	uid := models.UserID(r.URL.Query().Get("uid"))

	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.uploadApplication.UploadUserIcon(ctx, uid, blobs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
