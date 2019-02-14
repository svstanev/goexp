# goexp

Recursive descent expression parser in Go

## Installation

```
go get -u github.com/svstanev/goexp
```

## Usage

```
```

## Expression language

### Syntax Grammar

```
expression      -> logical_or
logical_or      -> logical_and (("||") logical_and)*;
logical_and     -> equality (("&&") equality)*;
equality        -> comparison (("==" | "!=") comparison)*;
comparison      -> addition (("<" | "<=" | ">" | ">=") addition)*;
addition        -> multiplication (("+" | "-") multiplication)*;
multiplication  -> unary (("*" | "/") unary)*;
unary           -> ("!" | "-")? call;
call            -> primary (("(" arguments? ")") | ("." IDENTIFIER))*;
primary         -> "false" | "true" | "nil" | IDENTIFIER | NUMBER | STRING | "(" expression ")";

arguments       -> expression ("," expression)*;
```

### Lexical Grammar

```
IDENTIFIER      -> ALPHA (ALPHA | DIGIT)*;
NUMBER          -> DIGIT* ("." DIGIT*)?;
STRING          -> "'" <any char except "'">* "'"
                  | '"' <any char except '"'>* '"';

DIGIT           -> '0'...'9'
ALPHA           -> 'a'...'z'|'A'...'Z'|'_'
```

