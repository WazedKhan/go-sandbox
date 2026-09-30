# Chapter 2: Pre-declared Types and Declaration

## The Predeclared Types
- Build in types are knows as pre-declared types, ex: boolean, floats and strings

`ex`: The language categorizes its data types into four major groups: Basic Types, Aggregate Types, Reference Types, and Interface Types.

## The Zero Value
Any variable declared in Go without explicit initial value is automatically assigned its zero value

`ex`: `var age int` here age will have value 0. by explicit initial value we avoids bugs like:

- In C and C++, local variables (automatic variables) declared without an explicit initial value contain whatever leftover data happens to be sitting in that memory location. This is known as indeterminate or garbage data.
```cpp
// C++ Danger Example
int* ptr; // Points to a random memory address (e.g., 0x7fff65a3)

if (ptr != nullptr) { // Evaluates to true because it's not null!
    *ptr = 42;        // ❌ Crash or silent memory corruption
}

```

### Explicit Type Conversion
- Automatically converting from one to another when needed is called `automatic type promotion`, Go doesn't allow automatic type promotion
- 
