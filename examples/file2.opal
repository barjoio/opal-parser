Paragraphs


Exercitation ea officia cillum est non cillum aute sit commodo sit id elit. Culpa commodo nisi irure ex ut laboris non veniam enim duis anim anim elit velit. In ipsum ea culpa velit consequat anim ullamco ex ea. Eiusmod sint consequat Lorem aliqua dolore eiusmod dolor excepteur ut. Occaecat sit et cillum ea qui ipsum. Irure deserunt reprehenderit dolor duis aliqua voluptate mollit do anim ex cupidatat ullamco incididunt.

Mollit amet in aliquip non ipsum. Ex occaecat cillum id mollit excepteur labore ea. Amet dolor exercitation ullamco quis nostrud officia veniam sunt elit reprehenderit magna. Voluptate qui nulla ullamco reprehenderit enim pariatur. Ullamco reprehenderit laborum ullamco ullamco est dolore et nostrud. Qui quis pariatur incididunt labore elit proident minim pariatur minim. Ex labore exercitation ullamco deserunt laborum eu minim nostrud esse.

Sint pariatur ad voluptate reprehenderit cupidatat deserunt proident tempor veniam. Labore incididunt et ea labore occaecat voluptate fugiat occaecat. Laboris dolor in nisi ullamco velit ullamco non nostrud tempor adipisicing esse. Ad amet culpa consectetur laboris ad esse culpa qui est duis voluptate est mollit magna. Commodo ullamco exercitation eiusmod ut qui commodo id laborum nulla irure laborum.



Document settings

.doc
Author: John Smith
Font: Inter 10
Indent size: 20
Margin: 15
Title: How to write documents using Opal



Text formatting


It has `b strong` significance to me.

I `i cannot` stress this enough.

Type `c OK` to accept.

That `bi really` has to go.

Can't pick one? Let's use them `bic all`.

`b C`reate, `b R`ead, `b U`pdate, and `b D`elete

That's fan`i freakin`tastic!

Werewolves are allergic to `h cinnamon`.

Where did all the `u cores` go?

We need `s ten` twenty VMs.

`sp super`script phrase

`sb sub`script phrase

This is a `l hyperlink, https://example.com`

Here is a footnote.`f`

Here is a reference.`r`



Bibliography


.footnotes
- Clarification about this statement.
- Opinions are my own. 

.references
- Andy Hunt & Dave Thomas. The Pragmatic Programmer: From Journeyman to Master. Addison-Wesley. 1999.
- Erich Gamma, Richard Helm, Ralph Johnson & John Vlissides. Design Patterns: Elements of Reusable Object-Oriented Software. Addison-Wesley. 1994.



Headings


.1: Heading level one
.2: Heading level two
.3: Heading level three
.4: Heading level four
.5: Heading level five
.6: Heading level six



Lists


.list/b
- Level 1 list item
- Level 1 list item
- Level 1 list item
  - Level 2 list item
    - Level 3 list item
      - Level 4 list item
        - Level 5 list item
          - Level 6 list item
- Level 1 list item
  - Level 2 list item
  - Level 2 list item

.list/n
- Level 1 list item
- Level 1 list item
- Level 1 list item
  - Level 2 list item
    - Level 3 list item
      - Level 4 list item
        - Level 5 list item
          - Level 6 list item
- Level 1 list item
  - Level 2 list item
  - Level 2 list item

.checklist
- Level 1 checklist item
x Level 1 checklist item
- Level 1 checklist item
  x Level 2 checklist item
    x Level 3 checklist item
      - Level 4 checklist item
        x Level 5 checklist item
          x Level 6 checklist item
x Level 1 checklist item
  - Level 2 checklist item
  x Level 2 checklist item



Images


.image: ~/Pictures/opal.png

.image/800x600
~/Pictures/opal.png

.image
~/Pictures/opal.png
Opal logo

.image/800x600
~/Pictures/opal.png



Code


.code/go -
package main

import (
	"fmt"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	http.HandleFunc("/", greet)
	http.ListenAndServe(":8080", nil)
}
----



Tables


.table: my_data.csv

.table/h
First name | Last name | Age
Adam       | Jones     | 21
Joe        | Bloggs    | 30

.table/h/f
First name | Last name | Age
Adam       | Jones     | 21
Joe        | Bloggs    | 30
Sum        | <         | !sum

.table/h
Position | Name   | Age
!iota    | Jones  | 21
         | Bloggs | 30

.table/h/f/442
First name | Last name | Age
Adam       | Jones     | 21
Joe        | Bloggs    | 30

.table/h -
| First name
| Last name
| Age

| Adam
| Jones
| 21

| Joe
| Bloggs
| 30

| Sum
|
| !sum
----

.table -
| Firefox
| Browser
| Mozilla Firefox is an open source web browser.\n
  It's designed for:\n
  .list
  - standards compliance
  - performance
  - portability
  `l Get Firefox, https://getfirefox.com`!
| Version 83.0

| Firefox
| Browser
| Mozilla Firefox is an open source web browser.\n  
  It's designed for:\n  
  .list
  - standards compliance
  - performance
  - portability
  `l Get Firefox, https://getfirefox.com`!
| Version 83.0
----



Includes


.include: example.opal

.include: https://opalml.org/cheatsheet.opal




// comment








Opal `b Markup`

.table/h/f
First name | Last name | Age
Adam       | Jones     | 21
Joe        | Bloggs    | 30
Sum        | <         | !sum

.1: Functional `i and` elegant
