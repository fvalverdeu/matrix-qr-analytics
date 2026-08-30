//go:build integration

package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"qr-service-go/internal/clients"
	"qr-service-go/internal/handlers"
	"qr-service-go/internal/models"
	"qr-service-go/internal/services"
)

const (
	integrationTolerance = 1e-9
	diagonalEpsilon      = 1e-10
)

type synchronizedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *synchronizedBuffer) appendLine(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf.WriteString(line)
	b.buf.WriteByte('\n')
}

func (b *synchronizedBuffer) appendRaw(text string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf.WriteString(text)
}

func (b *synchronizedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

type nodeTestServer struct {
	cmd      *exec.Cmd
	baseURL  string
	stdout   *synchronizedBuffer
	stderr   *synchronizedBuffer
	waitDone chan error
}

func TestGoNodeIntegration_WideMatrixFullFlow(t *testing.T) {
	node := startNodeTestServer(t)

	app := buildRealGoApp(node.baseURL, 1500*time.Millisecond)

	input := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}

	requestBody, err := json.Marshal(models.QRRequest{Matrix: input})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/qr", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fiber app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := readBody(resp)
		t.Fatalf("unexpected status: got %d want %d body=%s", resp.StatusCode, http.StatusOK, string(bodyBytes))
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("expected JSON response Content-Type, got %q", contentType)
	}

	var body models.QRResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	assertShape(t, "Q", body.Q, 2, 2)
	assertShape(t, "R", body.R, 2, 3)
	assertAllFinite(t, "Q", body.Q)
	assertAllFinite(t, "R", body.R)

	assertReconstruction(t, input, body.Q, body.R, integrationTolerance)
	assertOrthogonality(t, body.Q, integrationTolerance)

	expectedStats := recomputeStatistics(body.Q, body.R)
	assertApproxEqual(t, "statistics.max", body.Statistics.Max, expectedStats.Max, integrationTolerance)
	assertApproxEqual(t, "statistics.min", body.Statistics.Min, expectedStats.Min, integrationTolerance)
	assertApproxEqual(t, "statistics.sum", body.Statistics.Sum, expectedStats.Sum, integrationTolerance)
	assertApproxEqual(t, "statistics.average", body.Statistics.Average, expectedStats.Average, integrationTolerance)
	if body.Statistics.HasDiagonalMatrix != expectedStats.HasDiagonalMatrix {
		t.Fatalf("unexpected statistics.hasDiagonalMatrix: got %v want %v", body.Statistics.HasDiagonalMatrix, expectedStats.HasDiagonalMatrix)
	}
}

func TestGoNodeIntegration_NodeUnavailableReturns502(t *testing.T) {
	unreachableURL := reserveAndReleaseLocalURL(t)
	app := buildRealGoApp(unreachableURL, 250*time.Millisecond)

	requestBody, err := json.Marshal(models.QRRequest{Matrix: [][]float64{{1, 2, 3}, {4, 5, 6}}})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/qr", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fiber app.Test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		bodyBytes, _ := readBody(resp)
		t.Fatalf("unexpected status: got %d want %d body=%s", resp.StatusCode, http.StatusBadGateway, string(bodyBytes))
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("expected JSON response Content-Type, got %q", contentType)
	}

	var body models.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON error response: %v", err)
	}

	expected := models.ErrorResponse{
		Error: models.APIError{
			Code:    "STATISTICS_UNAVAILABLE",
			Message: "Statistics service is unavailable",
		},
	}
	if body != expected {
		t.Fatalf("unexpected error response: got %+v want %+v", body, expected)
	}
}

func startNodeTestServer(t *testing.T) *nodeTestServer {
	t.Helper()

	nodePath, err := exec.LookPath("node")
	if err != nil {
		t.Fatalf("node executable was not found in PATH; Node.js is required for Go-Node integration tests: %v", err)
	}

	statsServiceDir, err := findStatisticsServiceDir()
	if err != nil {
		t.Fatalf("failed to locate statistics-service-node directory: %v", err)
	}

	bootstrapRelative := filepath.Join("test-support", "start-test-server.js")
	bootstrapAbsolute := filepath.Join(statsServiceDir, bootstrapRelative)
	if _, err := os.Stat(bootstrapAbsolute); err != nil {
		t.Fatalf("node test bootstrap script not found at %s: %v", bootstrapAbsolute, err)
	}

	cmd := exec.Command(nodePath, bootstrapRelative)
	cmd.Dir = statsServiceDir

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe for Node process: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe for Node process: %v", err)
	}

	stdoutBuf := &synchronizedBuffer{}
	stderrBuf := &synchronizedBuffer{}
	readyCh := make(chan int, 1)

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start Node test server process: %v", err)
	}

	node := &nodeTestServer{
		cmd:     cmd,
		stdout:  stdoutBuf,
		stderr:  stderrBuf,
		waitDone: nil,
	}
	t.Cleanup(func() {
		node.stop(t)
	})

	waitDone := make(chan error, 1)
	node.waitDone = waitDone
	go func() {
		waitDone <- cmd.Wait()
		close(waitDone)
	}()

	go scanOutput(stdoutPipe, stdoutBuf, readyCh)
	go scanOutput(stderrPipe, stderrBuf, nil)

	startupDeadline := time.NewTimer(12 * time.Second)
	defer startupDeadline.Stop()

	var port int
	select {
	case port = <-readyCh:
	case err := <-waitDone:
		if err == nil {
			err = errors.New("process exited before readiness was reported")
		}
		t.Fatalf("node bootstrap exited before readiness: %v\nstdout:\n%s\nstderr:\n%s", err, stdoutBuf.String(), stderrBuf.String())
	case <-startupDeadline.C:
		t.Fatalf("timed out waiting for node bootstrap readiness line\nstdout:\n%s\nstderr:\n%s", stdoutBuf.String(), stderrBuf.String())
	}

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	node.baseURL = baseURL
	if err := waitForHealth(baseURL, 8*time.Second); err != nil {
		t.Fatalf("node test server readiness check failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdoutBuf.String(), stderrBuf.String())
	}

	return node
}

func scanOutput(pipe ioReader, sink *synchronizedBuffer, readyCh chan<- int) {
	scanner := bufio.NewScanner(pipe)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	sentReady := false
	for scanner.Scan() {
		line := scanner.Text()
		sink.appendLine(line)

		if !sentReady && readyCh != nil {
			const prefix = "TEST_SERVER_PORT="
			if strings.HasPrefix(line, prefix) {
				portStr := strings.TrimSpace(strings.TrimPrefix(line, prefix))
				port, err := strconv.Atoi(portStr)
				if err == nil && port > 0 {
					readyCh <- port
					sentReady = true
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		sink.appendRaw("scanner error: " + err.Error())
	}
}

func (n *nodeTestServer) stop(t *testing.T) {
	t.Helper()

	if n == nil || n.cmd == nil {
		return
	}

	select {
	case <-n.waitDone:
		return
	default:
	}

	_ = n.cmd.Process.Signal(os.Interrupt)

	graceTimer := time.NewTimer(1200 * time.Millisecond)
	defer graceTimer.Stop()

	select {
	case <-n.waitDone:
		return
	case <-graceTimer.C:
	}

	_ = n.cmd.Process.Kill()

	killTimer := time.NewTimer(3 * time.Second)
	defer killTimer.Stop()

	select {
	case <-n.waitDone:
		return
	case <-killTimer.C:
		t.Errorf("node test server did not exit after kill\nstdout:\n%s\nstderr:\n%s", n.stdout.String(), n.stderr.String())
	}
}

func waitForHealth(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 300 * time.Millisecond}
	url := baseURL + "/health"
	var lastErr error

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			lastErr = fmt.Errorf("unexpected health status: %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		time.Sleep(50 * time.Millisecond)
	}

	if lastErr == nil {
		lastErr = errors.New("health check timeout")
	}
	return lastErr
}

func buildRealGoApp(statisticsBaseURL string, statisticsTimeout time.Duration) *fiber.App {
	validator := services.NewMatrixValidator()
	statisticsClient := clients.NewHTTPStatisticsClient(statisticsBaseURL, statisticsTimeout)
	qrService := services.NewQRService(validator, statisticsClient)
	qrHandler := handlers.NewQRHandler(qrService)

	app := fiber.New(fiber.Config{AppName: "qr-service-go"})
	app.Use(recover.New())
	handlers.RegisterRoutes(app, qrHandler)

	return app
}

func findStatisticsServiceDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := wd
	for {
		candidates := []string{
			filepath.Join(current, "statistics-service-node"),
			filepath.Join(current, "..", "statistics-service-node"),
		}

		for _, dir := range candidates {
			bootstrap := filepath.Join(dir, "test-support", "start-test-server.js")
			if info, err := os.Stat(bootstrap); err == nil && !info.IsDir() {
				absDir, absErr := filepath.Abs(dir)
				if absErr != nil {
					return "", absErr
				}
				return absDir, nil
			}
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("could not find statistics-service-node/test-support/start-test-server.js from working directory %s", wd)
}

func reserveAndReleaseLocalURL(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to reserve local ephemeral port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("failed to release reserved local port: %v", err)
	}

	return fmt.Sprintf("http://127.0.0.1:%d", port)
}

func readBody(resp *http.Response) ([]byte, error) {
	return ioReadAll(resp.Body)
}

func assertShape(t *testing.T, name string, matrix [][]float64, rows, cols int) {
	t.Helper()
	if len(matrix) != rows {
		t.Fatalf("%s row count mismatch: got %d want %d", name, len(matrix), rows)
	}
	for i := range matrix {
		if len(matrix[i]) != cols {
			t.Fatalf("%s column count mismatch at row %d: got %d want %d", name, i, len(matrix[i]), cols)
		}
	}
}

func assertAllFinite(t *testing.T, name string, matrix [][]float64) {
	t.Helper()
	for i := range matrix {
		for j := range matrix[i] {
			v := matrix[i][j]
			if math.IsNaN(v) || math.IsInf(v, 0) {
				t.Fatalf("%s contains non-finite value at [%d,%d]=%v", name, i, j, v)
			}
		}
	}
}

func assertReconstruction(t *testing.T, a, q, r [][]float64, tol float64) {
	t.Helper()
	reconstructed := multiply(q, r)
	errNorm := frobeniusNormDiff(reconstructed, a)
	threshold := tol * math.Max(1, frobeniusNorm(a))
	if errNorm > threshold {
		t.Fatalf("reconstruction invariant failed: ||Q*R - A||_F=%g threshold=%g", errNorm, threshold)
	}
}

func assertOrthogonality(t *testing.T, q [][]float64, tol float64) {
	t.Helper()
	qtq := multiply(transpose(q), q)
	identity := makeIdentity(len(q))
	errNorm := frobeniusNormDiff(qtq, identity)
	threshold := tol * math.Max(1, float64(len(q)))
	if errNorm > threshold {
		t.Fatalf("orthogonality invariant failed: ||Q^TQ - I||_F=%g threshold=%g", errNorm, threshold)
	}
}

func recomputeStatistics(q, r [][]float64) models.Statistics {
	values := make([]float64, 0, len(q)*len(q[0])+len(r)*len(r[0]))
	for _, row := range q {
		values = append(values, row...)
	}
	for _, row := range r {
		values = append(values, row...)
	}

	maxValue := values[0]
	minValue := values[0]
	sum := 0.0
	for _, v := range values {
		if v > maxValue {
			maxValue = v
		}
		if v < minValue {
			minValue = v
		}
		sum += v
	}

	return models.Statistics{
		Max:               maxValue,
		Min:               minValue,
		Average:           sum / float64(len(values)),
		Sum:               sum,
		HasDiagonalMatrix: isDiagonal(q) || isDiagonal(r),
	}
}

func isDiagonal(matrix [][]float64) bool {
	rows := len(matrix)
	if rows == 0 {
		return false
	}
	cols := len(matrix[0])
	if rows != cols {
		return false
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if i != j && math.Abs(matrix[i][j]) > diagonalEpsilon {
				return false
			}
		}
	}
	return true
}

func assertApproxEqual(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	scale := math.Max(1, math.Max(math.Abs(got), math.Abs(want)))
	if math.Abs(got-want) > tol*scale {
		t.Fatalf("%s mismatch: got %.15g want %.15g tolerance=%g", name, got, want, tol)
	}
}

func multiply(a, b [][]float64) [][]float64 {
	rows := len(a)
	inner := len(a[0])
	cols := len(b[0])

	result := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			sum := 0.0
			for k := 0; k < inner; k++ {
				sum += a[i][k] * b[k][j]
			}
			result[i][j] = sum
		}
	}

	return result
}

func transpose(matrix [][]float64) [][]float64 {
	rows := len(matrix)
	cols := len(matrix[0])
	result := make([][]float64, cols)
	for i := 0; i < cols; i++ {
		result[i] = make([]float64, rows)
		for j := 0; j < rows; j++ {
			result[i][j] = matrix[j][i]
		}
	}
	return result
}

func makeIdentity(n int) [][]float64 {
	identity := make([][]float64, n)
	for i := 0; i < n; i++ {
		identity[i] = make([]float64, n)
		identity[i][i] = 1
	}
	return identity
}

func frobeniusNorm(matrix [][]float64) float64 {
	sum := 0.0
	for i := range matrix {
		for j := range matrix[i] {
			v := matrix[i][j]
			sum += v * v
		}
	}
	return math.Sqrt(sum)
}

func frobeniusNormDiff(a, b [][]float64) float64 {
	sum := 0.0
	for i := range a {
		for j := range a[i] {
			d := a[i][j] - b[i][j]
			sum += d * d
		}
	}
	return math.Sqrt(sum)
}

type ioReader interface {
	Read(p []byte) (n int, err error)
}

func ioReadAll(r ioReader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
