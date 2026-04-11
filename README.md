# Chronicle-A-Personal-Knowledge-Vault

This project started as an attempt by me to learn GO

Just building a personal project allowing me to note down my thoughts into different categories and being able to use those thoughts/learnings/questions at the right time. This project currently is just a simple append-only log that lets you enter notes and allows you to pull those thoughts based on the words you might think are used in those entries. If trying to see all the features that this currently has, use help in the terminal. The path to logFile is displayed on the terminal when the project starts. If you want to add a particular path for the logFile, apart from the default, just add CHRONICLE_LOG_FILE_PATH to an env file

## CQL (Chronicle Query Language)

I have added a basic query Language to this, allowing you to query your notes based on the words an entry contains, the entry's type and what date or time the entry was made. I am yet to enter a testing phase, so you might see a lot of bugs, kindly let me know by raising a pr, on the bugs.md

Here is a basic user manual

Use `recall` to search your notes with simple query expressions.

```bash
recall <predicate> AND/OR <predicate> ...
```

Supported Fields

date, time, type, contains

Operators

=, !=, >, <, >=, <=, AND, OR (Keep in mind AND is more binding than OR)

Examples

```bash
recall time > "7 PM"
recall date >= "2026-04-01"

recall type["imp","note"]              # matches any
recall contains ["ml","project"]       # must contain all

recall time > "7 PM" AND type["note"]
recall date >= "2026-04-01" AND contains ["project"]

```

Formats
Time: "7 PM", "07:30 PM", "19:30", "11:00"
Date: "2026-04-10", "10/04/2026"
