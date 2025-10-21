package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type RequestResult struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

type Report struct {
	TotalTime          time.Duration
	TotalRequests      int
	SuccessRequests    int
	StatusDistribution map[int]int
	mutex              sync.Mutex
}

func (r *Report) AddResult(result RequestResult) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if result.Error == nil {
		if result.StatusCode == 200 {
			r.SuccessRequests++
		}
		r.StatusDistribution[result.StatusCode]++
	} else {
		r.StatusDistribution[0]++ // 0 para erros de conexão
	}
}

func (r *Report) Print() {
	fmt.Println("\n========== RELATÓRIO DE TESTE DE CARGA ==========")
	fmt.Printf("Tempo total de execução: %v\n", r.TotalTime)
	fmt.Printf("Quantidade total de requests: %d\n", r.TotalRequests)
	fmt.Printf("Requests com status 200: %d\n", r.SuccessRequests)
	fmt.Println("\nDistribuição de status HTTP:")

	for status, count := range r.StatusDistribution {
		if status == 0 {
			fmt.Printf("  Erros de conexão: %d\n", count)
		} else {
			fmt.Printf("  Status %d: %d requests\n", status, count)
		}
	}
	fmt.Println("================================================")
}

func makeRequest(url string, client *http.Client) RequestResult {
	start := time.Now()

	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return RequestResult{
			StatusCode: 0,
			Duration:   duration,
			Error:      err,
		}
	}
	defer resp.Body.Close()

	return RequestResult{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}

func worker(id int, url string, jobs <-chan int, results chan<- RequestResult, wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()

	for range jobs {
		result := makeRequest(url, client)
		results <- result
	}
}

func runLoadTest(url string, totalRequests, concurrency int) *Report {
	report := &Report{
		TotalRequests:      totalRequests,
		StatusDistribution: make(map[int]int),
	}

	// Cliente HTTP reutilizável para melhor performance
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: concurrency,
		},
	}

	jobs := make(chan int, totalRequests)
	results := make(chan RequestResult, totalRequests)

	var wg sync.WaitGroup
	var processed int32

	// Iniciar workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(i, url, jobs, results, &wg, client)
	}

	// Goroutine para processar resultados
	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for result := range results {
			report.AddResult(result)
			current := atomic.AddInt32(&processed, 1)
			if current%100 == 0 || int(current) == totalRequests {
				fmt.Printf("\rProgresso: %d/%d requests completados", current, totalRequests)
			}
		}
	}()

	// Enviar jobs
	startTime := time.Now()
	for i := 0; i < totalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	// Aguardar conclusão dos workers
	wg.Wait()
	close(results)

	// Aguardar processamento dos resultados
	resultsWg.Wait()

	report.TotalTime = time.Since(startTime)
	fmt.Println() // Nova linha após o progresso

	return report
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")

	flag.Parse()

	// Validações
	if *url == "" {
		log.Fatal("Erro: O parâmetro --url é obrigatório")
	}

	if *requests <= 0 {
		log.Fatal("Erro: O parâmetro --requests deve ser maior que 0")
	}

	if *concurrency <= 0 {
		log.Fatal("Erro: O parâmetro --concurrency deve ser maior que 0")
	}

	if *concurrency > *requests {
		*concurrency = *requests
	}

	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Total de requests: %d\n", *requests)
	fmt.Printf("Concorrência: %d\n\n", *concurrency)

	report := runLoadTest(*url, *requests, *concurrency)
	report.Print()
}
