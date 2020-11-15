package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/harshitsinghai/felix/generatemd5"
	"github.com/harshitsinghai/felix/models"
	query "github.com/harshitsinghai/felix/queries"
	"github.com/harshitsinghai/felix/utils"
)

// GenerateShortURL adds URL to DB and gives back shortened string
func GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var record models.Record
	var requestBody models.GenerateRequest

	postBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(postBody, &requestBody)

	md5Hash := generatemd5.GenerateHash(requestBody.OriginalURL)
	record.ShortURL = md5Hash
	record.OriginalURL = requestBody.OriginalURL
	record.CreatedAt = time.Now()

	if requestBody.Expiry {
		dur, err := strconv.ParseUint(requestBody.ExpiresAfter, 0, 64)
		if err != nil {
			log.Println("Some error occured when converting to uint64")
		}

		expireLinkType := requestBody.ExpiryDateType
		record.ExpiresAt = utils.SetExpiryDate(expireLinkType, dur)
		log.Println("Time generated by GO", record.ExpiresAt)
	}

	exists, originalURL := query.FetchAlreadyExists(record.ShortURL)

	if exists {
		w.WriteHeader(http.StatusOK)
		responseMap := map[string]interface{}{"alreadyExists": true, "originalUrl": originalURL, "shortUrl": md5Hash}
		response, _ := json.Marshal(responseMap)
		w.Write(response)
		return
	}

	id, err := query.InsertURL(&record)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Some error occured"))
		return
	}

	record.AlreadyExists = false
	record.ID = id
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(record)
	w.Write(response)
}
