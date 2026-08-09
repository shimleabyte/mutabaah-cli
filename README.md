# mutabaah-cli

A small CLI project built while learning Go through implementation.

Project ini berfungsi sebagai **learning laboratory**: setiap konsep Go yang dipelajari langsung diterapkan ke satu project yang terus berkembang, bukan dipisahkan menjadi contoh-contoh syntax yang berdiri sendiri.

The project itself is intentionally simple. Fokus utamanya adalah **learning Go through building**, bukan membuat production-ready habit tracker.

## Learning Goals

Konsep Go yang akan dipelajari dan diimplementasikan secara bertahap:

* Declarations & types
* Structs
* Functions
* Slices
* Methods & receivers
* Pointers
* Function types
* Function variables
* Closures
* Interfaces
* Error handling
* Testing
* Packages & modularization

## Project Direction

Mutabaah CLI dimulai sebagai tracker sederhana untuk merepresentasikan aktivitas dengan target dan progress.

Current model:

```text
Activity
├── Name
├── Target
└── Progress
```

The project will evolve as new Go concepts are introduced.

## Current Capabilities

* Create `Activity` values
* Calculate completion percentage
* Generate activity summaries
* Store multiple activities in a slice
* Iterate over activities
* Compose functions

Example:

```text
Tadarus: 3/5 (60%)
Puasa: 15/30 (50%)
Sedekah: 7/10 (70%)
```

## Learning Journal

| Stage | Topic                     | Status  |
| ----- | ------------------------- | ------- |
| 01    | Go Declaration Syntax     | Done    |
| 02    | Structs & Struct Values   | Done    |
| 03    | Functions & Return Values | Done    |
| 04    | Function Composition      | Done    |
| 05    | Slices & Iteration        | Done    |
| 06    | Methods & Receivers       | Next    |
| 07    | Pointers                  | Planned |
| 08    | Function Types & Closures | Planned |
| 09    | Interfaces                | Planned |
| 10    | Error Handling            | Planned |
| 11    | Testing                   | Planned |
| 12    | Packages & Modularization | Planned |

Each stage is learned by implementing it directly into the project.

## Learning Notes

Detailed explanations, mistakes, mental models, and experiments are kept separately from this README.

The README is intentionally kept short so it remains useful as a project index rather than becoming a textbook.

## Reference

Primary reference for the first stage:

**Rob Pike — Go's Declaration Syntax (2010)**
**ChatGPT**, **Notebooklm**, as a project coach

The goal is not merely to reproduce the examples from the reference, but to understand the reasoning behind Go's syntax and apply it through implementation.
