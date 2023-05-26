// Copyright (c) 2023 Luka Ivanovic
// This code is licensed under MIT licence (see LICENCE for details)

package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type TokenType int8
type AtomType int8
type ExprType int8

const (
	TOKEN_LPAREN TokenType = iota
	TOKEN_RPAREN
	TOKEN_SYMBOL
	TOKEN_NUMBER
	ATOM_INT AtomType = iota
	ATOM_FLOAT
	ATOM_SYMBOL
	EXPR_ATOM ExprType = iota
	EXPR_LIST
)

type Token struct {
	tokenType  TokenType
	tokenValue string
}

type Atom struct {
	atomKind  AtomType
	atomValue interface{}
}

type Expr struct {
	exprKind  ExprType
	exprValue interface{}
}

func (t TokenType) String() string {
	switch t {
	case TOKEN_LPAREN:
		return "LPAREN"
	case TOKEN_RPAREN:
		return "RPAREN"
	case TOKEN_NUMBER:
		return "NUMBER"
	default:
		return "SYMBOL"
	}
}
func (t ExprType) String() string {
	switch t {
	case EXPR_ATOM:
		return "ATOM"
	default:
		return "LIST"
	}
}

func tokenize(source string) []Token {
	source = strings.ReplaceAll(source, "(", " ( ")
	source = strings.ReplaceAll(source, ")", " ) ")
	fields := strings.Fields(source)
	tokens := make([]Token, 0)
	for _, field := range fields {
		switch field {
		case "(":
			tokens = append(tokens, Token{TOKEN_LPAREN, field})
		case ")":
			tokens = append(tokens, Token{TOKEN_RPAREN, field})
		default:
			if _, err := strconv.ParseInt(field, 10, 64); err == nil {
				tokens = append(tokens, Token{TOKEN_NUMBER, field})
			} else if _, err := strconv.ParseFloat(field, 64); err == nil {
				tokens = append(tokens, Token{TOKEN_NUMBER, field})
			} else {
				tokens = append(tokens, Token{TOKEN_SYMBOL, field})
			}
		}
	}
	return tokens
}

func printTokens(tokens []Token) {
	fmt.Println("[")
	for _, token := range tokens {
		fmt.Printf("\t%s = %s\n", token.tokenType, token.tokenValue)
	}
	fmt.Println("]")
}

func nextToken(tokens *[]Token) Token {
	t := (*tokens)[0]
	*tokens = (*tokens)[1:]
	return t
}

func parse(tokens *[]Token) (Expr, error) {
	if len(*tokens) == 0 {
		return Expr{}, errors.New("unexpected EOF")
	}
	token := nextToken(tokens)
	switch token.tokenType {
	case TOKEN_RPAREN:
		return Expr{}, errors.New("unexpected )")
	case TOKEN_NUMBER, TOKEN_SYMBOL:
		return Expr{EXPR_ATOM, parseAtom(token)}, nil
	default: // TOKEN_LPAREN
		res := make([]Expr, 0)
		for (*tokens)[0].tokenType != TOKEN_RPAREN {
			tmp, err := parse(tokens)
			if err != nil {
				return Expr{}, err
			}
			res = append(res, tmp)
		}
		*tokens = (*tokens)[1:]
		return Expr{EXPR_LIST, res}, nil
	}
}

func parseAtom(token Token) Atom {
	if number, err := strconv.ParseInt(token.tokenValue, 10, 64); err == nil {
		return Atom{ATOM_INT, number}
	} else if number, err := strconv.ParseFloat(token.tokenValue, 64); err == nil {
		return Atom{ATOM_FLOAT, number}
	} else {
		return Atom{ATOM_SYMBOL, token.tokenValue}
	}
}

func main() {
	tokens := tokenize("(begin (define r 10) (* pi (* r r)))")
	printTokens(tokens)
	ast, err := parse(&tokens)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ast)
}
