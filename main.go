/*
Opinionated ory/dockertest wrapper

Performs the following quality of life improvements:

• Build a Dockerfile

• Give it a unique name

• Expose its ports (like -P)

• Register it for purge if the test cleanup handlers run

• Wait for its health check to pass

• Configure the container to be auto-removed (like --rm)

    IMPORTANT: You should use resource.Expire(seconds) to ensure that your
               returned resource is cleaned up in case of the test being killed.

    IMPORTANT: You _must_ configure a health check in your Dockerfile.

To use this for testing with a dockerfile which runs a webserver
located at test/Dockerfile:

    func TestSudoRun(t *testing.T) {
    	resource := dockertesting.RunDockerfile(t, "test/Dockerfile")
    	resource.Expire(300)
      hp := resource.GetHostPort("80/tcp")
      _, err := http.Get(hp)
      if err != nil {
        t.Fatalf("Failed to get from container: %v", err)
      }
      // when this closes, the container will be purged.
    }

*/
package dockertesting

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// Takes a testing.T and the filename of the Dockerfile.
//
// Returns a docker resource
func RunDockerfile(t *testing.T, filename string) *dockertest.Resource {
	t.Helper()
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	dockerUuid := uuid.New().String()
	// Build and run the given Dockerfile
	resource, err := pool.BuildAndRunWithOptions(
		filename,
		&dockertest.RunOptions{
			Name: "go-docker-testing-image-" + dockerUuid,
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
		})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	t.Logf("Created container.  waiting for it to start")

	t.Cleanup(func() {
		pool.Purge(resource)
	})

	for {
		res, ok := pool.ContainerByName(resource.Container.Name)
		if ok && res.Container.State.Health.Status == "healthy" {
			t.Logf("Container started")
			return res
		}
	}
}
