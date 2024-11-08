package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

// Дополнительная структура для сортировки.
type strFreq struct {
	str  string
	freq int
}

// Сортировка по значениям если они равны то по ключам.
func sortByValue(frequency map[string]int) []strFreq {
	data := make([]strFreq, 0, len(frequency))
	for key, value := range frequency {
		data = append(data, strFreq{key, value})
	}
	sort.SliceStable(data, func(i, j int) bool {
		if data[i].freq < data[j].freq {
			return false
		}
		if data[i].str > data[j].str && data[i].freq == data[j].freq {
			return false
		}
		return true
	})
	if len(data) > 10 {
		data = data[:10]
	}
	return data
}

// Расчет частоты встречыаемых слов на входе пустой map на
// выходе map где ключь слово значение частота.
func calcFreq(freq map[string]int, token []string) map[string]int {
	for _, tok := range token {
		if _, exist := freq[tok]; !exist {
			freq[tok] = 1
			continue
		}
		freq[tok]++
	}
	return freq
}

func Top10(str string) []string {
	token := strings.Fields(str)
	// проверка на пустую строку
	if len(str) == 0 {
		return []string{}
	}
	frequency := make(map[string]int)
	// подсчет вхождений
	frequency = calcFreq(frequency, token)
	// сортировка
	tenFreq := sortByValue(frequency)
	// возвращаемое значение
	out := make([]string, 0, 10)
	for _, val := range tenFreq {
		out = append(out, val.str)
	}
	return out
}
