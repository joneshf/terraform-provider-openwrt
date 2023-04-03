package acceptancetest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gotest.tools/v3/assert"
)

const (
	acceptanceTestDockerDockerfile = "acceptance-test.Dockerfile"
	acceptanceTestDockerName       = "joneshf/openwrt"
	acceptanceTestDockerHTTPPort   = "80/tcp"
	acceptanceTestDockerSSHPort    = "22/tcp"
	acceptanceTestDockerTag        = "acceptance-test"

	dockerContainerHealthy = "healthy"
)

var (
	acceptanceTestDockerCachedFilePath = filepath.Join(".cache", "docker", "acceptance-test.tar")
	acceptanceTestDockerImage          = fmt.Sprintf("%s:%s", acceptanceTestDockerName, acceptanceTestDockerTag)
)

// OpenWrtServer groups all the information needed to connect to the running OpenWrt server.
type OpenWrtServer struct {
	Hostname string
	HTTPPort uint16
	Password string
	Scheme   string
	SSHPort  uint16
	Username string
}

// LuCIRPCClient returns a [*lucirpc.Client] to interact with the [OpenWrtServer].
func (s OpenWrtServer) LuCIRPCClient(
	ctx context.Context,
	t *testing.T,
) *lucirpc.Client {
	t.Helper()

	client, err := lucirpc.NewClient(
		ctx,
		s.Scheme,
		s.Hostname,
		s.HTTPPort,
		s.Username,
		s.Password,
	)
	assert.NilError(t, err)
	return client
}

// ProviderBlock creates a stringified provider block for the OpenWrt provider.
func (s OpenWrtServer) ProviderBlock() string {
	return fmt.Sprintf(`
provider "openwrt" {
	hostname = %q
	password = %q
	port = %d
	username = %q
}
`,
		s.Hostname,
		s.Password,
		s.HTTPPort,
		s.Username,
	)
}

// RunOpenWrtServer constructs a running OpenWrt server.
// The hostname and port can be arbitrary,
// so they are returned to make it easier to interact with the OpenWrt server.
func RunOpenWrtServer(
	ctx context.Context,
	dockerPool dockertest.Pool,
	t *testing.T,
) OpenWrtServer {
	t.Helper()

	openWrt, err := dockerPool.Run(acceptanceTestDockerName, acceptanceTestDockerTag, []string{})
	assert.NilError(t, err)
	t.Cleanup(func() {
		err = openWrt.Close()
		assert.NilError(t, err)
	})
	err = openWrt.Expire(60)
	assert.NilError(t, err)
	err = dockerPool.Retry(checkHealth(ctx, dockerPool, *openWrt))
	assert.NilError(t, err)
	hostname := openWrt.GetBoundIP(acceptanceTestDockerHTTPPort)
	rawHTTPPort := openWrt.GetPort(acceptanceTestDockerHTTPPort)
	intHTTPPort, err := strconv.Atoi(rawHTTPPort)
	assert.NilError(t, err)
	httpPort := uint16(intHTTPPort)
	rawSSHPort := openWrt.GetPort(acceptanceTestDockerSSHPort)
	intSSHPort, err := strconv.Atoi(rawSSHPort)
	assert.NilError(t, err)
	sshPort := uint16(intSSHPort)

	return OpenWrtServer{
		Hostname: hostname,
		HTTPPort: httpPort,
		Password: "",
		Scheme:   "http",
		SSHPort:  sshPort,
		Username: "root",
	}
}

// Setup does a bit of setup so acceptance tests can run:
//  1. Connect to a running docker daemon.
//  2. Build and tag the image for acceptance tests.
//
// The tearDown function must be called after tests are finished.
func Setup(
	ctx context.Context,
) (tearDown func(), dockerPool *dockertest.Pool, err error) {
	var conn net.Conn
	tearDown = func() {
		if conn != nil {
			conn.Close()
		}
	}

	var socket string
	homeDir, err := os.UserHomeDir()
	if err == nil {
		log.Printf("Attempting to connect to colima socket")
		colimaSocket := filepath.Join(homeDir, ".colima", "docker.sock")
		conn, err := net.Dial("unix", colimaSocket)
		if err == nil {
			log.Printf("Connect to colima successful")
			socket = fmt.Sprintf("%s://%s", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
		} else {
			log.Printf("Could not connect to colima")
		}
	} else {
		log.Printf("Could not find user home directory, defaulting to docker socket")
	}

	log.Printf("Constructing docker pool on socket: %q", socket)
	dockerPool, err = dockertest.NewPool(socket)
	if err != nil {
		err = fmt.Errorf("could not construct docker pool: %w", err)
		return
	}

	log.Printf("Connecting to docker")
	err = dockerPool.Client.PingWithContext(ctx)
	if err != nil {
		err = fmt.Errorf("could not connect to docker: %w", err)
		return
	}

	log.Printf("Grabbing the top-level directory")
	gitRevParse := exec.Command("git", "rev-parse", "--show-toplevel")
	revParseOutput, err := gitRevParse.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("could not grab the top-level directory: %w", err)
		return
	}

	topLevelDirectory := bytes.TrimSpace(revParseOutput)
	err = loadImageCache(ctx, dockerPool, string(topLevelDirectory))
	if err != nil {
		err = fmt.Errorf("could load image cache: %w", err)
		return
	}

	log.Printf("Building acceptance test image")
	err = dockerPool.Client.BuildImage(docker.BuildImageOptions{
		CacheFrom: []string{
			acceptanceTestDockerImage,
		},
		ContextDir:   string(topLevelDirectory),
		Dockerfile:   acceptanceTestDockerDockerfile,
		Name:         acceptanceTestDockerImage,
		OutputStream: os.Stdout,
	})
	if err != nil {
		err = fmt.Errorf("could not build acceptance test image: %w", err)
		return
	}

	err = saveImageCache(ctx, dockerPool, string(topLevelDirectory))
	if err != nil {
		err = fmt.Errorf("could save image cache: %w", err)
		return
	}

	return
}

// TerraformSteps removes a bit of boilerplate around setting up the [resource.TestStep]s.
func TerraformSteps(
	t *testing.T,
	testStep resource.TestStep,
	testSteps ...resource.TestStep,
) {
	t.Helper()

	allTestSteps := append([]resource.TestStep{testStep}, testSteps...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"openwrt": providerserver.NewProtocol6WithError(openwrt.New("test", os.LookupEnv)),
		},
		Steps: allTestSteps,
	})

}

func checkHealth(
	ctx context.Context,
	dockerPool dockertest.Pool,
	resource dockertest.Resource,
) func() error {
	return func() error {
		container, err := dockerPool.Client.InspectContainerWithContext(resource.Container.ID, ctx)
		if err != nil {
			return err
		}

		status := container.State.Health.Status
		if status != dockerContainerHealthy {
			return fmt.Errorf(status)
		}

		return nil
	}
}

func loadImageCache(
	ctx context.Context,
	dockerPool *dockertest.Pool,
	topLevelDirectory string,
) error {
	log.Println("Attempting to load image cache")
	cachedImageFilePath := filepath.Join(topLevelDirectory, acceptanceTestDockerCachedFilePath)
	cachedImageFile, err := os.Open(cachedImageFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Image file does not exist; ignoring")
			return nil
		}

		return fmt.Errorf("problem opening cached image file: %w", err)
	}

	defer cachedImageFile.Close()
	err = dockerPool.Client.LoadImage(docker.LoadImageOptions{
		InputStream: cachedImageFile,
		Context:     ctx,
	})
	if err != nil {
		return fmt.Errorf("problem loading image: %w", err)
	}

	log.Println("Loaded cached image")
	return nil
}

func saveImageCache(
	ctx context.Context,
	dockerPool *dockertest.Pool,
	topLevelDirectory string,
) error {
	log.Println("Saving image to cache")
	cachedImageFilePath := filepath.Join(topLevelDirectory, acceptanceTestDockerCachedFilePath)
	err := os.MkdirAll(filepath.Dir(cachedImageFilePath), 0755)
	if err != nil {
		return fmt.Errorf("problem creating cache directory: %w", err)
	}

	cachedImageFile, err := os.Create(cachedImageFilePath)
	if err != nil {
		return fmt.Errorf("problem creating cached image file: %w", err)
	}

	defer cachedImageFile.Close()
	err = dockerPool.Client.ExportImage(docker.ExportImageOptions{
		Context:      ctx,
		Name:         acceptanceTestDockerImage,
		OutputStream: cachedImageFile,
	})
	if err != nil {
		return fmt.Errorf("problem exporting image: %w", err)
	}

	log.Println("Saved image to cache")
	return nil
}
