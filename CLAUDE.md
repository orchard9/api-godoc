# Claude Development Guide

Written by world class golang engineers that pride elegant efficiency and delivering a focused product.

## Memory Management

### Memory Structure
The project maintains five memory files in the `.memory/` directory:

1. **state.md**: Lists every file in the project with a one-line description of its purpose
2. **working-memory.md**: Tracks work progress
   - Recent: 1-3 lines describing completed work
   - Current: 1-3 lines describing work in progress
   - Future: 1-3 lines describing upcoming work
3. **semantic-memory.md**: Contains simple factual statements about the project
4. **vision.md**: Defines the stable long-term vision for the project
5. **tasks.md**: Lists all tasks with 1-3 lines describing what each task is and why it's important

### Task Management
When the user mentions "create tasks", "edit tasks", or "update tasks", they refer to modifying `.memory/tasks.md`. This is the authoritative task list for the project, not any temporary todo tracking systems.

## How to Handle Coding Task Requests

### 0. Gather Requirements
Ask questions until you have 100% confidence in understanding:
- API specification format and structure
- Resource grouping requirements
- Output format preferences
- Edge cases and special handling
- **CRITICAL**: For testing/validation tasks, always investigate existing behavior thoroughly before writing tests. Use available tools to understand what they actually produce.

### 1. Create Project Header
Write a clear header that describes the project's purpose and scope.

### 2. Load Memory and Follow Principles
- First: Load all relevant memory files
- Second: Build elegantly following DRY principles
- Third: Ensure flexibility for future development
- Fourth: Maintain strong adherence to the project vision

### 3. Create 80/20 Tests
When relevant to the task:
- Write simple tests that cover 80% of use cases with 20% effort
- Focus on core functionality first
- Outline expected behavior clearly

### 4. Review Test Coverage
Double-check that tests:
- Follow the 80/20 coverage principle
- Adhere to specifications and vision
- Identify and address any technical debt

### 5. Scaffold Code Structure
Create the file structure with 1-3 line descriptions for each file explaining:
- The file's purpose
- What functionality it provides
- How it fits into the overall architecture

### 6. Validate Architecture
Apply the 80/20 rule to pressure test:
- How code will be called and used
- Whether the architecture supports the requirements
- If changes are needed before implementation

### 7. Implement Code
Write and refine code until:
- All tests pass with "make ci"
- Code meets quality standards
- Implementation matches the design

### 8. User Acceptance Testing
Perform final validation against real-world scenarios:
- Build the binary with `make build`  
- Test CLI functionality with `./api-godoc --help` and `./api-godoc --version`
- Test against actual OpenAPI specifications (Stripe, GitHub, etc.)
- Verify all requirements are met from user's perspective
- Ensure the solution is intuitive and reliable

### 9. Update Memory
After completing the task:
- Update state.md with new files
- Update working-memory.md with completed work
- Add new facts to semantic-memory.md if applicable
- Ensure vision.md remains accurate
- Update `.memory/tasks.md`: mark task as completed, move it to the bottom, and add notes about:
  - How it was implemented
  - What could be improved
  - Any concerns or regressions to watch for

## How to Handle Bug Reports

Written by world class golang engineers that pride elegant efficiency and delivering a focused product.

### 0. Understand the Bug
Ask questions until you have 100% confidence in understanding:
- Exact symptoms and error messages
- Steps to reproduce the issue
- Expected vs actual behavior
- Environmental context (OS, version, API spec, etc.)

### 1. Reproduce the Failure
- Create a minimal test case that demonstrates the bug
- Write a failing test that captures the exact issue
- Ensure the test fails for the right reasons
- Document the reproduction steps clearly

### 2. Identify Root Cause
- Trace the code path that leads to the failure
- Understand why our existing tests didn't catch this
- Identify the specific logic or assumption that's incorrect
- Consider related areas that might have similar issues

### 3. Fix with Test-First Approach
- Start with the failing test from step 1
- Implement the minimal fix that makes the test pass
- Ensure the fix doesn't break existing functionality
- Follow our coding principles: reliable, elegant, efficient

### 4. Improve Test Coverage
Critical step - prevent regression:
- Add comprehensive tests that would have caught this bug
- Test edge cases and boundary conditions
- Add validation that catches similar issues
- Ensure tests fail meaningfully when broken

### 5. Validate the Fix
- Run `make ci` to ensure all tests pass
- Run `make uat` to validate real-world scenarios
- Test the specific reproduction case from step 0
- Verify related functionality wasn't impacted

### 6. Update Documentation
- Update code comments if the bug revealed unclear logic
- Add examples or warnings for edge cases
- Update user documentation if behavior changed
- Document the fix approach for future reference

### 7. Update Memory
After fixing the bug:
- Update working-memory.md with the bug fix details
- Add lessons learned to semantic-memory.md
- Note any architectural improvements made
- Update `.memory/tasks.md` with prevention measures if needed

## Handling "check gh" Command

When the user says "check gh", perform these actions:

1. **Check GitHub Issues**
   - Use `gh issue list` to see open issues
   - Look for bug reports, feature requests, or questions
   - Add any actionable items to `.memory/tasks.md`

2. **Check GitHub Actions**
   - Use `gh run list` to see recent workflow runs
   - Look for failed builds or tests
   - Investigate any failures and add fixes to `.memory/tasks.md`

3. **Update Task List**
   - Add new tasks for any issues found to `.memory/tasks.md`
   - Prioritize based on severity (build failures = high priority)
   - Include issue/run numbers in task descriptions for tracking

## How to Handle Testing and UAT Tasks

When implementing tests for existing functionality, follow this investigation-first approach:

### 1. Investigate Current Behavior First
- Run the actual tool with various inputs to understand real behavior
- Document what the tool actually produces, not what you expect
- Use `./build/tool --help`, sample inputs, different formats
- Capture actual outputs before writing any test expectations
- For UAT tasks, run `make uat` and examine existing test patterns

### 2. Start with Smoke Tests
- Begin with simple "does it run without crashing" tests
- Verify basic functionality works before testing edge cases
- Validate one feature at a time rather than comprehensive suites
- Test the most critical user paths first

### 3. Ground Tests in Reality
- Tests should validate actual behavior, not ideal behavior
- If behavior doesn't match expectations, decide: fix code or fix expectations
- Document any surprising behavior discovered during investigation
- Use actual tool outputs to write accurate test assertions

### 4. Build Complexity Gradually
- Start with basic input/output validation
- Add detailed content validation only after basics work
- Prefer multiple simple tests over complex multi-assertion tests
- Each test should validate one specific aspect of functionality

### 5. Validate Against Real-World Usage
- Test with actual user inputs (real API specifications)
- Verify edge cases that users might encounter
- Ensure error messages are helpful and accurate
- Test all supported output formats and options
