package gtimeentry

import (
	"github.com/dougEfresh/gtoggl/test"
	"os"
	"testing"
	"time"
)

var _, debugMode = os.LookupEnv("GTOGGL_TEST_DEBUG")

func togglClient() *TimeEntryClient {
	tu := &gtest.TestUtil{Debug: debugMode}
	client := tu.MockClient()
	ws, err := NewClient(client)
	if err != nil {
		panic(err)
	}
	return ws
}

/*
func TestUserCreate(t *testing.T) {
	tClient := togglClient()
	nc, err := tClient.Create("signup@blah.com", "StrongPasswrod", "UTC")
	if err != nil {
		t.Fatal(err)
	}

	if nc.Id != 3 {
		t.Fatal("!= 3")
	}

	if nc.Email != "signup@blah.com" {
		t.Fatal("!= signup@blah.com")
	}

	if nc.ApiToken != "808lolae4eab897cce9729a53642124effe" {
		t.Fatal("!= 808lolae4eab897cce9729a53642124effe")
	}
}

func TestUserUpdate(t *testing.T) {
	tClient := togglClient()
	c := &User{Id: 1, FullName: "John Swift", Email: "newemail@swift.com"}
	nc, err := tClient.Update(c)
	if err != nil {
		t.Fatal(err)
	}

	if nc.Email != "newemail@swift.com" {
		t.Fatal("!= newemail@swift.com")
	}
}

func TestUserReset(t *testing.T) {
	tClient := togglClient()
	token, err := tClient.ResetToken()
	if err != nil {
		t.Fatal(err)
	}

	if token != "123456789" {
		t.Fatal("!= 123456789")
	}
}

*/

func TestTimeEntryDelete(t *testing.T) {
	tClient := togglClient()
	err := tClient.Delete(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTimeEntryList(t *testing.T) {
	tClient := togglClient()
	te, err := tClient.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(te) < 1 {
		t.Fatal("<1")
	}

}

func TestTimeEntryCreate(t *testing.T) {
	tClient := togglClient()

	te := &TimeEntry{}
	te.Billable = false
	te.Duration = 1200
	te.Pid = 123
	te.Wid = 777
	te.Description = "Meeting with possible clients"
	te.Tags = []string{"billed"}

	nTe, err := tClient.Create(te)

	if err != nil {
		t.Fatal(err)
	}
	if nTe.Id != 3 {
		t.Error("!= 3")
	}
}

func TestTimeEntryUpdate(t *testing.T) {
	tClient := togglClient()
	te, err := tClient.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	te.Description = "new"
	nTe, err := tClient.Update(te)
	if err != nil {
		t.Fatal(err)
	}
	if nTe.Description != "new" {
		t.Error("!= new")
	}
}

func TestTimeEntryGet(t *testing.T) {
	tClient := togglClient()

	timeentry, err := tClient.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if timeentry.Id != 1 {
		t.Error("!= 1")
	}

	st, err := time.Parse(time.RFC3339, "2013-02-27T01:24:00+00:00")

	if err != nil {
		t.Fatal(err)
	}
	diff := st.Sub(timeentry.Start)
	if diff != 0 {
		t.Errorf("!= %s", diff)
	}
	st, err = time.Parse(time.RFC3339, "2013-02-27T07:24:00+00:00")
	diff = st.Sub(timeentry.Stop)
	if diff != 0 {
		t.Errorf("!= %s", diff)
	}

	/*
		if timeentry.FullName != "John Swift" {
			t.Error("!= John Swift:  " + timeentry.FullName)
		}

		if timeentry.ApiToken != "1971800d4d82861d8f2c1651fea4d212" {
			t.Error("!= J1971800d4d82861d8f2c1651fea4d212:  " + timeentry.ApiToken)
		}

		if timeentry.Email != "johnt@swift.com" {
			t.Error("!= johnt@swift.com" + timeentry.Email)
		}
	*/
}

func BenchmarkClientTransport_Get(b *testing.B) {
	b.ReportAllocs()
	tClient := togglClient()
	for i := 0; i < b.N; i++ {
		tClient.Get(1)
	}
}
