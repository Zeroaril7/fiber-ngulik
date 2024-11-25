package utils

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func ConvertString(v interface{}) string {
	result := ""
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		result = v.(string)
	case int:
		result = strconv.Itoa(v.(int))
	case int64:
		result = strconv.FormatInt(v.(int64), 10)
	case bool:
		result = strconv.FormatBool(v.(bool))
	case float64:
		result = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case []uint8:
		result = string(v.([]uint8))
	default:
		resultJSON, err := json.Marshal(v)
		if err == nil {
			result = strings.Trim(string(resultJSON), "\"")
		} else {
			log.Println("Error on lib/converter ConvertString() ", err)
		}
	}

	return result
}

func ConvertInt(v interface{}) int {
	result := int(0)
	switch v.(type) {
	case string:
		str := strings.TrimSpace(v.(string))
		result, _ = strconv.Atoi(str)
	case int:
		result = int(v.(int))
	case int64:
		result = int(v.(int64))
	case float64:
		result = int(v.(float64))
	case []byte:
		result, _ = strconv.Atoi(string(v.([]byte)))
	default:
		result = int(0)
	}

	return result
}

func ConvertInt64(v interface{}) int64 {
	result := int64(0)
	switch v.(type) {
	case string:
		str := strings.TrimSpace(v.(string))
		result, _ = strconv.ParseInt(str, 10, 64)
	case int:
		result = int64(v.(int))
	case int64:
		result = int64(v.(int64))
	case float64:
		result = int64(v.(float64))
	case []byte:
		result, _ = strconv.ParseInt(string(v.([]byte)), 10, 64)
	default:
		result = int64(0)
	}

	return result
}

func GetLocalTime() time.Time {
	return time.Now().Local()
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
