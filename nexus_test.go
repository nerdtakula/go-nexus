package nexus

import (
	"log"
	"os"
	"testing"
)

var client Client

func TestMain(m *testing.M) {
	client, err := New("http://localhost:8081/service/rest/v1")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	client.SetBasicAuth("admin", "admin123")

	code := m.Run()
	os.Exit(code)
}

// // Comment this out to run tests agaist an existing instance
// func TestMain(m *testing.M) {
// 	// Uses a sensible default on windows (tcp/http) and linux/osx (socket)
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	// Pulls an image, creates a container based on it and runs it
// 	resource, err := pool.Run("sonatype/nexus3", "latest", nil)
// 	if err != nil {
// 		log.Fatalf("Could not start resource: %s", err)
// 	}

// 	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
// 	err = pool.Retry(func() error {
// 		//
// 		client = New(fmt.Sprintf("localhost:%s/service/rest/v1", resource.GetPort("8081/tcp"))).SetBasicAuth("admin", "admin123")
// 		return client.Ping()
// 	})
// 	if err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	code := m.Run()

// 	// You can't defer this because os.Exit doesn't care for defer
// 	if err := pool.Purge(resource); err != nil {
// 		log.Fatalf("Could not purge resource: %s", err)
// 	}

// 	os.Exit(code)
// }
