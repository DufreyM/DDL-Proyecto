# Generador de Analizadores Léxicos y Sintácticos (YALex + YAPar)

## Descripción

Este proyecto implementa un generador de analizadores léxicos y sintácticos inspirado en herramientas clásicas como:

- Lex / Flex
- Yacc / Bison
- ocamllex
- ocamlyacc

El sistema recibe:

- un archivo `.yal` (especificación léxica)
- un archivo `.yalp` (especificación sintáctica)
- un archivo de entrada

y construye automáticamente:

1. Un Analizador Léxico
2. Un Parser SLR(1)
3. Una Tabla de Símbolos
4. Una Tabla SLR
5. Un sistema completo de parsing shift-reduce

---

# Objetivos

- Implementar un generador de analizadores léxicos
- Implementar un generador de analizadores sintácticos
- Aplicar teoría formal de compiladores
- Construir autómatas finitos
- Implementar algoritmos LR(0) y SLR(1)
- Generar código automáticamente
- Ejecutar parsing completo end-to-end

---

# Arquitectura General

El sistema sigue el pipeline clásico de compiladores:

.yal
 ↓
Regex
 ↓
Postfix
 ↓
AST
 ↓
AFN (Thompson)
 ↓
AFD (Subset Construction)
 ↓
Lexer
 ↓
Token Stream
 ↓
.yalp
 ↓
Grammar
 ↓
FIRST/FOLLOW
 ↓
LR(0)
 ↓
SLR Table
 ↓
Shift-Reduce Parser

---

# Módulos Principales

## yal/
Parser de archivos `.yal`

Responsable de:
- leer definiciones `let`
- expandir expresiones regulares
- extraer reglas léxicas

---

## regex/
Manejo de expresiones regulares

Incluye:
- concatenación explícita
- conversión a postfix
- manejo de operadores
- soporte de escapes
- construcción de AST

---

## automata/
Construcción de autómatas

Implementa:
- AFN usando Thompson
- ε-closures
- subset construction
- conversión AFN → AFD

---

## lexer/
Motor léxico

Responsable de:
- tokenización
- manejo de errores léxicos
- longest match
- tabla de símbolos

---

## yapar/
Parser de archivos `.yalp`

Responsable de:
- lectura de tokens
- extracción de producciones
- construcción de gramáticas

---

## syntax/
Parser Generator SLR(1)

Implementa:
- FIRST
- FOLLOW
- Items LR(0)
- closure()
- goto()
- colección canónica LR(0)
- tabla SLR(1)
- parser shift-reduce

---

## generator/
Generación automática de código Go

Genera:
generated_lexer.go

---

## graph/
Visualización de árboles

Genera:
tree.dot

---

# Funcionalidades Implementadas

## Analizador Léxico

✔ Parsing de `.yal`
✔ Expansión de expresiones regulares
✔ Soporte de operadores:
- |
- *
- +
- ?

✔ Conversión a postfix
✔ Construcción de AST
✔ Thompson Construction
✔ Subset Construction
✔ DFA funcional
✔ Longest Match
✔ Prioridad de reglas
✔ Manejo de errores léxicos
✔ Generación automática de lexer

---

# Analizador Sintáctico

✔ Parsing de `.yalp`
✔ Lectura de `%token`
✔ Extracción de producciones
✔ Construcción de gramática

---

# FIRST / FOLLOW

Implementación funcional de:
- FIRST(X)
- FOLLOW(X)

---

# LR(0)

Implementación completa de:
- LR Items
- closure()
- goto()
- colección canónica

---

# Tabla SLR(1)

Construcción de:
- ACTION
- GOTO
- reduce
- accept

---

# Shift-Reduce Parser

Parser funcional con:
- stack
- shift
- reduce
- accept
- manejo de errores sintácticos

---

# Tabla de Símbolos

El sistema implementa una tabla de símbolos básica:

LEXEME=abc TOKEN=ID LINE=1
LEXEME=+ TOKEN=PLUS LINE=1
LEXEME=123 TOKEN=INT_LIT LINE=1

Incluye:
- lexema
- token
- línea

---

# Estructura del Proyecto

yalex-full/

├── main.go
├── generated_lexer.go
├── lexer.yal
├── test.yalp
├── input.txt

├── yal/
├── yapar/
├── regex/
├── automata/
├── lexer/
├── syntax/
├── generator/
├── graph/

---

# Cómo Ejecutar

## Ejecutar sistema completo

go run main.go lexer.yal input.txt

---

# Ejemplo de Entrada

## lexer.yal

let digit = ['0'-'9']
let letter = ['a'-'z']

rule gettoken =
  | digit+ { return INT_LIT }
  | letter+ { return ID }
  | [+] { return PLUS }

---

## test.yalp

%token ID
%token PLUS
%token INT_LIT

%%

expr:
 expr PLUS term
 | term
;

term:
 INT_LIT
 | ID
;

---

## input.txt

abc+123

---

# Ejemplo de Salida

TOKENS:
ID -> abc
PLUS -> +
INT_LIT -> 123

SYMBOL TABLE:
LEXEME=abc TOKEN=ID LINE=1
LEXEME=+ TOKEN=PLUS LINE=1
LEXEME=123 TOKEN=INT_LIT LINE=1

TOKEN STREAM:
[ID PLUS INT_LIT]

FIRST(expr):
[INT_LIT ID]

FOLLOW(expr):
[$ PLUS]

Generated 7 states

INPUT ACCEPTED

PARSE SUCCESS

---

# Estados LR(0)

STATE 0
expr' -> [expr] (dot=0)
expr -> [expr PLUS term] (dot=0)
expr -> [term] (dot=0)
term -> [INT_LIT] (dot=0)
term -> [ID] (dot=0)

---

# Árbol de Expresión

El sistema genera:
tree.dot

Convertir a imagen:
dot -Tpng tree.dot -o tree.png

---

# Decisiones de Diseño

- El lexer utiliza DFA como estructura principal
- Se respeta:
  - longest match
  - prioridad de reglas
- Los parsers LR se implementaron modularmente
- Se separaron lexer y parser generator
- El parser utiliza SLR(1)

---

# Características

✔ Regex parser
✔ AST generation
✔ Thompson Construction
✔ DFA generation
✔ Lexer generation
✔ Symbol Table
✔ YAPar parser
✔ FIRST/FOLLOW
✔ LR(0)
✔ SLR(1)
✔ Shift-Reduce Parser
✔ Parsing completo end-to-end

---

# Limitaciones

- No se implementa minimización de DFA
- Manejo parcial de clases complejas
- No se implementan conflictos avanzados LR
- No incluye recuperación de errores sintácticos

---

# Pruebas Realizadas

Se probaron:
- identificadores
- enteros
- operadores
- expresiones válidas
- errores léxicos
- parsing SLR completo

---

# Demo

El proyecto incluye demostraciones de:

- generación de lexer
- construcción de DFA
- generación de estados LR
- construcción de tabla SLR
- parsing shift-reduce
- aceptación/rechazo de cadenas

---

# Autores

Proyecto realizado para el curso:

Diseño de Lenguajes de Programación

Integrantes:
- Leonardo Mejía
- Mía Fuentes
- María Girón

---

# Conclusión

Este proyecto implementa exitosamente un generador de analizadores léxicos y sintácticos completo, aplicando teoría formal de compiladores para construir un sistema funcional basado en:

- expresiones regulares
- autómatas finitos
- parsing LR(0)
- tablas SLR(1)
- parsing shift-reduce

El sistema logra ejecutar un pipeline completo de compilación desde especificaciones `.yal` y `.yalp` hasta el reconocimiento sintáctico final de cadenas de entrada.