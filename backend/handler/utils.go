package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"
)

const dateThreshold = time.Minute

func checkEncodedValue(v string) error {
	if len(v) == 0 {
		return errors.New("field empty")
	}
	_, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return errors.New("value is not base64 encoded")
	}

	return nil
}

func timeParse(r *http.Request) (time.Time, error) {
	dateString := r.Header.Get("Date")
	if len(dateString) == 0 {
		return time.Time{}, errors.New("date header not found")
	}
	return time.Parse(time.RFC3339, r.Header.Get("Date"))
}

func checkValidDate(date time.Time) bool {
	now := time.Now()
	return now.After(date) && now.Before(date.Add(dateThreshold))
}

func timeCheck(r *http.Request) error {
	date, err := timeParse(r)
	if err != nil {
		return err
	}
	if !checkValidDate(date) {
		return errors.New("date threshold exceeded")
	}

	return nil
}
