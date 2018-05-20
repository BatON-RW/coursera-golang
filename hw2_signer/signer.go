package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// ExecutePipeline обеспечивает ковейерную обработку задач.
func ExecutePipeline(inJob ...job) {
	in := make(chan interface{}, 1)
	wg := &sync.WaitGroup{}

	for _, j := range inJob {
		wg.Add(1)
		in = execute(j, in, wg)
	}

	wg.Wait()
	return
}

func execute(j job, in chan interface{}, wg *sync.WaitGroup) chan interface{} {
	out := make(chan interface{}, 1)
	go func(res chan interface{}) {
		defer wg.Done()
		defer close(res)
		j(in, res)
	}(out)
	return out
}

// SingleHash считает значение crc32(data)+"~"+crc32(md5(data)) ( конкатенация двух строк через ~),
// где data - то что пришло на вход
func SingleHash(in, out chan interface{}) {
	for rawData := range in {
		data := rawData.(int)
		out <- DataSignerCrc32(strconv.Itoa(data)) + "~" + DataSignerCrc32(DataSignerMd5(strconv.Itoa(data)))
	}

}

// MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки),
// где th=0..5 ( т.е. 6 хешей на каждое входящее значение ), потом берёт конкатенацию результатов в порядке расчета (0..5),
// где data - то что пришло на вход
func MultiHash(in, out chan interface{}) {
	for rawData := range in {
		data := rawData.(string)
		result := ""
		for th := 0; th <= 5; th++ {
			result = result + DataSignerCrc32(strconv.Itoa(th)+data)
		}
		out <- result
	}
}

// CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/),
// объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
func CombineResults(in, out chan interface{}) {
	resSl := []string{}
	for rawData := range in {
		data := rawData.(string)
		resSl = append(resSl, data)
	}
	sort.Slice(resSl, func(i, j int) bool {
		return resSl[i] < resSl[j]
	})
	out <- strings.Join(resSl, "_")
}
