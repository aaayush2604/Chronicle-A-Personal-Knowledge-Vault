# Chronicle-A-Personal-Knowledge-Vault

This project started as an attempt by me to learn GO

Just building a personal project allowing me to note down my thoughts into different categories and being able to use those thoughts/learnings/questions at the right time. This project currently is just a simple append-only log that lets you enter notes and allows you to pull those thoughts based on the words you might think are used in those entries. If trying to see all the features that this currently has, use help in the terminal. The path to logFile is displayed on the terminal when the project starts. If you want to add a particular path for the logFile, apart from the default, just add CHRONICLE_LOG_FILE_PATH to an env file

## CQL (Chronicle Query Language)

I have added a basic query Language to this, allowing you to query your notes based on the words an entry contains, the entry's type and what date or time the entry was made. I am yet to enter a testing phase, so you might see a lot of bugs, kindly let me know by raising a pr, on the bugs.md

# User Manual

## Remembering Entries

Use `rem` or `remember` to add new entries.

```bash
rem [@type] [#tags...] <text>
remember [@type] [#tags...] <text>
```

Supported entry types:

- `note` (`n`)
- `idea` (`i`)
- `question` (`q`)
- `learning` (`l`)
- `important` (`imp`)

If no type is specified, the entry is stored as a note.

### Examples

```bash
rem "Read chapter 3"

rem @learning "Definition of uniform convergence"

rem @idea #project #go "Use an append-only log"

remember @question #math "Why does this proof work?"
```

---

## Recall

Use `recall` to search your entries using query expressions.

```bash
recall <predicate> AND/OR <predicate> ...
```

**Note:** `AND` has higher precedence than `OR`.

---

## Supported Predicates

- `all`
- `contains[...]`
- `type[...]`
- `tags[...]`
- `date`
- `time`
- `len`

---

## Operators

Comparison operators:

```text
=    !=    >    <    >=    <=
```

Logical operators:

```text
AND    OR
```

---

## List Semantics

The following predicates accept lists:

- `contains[...]`
- `type[...]`
- `tags[...]`

Lists use **OR semantics**.

For example:

```bash
type[note,idea]
tags[go,database]
contains["ml","project"]
```

match entries satisfying **any one** of the values.

To require multiple conditions, use `AND`:

```bash
tags[go] AND tags[database]

contains["ml"] AND contains["project"]
```

---

## Examples

Return all entries:

```bash
recall all
```

Filter by type:

```bash
recall type[note]

recall type[note,idea]
```

Filter by tags:

```bash
recall tags[go]

recall tags[go,database]
```

Search contents:

```bash
recall contains["project"]

recall contains["ml","project"]
```

Date and time queries:

```bash
recall time > "7 PM"

recall date >= "2026-04-01"

recall len > 100
```

Combine predicates:

```bash
recall time > "7 PM" AND type[note]

recall date >= "2026-04-01" AND contains["project"]

recall tags[go] AND contains["database"]

recall type[note,idea] OR tags[important]
```

---

## Time Formats

Accepted formats:

```text
"7 PM"
"07:30 PM"
"19:30"
"11:00"
```

---

## Date Formats

Accepted formats:

```text
"2026-04-10"
"10/04/2026"
```
