package compute

import (
	"strings"
)

const maxLenToken = 3

type state int

const (
	waitState state = iota
	computeState
	// 	stateNumber количество состояний
	stateNumber
)

type symbol int

const (
	otherSymbol symbol = iota
	spaceSymbol
	quoteSymbol
	// symbolNumber количество символов
	symbolNumber
)

type transaction struct {
	event     func(*dataStateMachine, rune)
	nextState state
}

var transitions = [stateNumber][symbolNumber]transaction{
	waitState: {
		otherSymbol: transaction{
			event:     addingEvent,
			nextState: computeState,
		},

		spaceSymbol: transaction{
			event:     skipEvent,
			nextState: waitState,
		},

		quoteSymbol: transaction{
			event:     quoteEvent,
			nextState: computeState,
		},
	},
	computeState: {
		otherSymbol: transaction{
			event:     addingEvent,
			nextState: computeState,
		},

		spaceSymbol: transaction{
			event:     nextTokenEvent,
			nextState: waitState,
		},

		quoteSymbol: transaction{
			event:     nextTokenEvent,
			nextState: waitState,
		},
	},
}

type dataStateMachine struct {
	sb           strings.Builder
	tokens       []string
	quote        rune
	currentState state
}

func parsingQuery(req string) ([]string, error) {
	data := &dataStateMachine{
		tokens: make([]string, 0, maxLenToken),
	}

	runes := []rune(req)
	number := len(runes)
	var previous rune
	var current rune
	for i := 0; i < number; i++ {
		previous = current
		current = runes[i]

		// Проверка на конец строки запроса
		// 1. Текущий символ переноса строки
		// 2. Текущий символ не часть того что между скобок
		if isNewlineSymbol(current) && isRuneEmpty(data.quote) {
			nextState(data, spaceSymbol, current)
			continue
		}

		// Проверка на открывающую кавычку
		// 1. Символ является кавычкой
		// 2. Ещё не было открывающей кавычки
		// 3. Предыдущий символ пустой или пробел
		if isQuoteSymbol(current) && isRuneEmpty(data.quote) && (isRuneEmpty(previous) || isSpaceSymbol(previous)) {
			nextState(data, quoteSymbol, current)
			continue
		}

		// Проверка на закрывающую кавычку
		// 1. Текущий символ равен открывающей кавычки
		// 2. Предыдущей символ не должен быть экранирующим
		if current == data.quote && previous != '\\' {
			nextState(data, quoteSymbol, current)
			continue
		}

		// Проверка на разделяющий пробел
		// 1. Символ является пробелом
		// 2. Нет открытой кавычки
		if isSpaceSymbol(current) && isRuneEmpty(data.quote) {
			nextState(data, spaceSymbol, current)
			continue
		}

		// Общая ситуация для всех символов, которые не попали под условие выше
		nextState(data, otherSymbol, current)
	}

	// Оставшиеся символы в strings.Builder выгружаем в массив токенов
	// rune можно подставить любую на консистентность токена не повлияет
	nextTokenEvent(data, '\x00')

	return data.tokens, nil
}

// nextState переход между состояниями
func nextState(d *dataStateMachine, st symbol, symbol rune) {
	t := transitions[d.currentState][st]
	t.event(d, symbol)
	d.currentState = t.nextState
}

// addingEvent - добавление нового символа в формирующийся токен
func addingEvent(d *dataStateMachine, symbol rune) {
	d.sb.WriteRune(symbol)
}

// quoteEvent - фиксируем тип открывающей скобки, так как их два типа " и '
func quoteEvent(d *dataStateMachine, symbol rune) {
	d.quote = symbol
}

// skipEvent - метод используется, чтобы пропустить символ
func skipEvent(_ *dataStateMachine, _ rune) {
	return
}

// nextTokenEvent - метод для сохранения токена
func nextTokenEvent(d *dataStateMachine, _ rune) {
	if d.sb.Len() == 0 {
		return
	}

	d.tokens = append(d.tokens, d.sb.String())
	d.sb.Reset()
	d.quote = '\x00'
}

// isSpaceSymbol - метод, которая проверяет, является ли переданный символ символом пробела, новой строки
// или табуляции.
// Проверка на перенос строки '\t' и табуляцию '\n'
func isSpaceSymbol(symbol rune) bool {
	return symbol == ' '
}

// isSpaceSymbol - метод, которая проверяет, является ли переданный символ символом пробела, новой строки
// или табуляции.
// Проверка на перенос строки '\t' и табуляцию '\n'
func isQuoteSymbol(symbol rune) bool {
	return symbol == '\'' || symbol == '"'
}

// isWhitespaceSymbol - метод, которая проверяет, является ли переданный символ символом пробела, новой строки
// или табуляции.
// Проверка на перенос строки '\t' и табуляцию '\n'
func isNewlineSymbol(symbol rune) bool {
	return symbol == '\n'
}

// isRuneEmpty Проверка руны на пустоту
func isRuneEmpty(r rune) bool {
	return r == '\x00'
}
