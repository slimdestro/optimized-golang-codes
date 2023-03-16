/**
	@ internet speedtest golang implementation
	@ slimdestro
*/
package speedtest

import (
	"net/http"
	"time"
)

// Zone represents a geographical region
type Zone struct {
	Name string
	URL  string
}

// SpeedTest is used to test user internet speed in different zones with different size of payload
type SpeedTest struct {
	Zones []Zone
	Size  int
}

// NewSpeedTest creates a new SpeedTest instance
func NewSpeedTest(zones []Zone, size int) *SpeedTest {
	return &SpeedTest{
		Zones: zones,
		Size:  size,
	}
}

/** 
	# use this function to generate real payload
*/

func generatePayload(size int) *strings.Reader {
	// Generate a string of the specified size
	str := ""
	for i := 0; i < size; i++ {
		str += "a"
	}

	// Create a reader from the string
	return strings.NewReader(str)
}

// generatePayload is a helper function to generate a payload of the specified size
func generatePayload(size int) *strings.Reader {
	// Generate a string of the specified size
	str := ""
	for i := 0; i < size; i++ {
		str += "a"
	}

	// Create a reader from the string
	return strings.NewReader(str)
}

// Test runs the speed test for the given zone and size
func (s *SpeedTest) Test(zone Zone) (time.Duration, error) {
	start := time.Now()
	resp, err := http.Get(zone.URL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// read the response body
	buf := make([]byte, s.Size)
	_, err = resp.Body.Read(buf)
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

// Example usage
func main() {
	zones := []Zone{
		{Name: "US East", URL: "https://site1"},
		{Name: "US West", URL: "https://site2"},
	}

	st := NewSpeedTest(zones, 1024) // 1KB

	for _, zone := range st.Zones {
		duration, err := st.Test(zone)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: %s\n", zone.Name, duration)
	}
}