package service

import (
	"fmt"
	"github.com/google/uuid"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type OrderData struct {
	Title string
	Text  []byte
}

type Order struct {
	Id   string
	Path string
}

// orderPath returns the local filepath of an order with the identifier `id`
func OrderPath(id string) string {
	c := GetConfig()
	return c.Datapath + string(os.PathSeparator) + id + ".txt"
}

func (p *OrderData) save(writeFile func(name string, data []byte, perm fs.FileMode) error) error {
	return writeFile(OrderPath(p.Title), p.Text, 0600)
}

func ReadOrderDetails(title string, readFile func(name string) ([]byte, error)) (*OrderData, error) {
	filename := OrderPath(title)

	text, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	return &OrderData{Title: title, Text: text}, nil
}

func CollectOrderDetails(orders chan Order, dir []os.DirEntry, basepath string) {
	go func() {
		for _, d := range dir {
			ext := filepath.Ext(d.Name())
			name := strings.TrimSuffix(d.Name(), ext)
			orders <- Order{Id: name, Path: strings.Join([]string{basepath, "/orders/", name}, "")}
		}
		close(orders)
	}()
}

func ReadOrders(
	readDir func(name string) ([]fs.DirEntry, error),
	config *Config,
) []Order {

	dir, err := readDir(config.Datapath)
	if err != nil {
		log.Println("[OrderListHandler] Error reading orders ", err.Error())
		return []Order{}
	}

	files := make([]Order, 0, len(dir))
	log.Printf("[OrderListHandler] Found '%v' orders", len(dir))

	orders := make(chan Order)

	CollectOrderDetails(orders, dir, config.Basepath)

	for order := range orders {
		files = append(files, order)
	}

	return files
}

func SaveOrder(title string, text string, writeFile func(name string, data []byte, perm fs.FileMode) error) (
	OrderData,
	error,
) {
	order := OrderData{Title: title, Text: []byte(text)}

	return order, order.save(writeFile)
}

// fanOut Creates a number of channels equal to `numWorkers`, to save orders concurrently.
// Returns an array of channels that return the save result
func fanOut(
	input <-chan string,
	numWorkers int,
	writeFile func(name string, data []byte, perm fs.FileMode) error,
) []<-chan SaveOrderResult {
	outputs := make([]<-chan SaveOrderResult, numWorkers)

	for i := 0; i < numWorkers; i++ {
		outputs[i] = goSaveOrder(input, writeFile)
	}

	return outputs
}

type SaveOrderResult struct {
	order OrderData
	err   error
}

// goSaveOrder Saves the orders from the channel `orderIds` and sends the result to returned channel
func goSaveOrder(
	orderIds <-chan string,
	writeFile func(name string, data []byte, perm fs.FileMode) error,
) <-chan SaveOrderResult {
	output := make(chan SaveOrderResult)

	go func() {
		defer close(output)

		for id := range orderIds {
			text := fmt.Sprintf("Order details for %s", id)
			order, err := SaveOrder(id, text, writeFile)

			if err != nil {
				log.Printf("[orderCreateHandler] Error saving order %s: %s\n", id, err.Error())
			}

			output <- SaveOrderResult{
				order: order,
				err:   err,
			}
		}
	}()

	return output
}

// fanIn takes an array of input channels and merges them into a single channel
func fanIn(inputChannels []<-chan SaveOrderResult) <-chan SaveOrderResult {
	output := make(chan SaveOrderResult)
	var wg sync.WaitGroup
	wg.Add(len(inputChannels))

	for _, input := range inputChannels {
		go func(input <-chan SaveOrderResult) {
			defer wg.Done()
			for data := range input {
				output <- data
			}
		}(input)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// GenerateOrders generates `count` orders and saves them using the function `writeFile`
func GenerateOrders(count int, writeFile func(name string, data []byte, perm fs.FileMode) error) []SaveOrderResult {
	input := make(chan string)

	numWorkers := 4
	outputChannels := fanOut(input, numWorkers, writeFile)

	mergedOutput := fanIn(outputChannels)

	go func() {
		defer close(input)
		for i := 1; i <= count; i++ {
			input <- uuid.New().String()
		}
	}()

	var saveOrderResults []SaveOrderResult
	for order := range mergedOutput {
		saveOrderResults = append(saveOrderResults, order)
	}

	return saveOrderResults
}
