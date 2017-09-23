package main

import (
	"bufio"
	"log"
	"os"
	"net/http"
	"io/ioutil"
	"strings"
	"sync"
)


func init() {
	log.SetFlags(0)
}


/* Обработать ошибку. Если ошибка, то выход из приложения с кодом */
func errorHandling(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}


/* Отправить GET-запрос, вернуть количество "Go" */
func sendGetRequest(url string) int {

	response, err := http.Get(url)
	defer response.Body.Close()
	errorHandling(err)

	body, err := ioutil.ReadAll(response.Body)
	errorHandling(err)

	goStringCount := counterGoString(string(body))

	//log.Printf("Горутина для %s отработала", url)		// Раскомментировать для наглядности!

	log.Printf("Count for %s: %d", url, goStringCount)

	return goStringCount
}


/* Подсчитать количество строк "Go" */
func counterGoString(request string) int {

	goStringCount := strings.Count(request, "Go")
	return goStringCount
}


/* Подсчитать общее количество строк "Go" на всех страницах */
func totalCounterGoString(bufferedChannel chan int, totalCount *int) {
	for {
		// Читать из канала
		goStringCount, readOk := <- bufferedChannel

		//log.Printf("readOk: %v", readOk)		// Раскомментировать для наглядности!

		// Не пустой ли уже канал
		if !readOk {
			break
		}

		// Добавить считанное количество в общую сумму
		*totalCount += goStringCount

		//log.Printf("Промежуточное значение: %v", *totalCount)	// Раскомментировать для наглядности!
	}
}


/* Главная */
func main() {

	parallelGoroutinesCount := 5		// Максимальное число одновременно работающих горутин
	bufferedChannel	:= make(chan int, parallelGoroutinesCount)		// Буферизированный канал
	var wg sync.WaitGroup 		// Текущее количество работающих горутин
	total := 0			// Суммарное количество строк "Go" на всех страницах

	scanner := bufio.NewScanner(os.Stdin)	// Из пайпа в Scanner
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	// Бежать по URL-ам
	for scanner.Scan() {

		// Ссылка из Сканера
		link := scanner.Text()

		// Увеличить счётчик текущего количества работающих горутин
		wg.Add(1)

		// Количество строк "Go" - в канал
		go func() {
			defer wg.Done()
			bufferedChannel <- sendGetRequest(link)
		}()
		//log.Printf("Запустилась горутина для %s", link)		// Раскомментировать для наглядности!
	}

	// Закрыть канал когда все горутины отработают
	go func () {
		wg.Wait()
		close(bufferedChannel)
	}()

	// Вычитывать из канала количество строк и суммировать
	totalCounterGoString(bufferedChannel, &total)

	log.Printf("Total: %v", total)
}
