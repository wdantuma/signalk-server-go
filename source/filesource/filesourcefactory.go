package filesource

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/wdantuma/signalk-server-go/source"
	"github.com/wdantuma/signalk-server-go/source/filesource/actisensesource"
	"github.com/wdantuma/signalk-server-go/source/filesource/candumpsource"
)

func CreateFileSource(filename string) (source.CanSource, error) {

	// very simple factory impl for now

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}
	parts := strings.Split(scanner.Text(), "#")
	if len(parts) == 2 {
		return candumpsource.NewCanDumpSource(filename)
	}
	parts = strings.Split(scanner.Text(), ",")

	if len(parts) >= 14 {
		return actisensesource.NewActiSenseSource(filename)
	}

	return nil, errors.New("No matching format found")
}
