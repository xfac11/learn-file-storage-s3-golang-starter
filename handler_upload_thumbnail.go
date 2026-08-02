package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if !(mediaType == "image/jpeg" || mediaType == "image/png") {
		respondWithError(w, http.StatusBadRequest, "Not supported media type. Use a jpeg or png", err)
		return
	}
	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to find the video", err)
		return
	}

	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "The user is not the owner of the video", err)
		return
	}

	fileExtension := strings.Split(mediaType, "/")[1]
	filePath := filepath.Join(cfg.assetsRoot, videoID.String())
	newFile, err := os.Create(filePath + "." + fileExtension)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not save the image on the servers file system", err)
		return
	}

	_, err = io.Copy(newFile, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not copy the file from the request to the server", err)
		return
	}

	thumbnailURL := fmt.Sprintf("http://localhost:%s/assets/%s.%s", cfg.port, videoID, fileExtension)
	video.ThumbnailURL = &thumbnailURL

	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update the video on the database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
