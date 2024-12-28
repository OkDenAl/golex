# golex
It is a simple lexical analyzers generator written on pure Go
with custom regular expression parser.

## Quick start

1) The token is specified as a Latin word written on the same line with a regular expression. 
If the token is found, a token with that name will be returned.
```
%%
/(0|(1)+)/ Num
%%
```

2) If a token does not imply any reaction to its detection, then it is allowed to skip this token using the `continue` 
construction following the name of the token. However, you still need to specify the name of the token itself.

```
%%
/[\n\t ]/ Skip continue
%%
```

3) In the rules section, it is allowed to use regular expressions that depend on the starting conditions 
of the following type: `<n>/r/`, where n is the starting condition, r is the regular expression. 
Initially, the `INIT` start condition is enabled, which is not specified in the rules. 
The starting conditions are switched by specifying the `begin` construction.
```
%%
/\"/ RegularStart begin(REGULAR) edit
<REGULAR>/\"/ RegularEnd begin(INIT) edit
<REGULAR>/\n/ RegularNewLine begin(INIT) edit
%%
```
The edit construct signals that it is necessary to make a change to the standard token handler, 
that is, to redefine the standard lex behavior.

4) In case of simplification of the work, it is possible to specify named regular expressions,
using this expression later inside the regular expression record.
```
NUM             /(0|(1)+)/
%%
/{NUM}/ Num
%%
```

## Grammar (RBNF)

```
Program         ::= (NamedRegExpr)* (NewLine)* (State)? (NewLine)* Rules
NamedRegExpr    ::= Name "/" RegExpr "/" NewLine
State           ::= "%x" (Name)+ NewLine
Rules           ::= "%%" (NewLine)+ (Rule)+ "%%"
Rule            ::= (StartCondition)? "/" RegExpr "/"  Name (SwitchCondition)? (NewLine)+
StartCondition  ::= "<" Name ">"
SwitchCondition ::= "BEGIN" "(" Name ")"

RegExpr         ::= Union
Union           ::= Concatenation ("|" Concatenation)*
Concatenation   ::= (BasicExpr)+
BasicExpr       ::= Element ("*"|"+"|"?")?
Element         ::= Group | Set | Escape | ValidIndependentCharacter
Group           ::= "(" RegExpr ")"
Escape          ::= "\" AnyCharacter
Set             ::= "[" ("^")? FirstSetItem SetItems "]"
SetItems        ::= SetItem SetItems
SetItem         ::= Range | Escape | SetCharacter
Range           ::= ("\" (AnyCharacter)? | RangeCharacter) "-"
                    ("\" (AnyCharacter)? | RangeCharacter)


Lexems
Name ::= [A-Z][A-z0-9_]*
NewLine ::= \\n
Continue ::= continue
Edit ::= edit
ValidIndependentCharacter ::= [^()|\/]
AnyCharacter ::= .
SetCharacter ::= .[^]]*
RangeCharacter ::= [^\]
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
\w| word, [a-z-A-Z0-9_]|\w{4}|v411|v4.1

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

> *Supported ASCII range from 0 to 1000