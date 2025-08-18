# Hexlet Basics Contribution Guide

Welcome! We are glad you are interested in the project. This guide will help you contribute.

## 🧩 Repository structure

The project consists of several key components:
1. **Main repository (hexlet-basics)** - contains educational materials:
- `/<language>/` - directories for each language (Python, JS, etc.)
- `/theory/` - theoretical lessons in Markdown
- `/exercises/` - practical tasks
- `/tests/` - tests for checking solutions
2. **Auxiliary repositories**:
- [hexlet-basics-ru](https://github.com/hexlet-basics/hexlet-basics-ru) - Russian version
- [codebattle](https://github.com/hexlet-codebattle) - solution checking system

## 🚀 How to participate

### For beginners
1. Start with simple tasks:
- Correcting typos in lessons
- Improving wording
- Updating documentation
2. Look for issues with the tags:
- `good first issue`
- `documentation`
- `bug`

### For experienced contributors
1. Working on complex issues:
- Developing new exercises
- Refactoring tests
- Optimizing infrastructure
2. Reviewing pull requests

### Step-by-step process:
1. Find an issue in [Issues](https://github.com/hexlet-basics/hexlet-basics/issues) or create a new one
2. Fork the repository
3. Create a branch for your issue:
```bash
git checkout -b fix/typo-python-variables
```
4. Make changes (following the project style)
5. Write tests if necessary
6. Commit changes:
```bash
git commit -m "docs: fix typo in Python variables lesson"
```
7. Submit a Pull Request to the `main` branch
