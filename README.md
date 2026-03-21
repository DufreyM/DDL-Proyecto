# Generador de Analizadores Léxicos (YALex → Go)

## Descripción

Este proyecto implementa un **generador de analizadores léxicos** basado en especificaciones escritas en **YALex**.

A partir de un archivo `.yal`, el sistema:

1. Parsea las definiciones léxicas
2. Convierte expresiones regulares a **postfix**
3. Construye un **Árbol de Expresión (AST)**
4. Genera un **AFN (Thompson)**
5. Convierte el AFN a **AFD (Subset Construction)**
6. Genera un **lexer funcional en Go**
7. Ejecuta el lexer sobre un archivo de entrada

---

## Objetivos

* Implementar un generador de analizadores léxicos
* Aplicar teoría de autómatas finitos
* Transformar expresiones regulares en autómatas
* Generar código ejecutable automáticamente

---

## Arquitectura

El sistema sigue el pipeline clásico de compiladores:

```
YAL → Regex → Postfix → AST → AFN → AFD → Lexer
```

### Módulos principales

* `yal/` → Parser de archivos `.yal`
* `regex/` → Manejo de expresiones regulares (postfix + AST)
* `automata/` → Construcción de AFN y AFD
* `lexer/` → Ejecución del lexer
* `generator/` → Generación de código (`generated_lexer.go`)
* `graph/` → Generación de árboles (DOT)

---

## Funcionamiento Interno

### 1. Parsing de YAL

Se leen:

* `let` → definiciones auxiliares
* `rule gettoken` → reglas léxicas

---

### 2. Expresiones Regulares

Las expresiones:

* se expanden (`let`)
* se convierten a postfix
* se transforman en AST

---

### 3. Construcción de AFN (Thompson)

Cada expresión regular se convierte en un AFN usando:

* concatenación
* unión (`|`)
* Kleene (`*`)
* plus (`+`)
* opcional (`?`)

---

### 4. Conversión AFN → AFD

Se aplica **subset construction**:

* ε-closure
* move
* construcción de estados deterministas

---

### 5. Generación de Lexer

Se genera automáticamente:

```
generated_lexer.go
```

Incluye:

* tabla de transiciones
* estados finales
* función `NextToken`

---

## Estructura del Proyecto

```
yalex-full/
│
├── main.go
├── generated_lexer.go
│
├── yal/
├── regex/
├── automata/
├── lexer/
├── generator/
├── graph/
```

---

## Cómo Ejecutar

### 1. Ejecutar generador + lexer

```
go run main.go archivo.yal archivo.txt
```

Ejemplo:

```
go run main.go arnoldc.yal test.arnoldc
```

---

### 2. Ejecutar lexer generado

```
go run generated_lexer.go
```

---

## Salida

El sistema imprime:

```
KW_PRINT -> TALK TO THE HAND
IDENT -> x
INT_LIT -> 10
```

o errores:

```
LEXICAL ERROR line 3: @
```

---

## Árbol de Expresión

Se genera un archivo:

```
tree.dot
```

Convertir a imagen:

```
dot -Tpng tree.dot -o tree.png
```

---

## Decisiones de Diseño

* Los **strings** y **números** se manejan directamente en el lexer para evitar ambigüedades en el DFA
* Se respeta:

  * **longest match**
  * **prioridad de reglas**
* El DFA se usa como estructura principal de reconocimiento

---

## Características

✔ Soporte de expresiones regulares
✔ Construcción de AFN (Thompson)
✔ Conversión a AFD
✔ Generación automática de código
✔ Manejo de errores léxicos
✔ Soporte para keywords complejos (ArnoldC)

---

## Limitaciones

* Soporte parcial de clases como `[^\n]`
* No incluye minimización de DFA
* El árbol puede crecer exponencialmente en regex grandes

---

## Pruebas

Se probó con:

* ArnoldC (lenguaje basado en frases)
* Casos válidos e inválidos
* Identificadores, enteros y strings

---

## Demo

Se incluye un video demostrando:

* generación del lexer
* ejecución
* manejo de errores

---

## Autor

Proyecto realizado como parte del curso de **Diseño de Lenguajes de Programación**

---

## Conclusión

Este proyecto demuestra la implementación completa de un generador léxico, aplicando teoría formal de autómatas para producir una herramienta funcional capaz de analizar lenguajes definidos por expresiones regulares.

---
