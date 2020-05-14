package nomadclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)


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

func TestNewClientHTTP(t *testing.T) {
	Convey("Given valid input values", t, func(){
		c, err := NewClient("http://localhost:4646", "", false)
		Convey("When a new Client is created", func(){
			Convey("Then a client is returned with no errors", func(){
				So(err, ShouldBeNil)
				So(c, ShouldNotBeNil)
			})
		})
	})
}

func TestCreateTLSConfigCert(t *testing.T) {
	Convey("Given valid input values", t, func(){
		c, err := createTlsConfig("", true)
		Convey("When a new config is created", func(){
			Convey("Then config is returned with no errors", func(){
				So(err, ShouldBeNil)
				So(c, ShouldNotBeNil)
			})
		})
	})
}

func TestCreateTLSConfigCertFails(t *testing.T) {
	Convey("Given invalid input values", t, func(){
		c, err := createTlsConfig("/does/not/exist", false)
		Convey("When a new config is created", func(){
			Convey("Then no config is returned with errors", func(){
				So(err, ShouldNotBeNil)
				So(c, ShouldBeNil)
				So(err.Error(), ShouldStartWith, "open /does/not/exist: no such")
			})
		})
	})
}

func TestCreateTLSConfigBadCertFails(t *testing.T) {
	Convey("Given invalid input values", t, func(){
		c, err := createTlsConfig("testdata/testBadCertFile", false)
		Convey("When a new config is created", func(){
			Convey("Then no config is returned with errors", func(){
				So(err, ShouldNotBeNil)
				So(c, ShouldBeNil)
				So(err.Error(), ShouldStartWith, "failed to append ca cert to pool")
			})
		})
	})
}