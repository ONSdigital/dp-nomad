package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	nomad "github.com/ONSdigital/dp-nomad"
	. "github.com/smartystreets/goconvey/convey"
)

type MockedHTTPResponse struct {
	StatusCode int
	Body       string
}

const Name = "nomad"

var ctx = context.Background()

func getMockNomad(expectRequest http.Request, mockedHTTPResponse MockedHTTPResponse) *nomad.Client {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectRequest.Method {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unexpected HTTP method used"))
			return
		}
		w.WriteHeader(mockedHTTPResponse.StatusCode)
		fmt.Fprintln(w, mockedHTTPResponse.Body)
	}))

	client, _ := nomad.NewClient(ts.URL, "testdata/testCertFile", false)

	return client
}

func TestCheckOk(t *testing.T) {
	Convey("Given that the vault client is available", t, func() {
		//c, _ := nomad.NewClient(ts.URL, "", true)
		mockedNomad := getMockNomad(
			http.Request{Method: "GET"},
			MockedHTTPResponse{StatusCode: 200, Body: "{\"status\": \"OK\"}"})

		Convey("When a check is performed", func() {
			checkState := health.NewCheckState(nomad.ServiceName)

			Convey("Then the Checker returns no error", func() {
				//err := c.Checker(context.Background(), checkState)
				err := mockedNomad.Checker(ctx, checkState)
				So(err, ShouldBeNil)
				So(checkState.Status(), ShouldEqual, health.StatusOK)
				So(checkState.StatusCode(), ShouldEqual, 200)
			})

		})
	})
}

func TestCheckStatus204(t *testing.T) {
	Convey("Given that the vault client is available", t, func() {
		mockedNomad := getMockNomad(
			http.Request{Method: "GET"},
			MockedHTTPResponse{StatusCode: 204, Body: `{"status": "OK"}`})

		Convey("When a check is performed", func() {
			checkState := health.NewCheckState(nomad.ServiceName)

			Convey("Then the Checker returns no error", func() {
				err := mockedNomad.Checker(ctx, checkState)
				So(err, ShouldBeNil)
				So(checkState.StatusCode(), ShouldEqual, 204)
			})

		})
	})
}

func TestCheckFail(t *testing.T) {
	Convey("Given that the vault client is available", t, func() {
		mockedNomad := getMockNomad(
			http.Request{Method: "GET"},
			MockedHTTPResponse{StatusCode: 500, Body: `{"status": "CRITICAL"}`})

		Convey("When a check is performed", func() {
			checkState := health.NewCheckState(nomad.ServiceName)

			Convey("Then the Checker returns no error", func() {
				err := mockedNomad.Checker(ctx, checkState)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid response")
				So(checkState.StatusCode(), ShouldEqual, 0)
			})

		})
	})
}