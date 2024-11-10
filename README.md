# golex
It is a simple lexical analyzers generator written on pure Go
with custom regular expression parser.

## Grammar (RBNF)

```
Program         ::= (NamedRegExpr)* (NewLine)* (State)? (NewLine)* Rules
NamedRegExpr    ::= Name "/" RegExpr "/" NewLine
State           ::= "%x" (Name)+ NewLine
Rules           ::= "%%" (NewLine)+ (Rule)+ "%%"
Rule            ::= (StartCondition)? "/" RegExpr "/"  Name (SwitchCondition)? (NewLine)+
StartCondition  ::= "<" Name ">"
SwitchCondition ::= "BEGIN" "(" Name ")"

RegExpr         ::= Union | SimpleExpr
Union           ::= SimpleExpr "|" RegExpr
SimpleExpr      ::= Concatenation | BasicExpr
Concatenation   ::= BasicExpr SimpleExpr
BasicExpr       ::= Element ("*"|"+"|"?"|Repetition)?
Repetition      ::= ("{" Number ("}" | "," ("}" | Number "}")))
Element         ::= Group | Set | Escape | Repetition | ValidIndependentCharacter
Group           ::= "(" RegExpr ")"
Escape          ::= "\" EscapeCharacter
Set             ::= "[" ("^")? FirstSetItem SetItems "]"
SetItems        ::= SetItem SetItems
SetItem         ::= Range | Escape | SetCharacter
Range           ::=  (Escape | RangeStartCharacter) "-" RangeEndCharacter


Lexems
Name ::= [A-Z][A-z0-9_]*
NewLine ::= \\n
ValidIndependentCharacter ::= [^()|/]
EscapeCharacter ::= .
SetCharacter ::= .[^]]*
RangeStartCharacter ::= .
RangeEndCharacter ::= .
```

## Supported regular expressions features

|Char class|Description|Example|Valid match|Invalid|
:---|:---|:---|:---|---
[ ]|class definition|[axf]|a, x, f|b
[ - ]|class definition range|[a-c]|a, b, c|d
[ \ ]|escape inside class|[a-f\.]|a, b, .| g
[^ ]|Not in class|[^abc]|d, e| a
.|match any chars except new line|b.ttle|battle, bottle| bttle
\s|white space, [\n\r\f\t ]|good\smorning|good morning|good.morning
\S|no-white space, [^\n\r\f\t]|good\Smorning|good.morning|good morning
\d| digit|\d{2}|23|1a
\D| non-digit|\D{3}|foo, bar|fo1
\w| word, [a-z-A-Z0-9_]|\w{4}|v411|v4.1
\W|non word, [^a-z-A-Z0-9_]|.$%?|.$%?|.ab?

|Special character|Description
:---|:---
\\ |general escape|
\n|new line|
\r|carriage return|
\t|tab|

|Sequence|Description|Example|Valid match|Invalid|
:---|:---|:---|:---|---
\||alternation|apple\|orange|apple, orange|melon
( )| subpattern |foot(er\|ball)|footer or football|footpath
+| one or more quantifier|ye+ah|yeah, yeeeah|yah  
*| zero or more quantifier|ye*ah|yeeah, yeeeah, yah|yeh 
?| zero or one quantifier|yes?|yes, ye|yess
{n}|n times exactly|fo{2}|foo|fooo
{n,m}|from n to m times|go{2,3}d|good,goood|gooood
{n,}|at least n times|go{2,}|goo, gooo|go

> *Supported ASCII range from 0 to 1000