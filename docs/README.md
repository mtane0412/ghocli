# gho Documentation

## Overview

This directory contains documentation for the gho (Ghost Admin API CLI) project.

## Documentation Index

### 📊 [PROJECT_STATUS.md](./PROJECT_STATUS.md)

**Purpose**: Understand the current state of the project

**Contents**:
- Implementation phase progress
- Completed features
- Current project structure
- Test coverage
- Dependencies

**When to read**:
- When you want to check the project status
- When you want to know how much has been implemented

---

### 📋 [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md)

**Purpose**: Understand the overall implementation plan

**Contents**:
- Technology stack
- Implementation plan for all phases
- Goals and implementation details for each phase
- Verification methods
- Development workflow

**When to read**:
- When you want to know what to implement next
- When you want to understand the overall picture of implementation
- Before starting a new phase

---

### 🏗️ [ARCHITECTURE.md](./ARCHITECTURE.md)

**Purpose**: Understand the system architecture

**Contents**:
- Project structure
- Layer composition
- Component design
- Authentication flow
- API request flow
- Error handling
- Test strategy
- Security considerations

**When to read**:
- When you want to understand code design
- When deciding where to implement new features
- During code review

---

### 👨‍💻 [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md)

**Purpose**: Learn development methods

**Contents**:
- Development environment setup
- Development workflow (TDD)
- Coding conventions
- How to write tests
- Quality check methods
- How to add new API resources
- Debugging methods
- Troubleshooting

**When to read**:
- When joining the project for the first time
- Before writing code
- When you want to check how to write tests
- Before running quality checks

---

### 🚀 [NEXT_STEPS.md](./NEXT_STEPS.md)

**Purpose**: Check what to do next

**Contents**:
- Current status
- Task list for the next phase
- How to start implementation
- Reference information
- Implementation notes

**When to read**:
- When you want to check the next task
- When starting a new phase
- When you don't know where to start

---

## How to Read Documentation

### For First-Time Project Contributors

1. Read **PROJECT_STATUS.md** to understand the project status
2. Read **ARCHITECTURE.md** to understand system design
3. Read **DEVELOPMENT_GUIDE.md** to learn development methods
4. Read **NEXT_STEPS.md** to check next tasks

### When Starting a New Phase

1. Check task list in **NEXT_STEPS.md**
2. Check detailed plan in **IMPLEMENTATION_PLAN.md**
3. Check implementation methods in **DEVELOPMENT_GUIDE.md**
4. Start implementation

### During Code Review

1. Check design principles in **ARCHITECTURE.md**
2. Check coding conventions in **DEVELOPMENT_GUIDE.md**

### During Troubleshooting

1. Check troubleshooting section in **DEVELOPMENT_GUIDE.md**
2. Check system structure in **ARCHITECTURE.md**

## Documentation Updates

Documentation must be kept up to date at all times.

### When to Update

| Document | Update Timing |
|----------|--------------|
| PROJECT_STATUS.md | When a phase is completed |
| IMPLEMENTATION_PLAN.md | When plans change |
| ARCHITECTURE.md | When architecture changes |
| DEVELOPMENT_GUIDE.md | When development methods change |
| NEXT_STEPS.md | When phases or tasks are completed |

### Update Procedure

1. Edit the document
2. Add "docs:" prefix to the commit message

```bash
git commit -m "docs: update PROJECT_STATUS.md (Phase 2 complete)"
```

## Feedback

If you have suggestions for improving documentation, please share them via:

1. Create a GitHub Issue
2. Submit a Pull Request
3. Include in commit message

## Quick Reference

### Project Information

- **Project Name**: gho
- **Description**: Ghost Admin API CLI
- **Language**: Go 1.22+
- **CLI Framework**: Kong

### Directory Structure

```
gho/
├── cmd/gho/          # Entry point
├── internal/         # Internal packages
│   ├── cmd/         # CLI commands
│   ├── config/      # Configuration management
│   ├── secrets/     # Keyring integration
│   ├── ghostapi/    # Ghost API client
│   └── outfmt/      # Output formatting
└── docs/            # Documentation
```

### Quality Check Commands

```bash
make test         # Run tests
make type-check   # Run type check
make lint         # Run lint
make build        # Build
```

### Important Links

- [Ghost Admin API Documentation](https://ghost.org/docs/admin-api/)
- [Kong CLI Framework](https://github.com/alecthomas/kong)
- [99designs/keyring](https://github.com/99designs/keyring)
