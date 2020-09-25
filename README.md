# Toy Language Interpreter

## Reserved words
* print
* if
* else
* while

## Punctuation and operators
(	+	=	<  
)	-	==	>  
{	*	!=	<=  
}	/	&&	>=  
,	%	||  

## Other lexical rules
* Each number consists of one or more digits, and denotes a non-negative integer.

* Each identifier consists of one or more letters that do not form a reserved word. Reserved words and identifiers are case sensitive. That is, if denotes a reserved word, but If and iF and IF are each distinct identifiers.

* Whitespace characters include blanks, tabs, line feeds, and carriage returns.

## Statements
A program is a sequence of statements. Each statement is one of the following:

#### Statement type
**assignment** identifier = expression  
**print**	  print expression1 , expression2 ... , expressionN  
**selection**	if ( expression ) statement1 else statement2  
**iteration**	while ( expression ) statement  
**compound**	{ statement1 statement2 ... statementN }  

## Expressions
Binary operators have the same meanings, precedence, and associativity as in the C language. Parentheses force an evaluation order. There are no unary operators.

## Types
* Each identifier denotes a variable that has integer type and global scope.
* Arithmetic operators (+, −, \*, /, %) return integer values.
* Relational and logical operators (==, !=, <, >, <=,>=, &&, ||) also return integer values (1 for true, 0 for false).
* The conditional expressions in if and while statements evaluate to 1 or 0.

## Limitations
* Although all print statements are evaluated correctly, only the last one is actually printed.
* The branches of the conditionals are parsed as blocks ie. with {} around them

## Acknowledgements
* The design of the interpreter was inspired by the excellent book "Writing an interpreter with Go" by Thorsten Ball 



