package client_test

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-net/http/httptest"
	"testing"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	nomad "github.com/ONSdigital/dp-nomad"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckOk(t *testing.T) {
	ts := httptest.NewTestServer(200)
	defer ts.Close()
	Convey("Given that the vault client is available", t, func() {
		c, _ := nomad.NewClient(ts.URL, "", true)
		Convey("When a check is preformed", func() {
			checkState := health.NewCheckState(nomad.ServiceName)

			Convey("Then the Checker returns no error", func() {
				err := c.Checker(context.Background(), checkState)
				So(err, ShouldBeNil)
				So(checkState.Status(), ShouldEqual, health.StatusOK)
				So(checkState.StatusCode(), ShouldEqual, 200)
			})

		})
	})
}

func TestCheckGetError(t *testing.T) {
	Convey("Given that the nomad client is unavailable", t, func() {
		c, _ := nomad.NewClient("https://localhost:4645", "", true)
		Convey("When a check is preformed", func() {
			checkState := health.NewCheckState(nomad.ServiceName)

			Convey("Then the Checker returns an error", func() {
				err := c.Checker(context.Background(), checkState)
				fmt.Print("Error ", err)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "connection refused")
			})

		})
	})
}