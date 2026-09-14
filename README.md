# Chronicle-A-Personal-Knowledge-Vault

This project started as an attempt by me to learn GO

Just building a personal project allowing me to note down my thoughts into different categories and being able to use those thoughts/learnings/questions at the right time. Chronicle is a terminal app over an append-only log: you write an entry, give it a type and some tags, and later pull it back with a small query language.

Everything is stored as plain text in a single log file, one line per record. Nothing is ever rewritten in place, so editing and deleting are appended as new records and the log replays into the current state at startup. The path is shown when the app launches, and you can pick an existing log or point it at a new one.

## Commands

```text
rem / remember   write an entry
recall           find entries
revise           change the type or tags of entries
forget           delete entries, after showing what matched
help             the in app manual
```

Type `help` in the terminal for the full manual, which stays in sync with the behaviour described below.

## CQL (Chronicle Query Language)

`recall`, `forget` and the `where` clause of `revise` all take the same query expressions, letting you select entries by the words they contain, their type, their tags, their id, and the date or time they were written. Predicates combine with `AND` and `OR`, and group with parentheses.

If you hit a bug, kindly let me know by raising a pr on bugs.md

# User Manual

## Remembering Entries

Use `rem` or `remember` to add new entries.

```bash
rem [@type] [#tags...] <text>
remember [@type] [#tags...] <text>
```

The type and the tags may come in **either order**, and everything after them is the entry.

Supported entry types:

- `@note` (`@n`)
- `@idea` (`@i`)
- `@question` (`@q`)
- `@learning` (`@l`)
- `@important` (`@imp`)

If no type is specified, the entry is stored as a note. An unknown type is an error rather than a silent note.

### Examples

```bash
rem Read chapter 3

rem @learning Definition of uniform convergence

rem @idea #project #go Use an append-only log

rem #project #go @idea same thing, markers swapped

remember @question #math Why does this proof work?
```

### Entry text is stored exactly as typed

Nothing after the markers is interpreted, so case, punctuation, symbols and other alphabets all survive:

```bash
rem cost is 50% of budget
rem see http://example.com and/or ask
rem café naïve
```

A `#` or `@` inside the text is just text. Only the markers **before** the entry are read as a type and tags:

```bash
rem #work meeting notes #work
# -> tag: work, text: "meeting notes #work"
```

### Quoting

Wrap the whole entry in quotes to keep it literal. The quotes are removed, and nothing inside them is read as a marker:

```bash
rem "@todo buy milk"
# -> note, text: "@todo buy milk"

rem @idea "Read chapter 3"
# -> idea, text: "Read chapter 3"
```

---

## Recall

Use `recall` to search your entries using query expressions.

```bash
recall <predicate> AND/OR <predicate> ...
```

**Note:** `AND` has higher precedence than `OR`. Use parentheses to group.

---

## Revising

Use `revise` to change the type or tags of existing entries.

```bash
revise [@type] [#tags...] where <predicate>
```

A type replaces the old one, tags are added to the existing tags (duplicates are ignored), and the entry text is left alone. Revising updates the entry's time to now, so an entry's time is when it last changed rather than when it was first written.

```bash
revise @important where id[3,7]

revise @important where tags[deadline]

revise #go where contains["golang"]

revise @idea #later where all
```

---

## Forgetting

Use `forget` to delete entries.

```bash
forget <predicate>
```

It shows what matched and asks before deleting. Answer `y` to delete all of them, `n` to cancel, or a comma separated list of ids to delete only those. The query on its own never deletes anything.

```bash
forget id[3,7]

forget contains["draft"]

forget type[note] AND date < "2025"
```

---

## Supported Predicates

| Predicate | Matches |
| --- | --- |
| `all` | every entry |
| `contains["a","b"]` | entries whose text holds any of these words |
| `type[note,idea]` | entries of any of these types |
| `tags[go,db]` | entries carrying any of these tags |
| `id[3,7]` | entries with these ids, the number shown in `{ }` beside each entry |
| `date <op> "..."` | the day the entry was written |
| `time <op> "..."` | the time of day it was written |
| `len <op> <number>` | how many characters the entry holds |

`contains` matches whole words only, and a handful of very common words (`the`, `and`, `with`, `a` and similar) are not searchable.

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
- `id[...]`

Lists use **OR semantics**.

For example:

```bash
type[note,idea]
tags[go,database]
id[3,7,12]
contains["ml","project"]
```

match entries satisfying **any one** of the values.

To require multiple conditions, use `AND`:

```bash
tags[go] AND tags[database]

contains["ml"] AND contains["project"]
```

---

## Reserved Words

```text
all  and  forget  or  recall  rem  remember  revise  where
```

These belong to the query language, so they cannot appear inside `type[...]` or `tags[...]`. A tag with one of these names can still be written, it just cannot be searched for.

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

Pick entries by id, the number shown in braces next to each entry:

```bash
recall id[3]

recall id[3,7,12]
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

recall (type[note] OR type[idea]) AND len > 100
```

---

## Time Formats

```text
"7 PM"        the 7 PM hour
"07:30 PM"    that minute
"19:30"       24 hour clock
"19:30:45"    to the second
```

## Date Formats

```text
"10-04-2026"  DD-MM-YYYY
"2026-04-10"  YYYY-MM-DD
"04-2026"     a whole month
"2026-04"     a whole month
"2026"        a whole year
```

Both `-` and `/` work, and the four digit year tells Chronicle which end is which, so `10/04/2026` and `2026-04-10` are the same day.

## How Date and Time Compare

A value covers a span as wide as the detail you give. `"7 PM"` covers 19:00:00 to 19:59:59, `"2026"` covers the whole year.

```text
=    inside the span          time = "7 PM"   matches 7:30 PM
!=   outside the span
>    after the span starts    time > "7 PM"   matches 7:30 PM
>=   at or after the start
<    before the span starts   time < "7 PM"   excludes 7:30 PM
<=   at or before the end     time <= "7 PM"  matches 7:30 PM
```

Things to watch out for:

- `date > "2026"` still matches 2026, since the year starts before its later days. For entries after a whole year use `date >= "2027"`.
- The year must be written in full, since it is what tells day and year apart. `"10-04-26"` is rejected, `"10-04-2026"` is not.
- A value that is not a real date or time is rejected rather than ignored, so `date > "yesterday"` and `time > "25:00"` are errors, not empty results.

---

## Notes on Storage

Entries are never rewritten. A `revise` appends a new record for the same id and a `forget` appends a deletion record, so the log is a full history and the current state is what you get after replaying it. That means the file only grows, and the ids you see stay stable for the life of an entry.

---

## Session Commands

```text
help      show the in app help
version   show the Chronicle version
clear     clear the screen
index     print the search index built at startup
exit      leave Chronicle
quit      leave Chronicle
```
