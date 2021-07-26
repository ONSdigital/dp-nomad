package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func getMockNomad(expectRequest http.Request, mockedHTTPResponse MockedHTTPResponse) (*nomad.Client, *httptest.Server) {
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

	return client, ts
}

func TestCheckOk(t *testing.T) {
	Convey("Given that Nomad is available", t, func() {
		mockedNomad, ts := getMockNomad(
			http.Request{Method: "GET"},
			MockedHTTPResponse{StatusCode: 200, Body: `{"status": "OK"}`},
		)
		defer ts.Close()

		Convey("When a check is performed", func() {
			checkState := health.NewCheckState(nomad.ServiceName)

			Convey("Then the Checker returns no error and status ok", func() {
				err := mockedNomad.Checker(ctx, checkState)
				So(err, ShouldBeNil)
				So(checkState.Status(), ShouldEqual, health.StatusOK)
				So(checkState.StatusCode(), ShouldEqual, 200)
			})

		})
	})
}

func TestCheckUnexpectedStatus(t *testing.T) {
	Convey("Given that Nomad is available, but returns an unexpected status code", t, func() {
		mockedNomad, ts := getMockNomad(
			http.Request{Method: "GET"},
			MockedHTTPResponse{StatusCode: 204, Body: `{"status": "OK"}`},
		)
		defer ts.Close()

		Convey("When a check is performed", func() {
			checkState := health.NewCheckState(nomad.ServiceName)

			Convey("Then the Checker returns an error, and the status is critical", func() {
				err := mockedNomad.Checker(ctx, checkState)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unexpected return code")
				So(checkState.Status(), ShouldEqual, health.StatusCritical)
				So(checkState.StatusCode(), ShouldEqual, 204)
			})

		})
	})
}

func TestCheckFail(t *testing.T) {
	Convey("Given that Nomad is available, but returns an error", t, func() {
		mockedNomad, ts := getMockNomad(
			http.Request{Method: "GET"},
			MockedHTTPResponse{StatusCode: 500, Body: `{"status": "CRITICAL"}`},
		)
		defer ts.Close()

		Convey("When a check is performed", func() {
			checkState := health.NewCheckState(nomad.ServiceName)

			Convey("Then the Checker returns an error,a nd the status is critical", func() {
				err := mockedNomad.Checker(ctx, checkState)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid response")
				So(checkState.Status(), ShouldEqual, health.StatusCritical)
				So(checkState.StatusCode(), ShouldEqual, 500)
			})

		})
	})
}

func TestCheckServerUnavailable(t *testing.T) {
	Convey("Given that Nomad is unavailable", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
		}))
		defer ts.Close()
		mockedNomad, _ := nomad.NewClient(ts.URL, "testdata/testCertFile", false)

		mockedNomad.Client.SetMaxRetries(0)
		mockedNomad.Client.SetTimeout(1 * time.Second)

		Convey("When a check is performed", func() {
			checkState := health.NewCheckState(nomad.ServiceName)

			Convey("Then the Checker returns an error, that timeout has been exceeded", func() {
				err := mockedNomad.Checker(ctx, checkState)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "Client.Timeout exceeded")
				So(checkState.Status(), ShouldEqual, health.StatusCritical)
				So(checkState.StatusCode(), ShouldEqual, 0)
			})

		})
	})
}
