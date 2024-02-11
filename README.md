```go
package example

import (
	"fmt"
	"net/http"
	"testing"

	"goception"
)

type TestSuite struct {
	Actor *goception.Actor
}

func (suite *TestSuite) Before(t *testing.T) {
	suite.Actor.SetBaseURL("...").HaveHeader("Authorization", "Bearer token...")
	//start transaction db
	//load fixtures
	//prepare db
	//etc
}

func (suite *TestSuite) After(t *testing.T) {
	// Perform actions after each test
	//fixture clean all tables
	//rollback transactions for testing db
	//etc
}

func (suite *TestSuite) TestExample1(t *testing.T) {
	var responseField string

	suite.Actor.SetBaseURL("http://example.com").
		HaveHeader("Custom-Header", "value").
		SendGet("/path").
		ExpectResponse(http.StatusOK).
		ParseFieldFromJSONPath("path.to.field", &responseField)

	fmt.Println("Response field:", responseField)
}

func (suite *TestSuite) TestExample2(t *testing.T) {
	var responseField string

	suite.Actor.SetBaseURL("http://example.com").
		HaveHeader("Custom-Header", "value").
		SendGet("/path").
		ExpectResponse(http.StatusOK).
		ParseFieldFromJSONPath("path.to.field", &responseField)

	fmt.Println("Response field:", responseField)
}

func (suite *TestSuite) TestExample3(t *testing.T) {
	var responseField string

	suite.Actor.SetBaseURL("http://example.com").
		HaveHeader("Custom-Header", "value").
		SendGet("/path").
		ExpectResponse(http.StatusOK).
		ParseFieldFromJSONPath("path.to.field", &responseField)

	fmt.Println("Response field:", responseField)
}

func TestRun(t *testing.T) {
	goception.RunSuites(t, &TestSuite{
		Actor: goception.NewActor(t),
	})
}

```