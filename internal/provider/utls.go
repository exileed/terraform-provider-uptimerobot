package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/exileed/uptimerobotapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

const timeoutMinutes = 60000

func retryTime(retryFunc func() error, minutes int) error {
	wait := 2
	return resource.RetryContext(context.Background(), time.Duration(minutes)*time.Minute, func() *resource.RetryError {
		err := retryFunc()

		if err == nil {
			return nil
		}

		log.Printf(fmt.Sprintf("[DEBUG] Error response %s", err.Error()))

		rand.Seed(time.Now().UnixNano())
		randomNumberMilliseconds := rand.Intn(1001)

		log.Printf(fmt.Sprintf("[DEBUG] Retrying server error code %d", randomNumberMilliseconds))

		if apiErr, ok := err.(uptimerobotapi.APIError); ok && (apiErr.StatusCode == 500 || apiErr.StatusCode == 502 || apiErr.StatusCode == 503) {
			timeSleep := time.Duration(wait)*time.Second + time.Duration(randomNumberMilliseconds)

			log.Printf(fmt.Sprintf("[DEBUG] Retrying server error code. Sllep: %s", timeSleep))
			time.Sleep(timeSleep)
			wait = wait * 2
			return resource.RetryableError(apiErr)
		}

		if apiErr, ok := err.(uptimerobotapi.APIError); ok && (apiErr.StatusCode == 409 || apiErr.StatusCode == 429) {
			timeSleep := time.Duration(wait)*time.Second + time.Duration(randomNumberMilliseconds)

			log.Printf("[DEBUG] Retrying quota/server error code...")
			log.Printf(fmt.Sprintf("[DEBUG] Retrying quota/server error code %s", timeSleep))

			time.Sleep(timeSleep)
			wait = wait * 2
			return resource.RetryableError(apiErr)
		}

		// Deal with the broken API
		if strings.Contains(fmt.Sprintf("%s", err), "Invalid Input: Bad request for \"") && strings.Contains(fmt.Sprintf("%s", err), "\"code\":400") {
			log.Printf("[DEBUG] Retrying invalid response from API")
			return resource.RetryableError(err)
		}
		if strings.Contains(fmt.Sprintf("%s", err), "Service unavailable. Please try again") {
			log.Printf("[DEBUG] Retrying service unavailable from API")
			return resource.RetryableError(err)
		}
		if strings.Contains(fmt.Sprintf("%s", err), "Eventual consistency. Please try again") {
			log.Printf("[DEBUG] Retrying due to eventual consistency")
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func mapKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		panic(errors.New("not a map"))
	}

	keys := v.MapKeys()
	s := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		s[i] = keys[i].String()
	}

	return s
}

func intToString(m map[string]int, value int) string {
	for k, v := range m {
		if int(v) == value {
			return k
		}
	}
	return ""
}
