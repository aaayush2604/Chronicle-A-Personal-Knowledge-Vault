# Chronicle-A-Personal-Knowledge-Vault

## Purpose

This project is a local, append-only personal knowledge store.  
It is designed to capture short thoughts, ideas, and learnings with minimal friction, while preserving their history over time.

The system prioritizes **durability, correctness, and maintainability** over features and user interface.  
It is also a learning-focused project intended to explore how backend systems store data, handle failures, and evolve safely.

---

## Goals

- Allow quick capture of short thoughts and learnings
- Ensure that saved thoughts are not lost once acknowledged
- Enable searching and browsing of past entries
- Preserve the historical timeline of thoughts

---

## Non-Goals

This project intentionally does **not** aim to provide:

- A graphical user interface (CLI only)
- AI-based tagging or classification
- Cloud sync or multi-device support
- Multi-user or collaboration features
- Heavy performance optimizations or distributed operation

---

## System Overview

The system stores personal knowledge as a sequence of **append-only entries** on disk.  
Each entry represents a single thought recorded at a specific moment in time.

The disk-based log acts as the **source of truth**, while in-memory data structures and indexes are used only to improve retrieval speed and usability.

---

## Core Invariants

The following rules must always hold true:

- Once a thought is acknowledged as saved, it must not be lost
- The append-only log on disk is the system’s source of truth
- In-memory state and indexes can always be rebuilt from disk
- The system must be able to recover to a usable state after a crash

---

## Source of Truth

The append-only log stored on disk is the authoritative source of data.

All other data structures—such as in-memory lists, indexes, or caches—are derived from this log and may be safely discarded or rebuilt if necessary.

---

## Derived Data

The following are considered derived and non-authoritative:

- In-memory collections of entries
- Keyword or theme-based indexes
- Time-based indexes
- Cached or summarized views of data

Loss of derived data must **not** result in loss of stored knowledge.

---

## Reliability Definition

For this system, reliability means:

> Once a thought has been successfully written and acknowledged, it must not be lost—even if the application crashes, is forcefully terminated, or restarted later.

Temporary unavailability is acceptable; silent data loss is not.

---

## Failure Scenarios

The system should be designed with the following failure scenarios in mind:

- Application crashes during a write operation
- Process is terminated by the operating system
- Log file contains a partially written or corrupted entry
- Index data is missing or corrupted
- Disk runs out of available space

These scenarios will be handled incrementally as the system evolves.

---

## Scalability Expectations

This system is intended for **single-user, local usage**.

It should comfortably handle:

- Thousands of entries
- Daily usage over long periods of time

Large-scale, distributed, or high-throughput use cases are explicitly out of scope.

---

## Maintainability Goals

The system should be easy to understand and evolve over time.  
Key maintainability goals include:

- Clear separation between storage and indexing logic
- Ability to rebuild indexes without touching stored data
- Simple, readable code paths
- Minimal hidden assumptions or implicit behavior

---

## Open Questions

The following design questions are intentionally left open and will be addressed iteratively:

- What is the most appropriate on-disk format for entries?
- How should partial or corrupted writes be detected?
- When should snapshots or log compaction be introduced?
- How often should indexes be rebuilt?

---

## Summary

This project focuses on building a **simple, reliable personal knowledge storage system**.  
It prioritizes correctness, durability, and clarity over features and performance.

The design intentionally starts small and evolves based on real usage, observed limitations, and learning goals.
