package nomadclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)


// NewClient(NomadEndpoint, NomadCACert string, NomadTLSSkipVerify bool)
func TestNewClient(t *testing.T) {
	Convey("Given acceptable input values", t, func(){
		c, err := NewClient("https://localhost:4646", "testdata/testCertFile", false)
		Convey("When a newClient is created", func() {
			Convey("Then a client is returned with no errors", func(){
				So(c, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNewClientFails(t *testing.T) {
	Convey("Given invalid input values", t, func(){
		c, err := NewClient("https://localhost:4646", "", false)
		Convey("When a new Client is created", func(){
			Convey("Then no client is returned with errors", func(){
				So(err, ShouldNotBeNil)
				So(c, ShouldBeNil)
				So(err.Error(), ShouldStartWith, "invalid configuration with https")
			})
		})
	})
}