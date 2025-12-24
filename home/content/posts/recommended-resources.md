+++
title = "Recommended Resources"
date = "2025-12-24"
description = "A collection of share-worthy technical readings"
+++

# Recommended Resources

## Generally Useful:

- [Let’s Talk about Logging](https://dave.cheney.net/2015/11/05/lets-talk-about-logging): log an error or return an error, but do not do both
    - Simpler: https://github.com/uber-go/guide/blob/master/style.md#handle-errors-once
- [Julia Evan’s Debugging Manifesto](https://jvns.ca/blog/2022/12/08/a-debugging-manifesto/): a great general debugging guide. In particular, I’m a fan of “trust nobody and nothing”
- [UTC Is Enough For Everyone, Right?](https://zachholman.com/talk/utc-is-enough-for-everyone-right): a fun explainer on how time is really hard
- [Notes on Structured Concurrency, or Go Statement Considered Harmful](https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/): a post written to be “incendiary”, but really makes a good case for proper ownership and lifetime management of (concurrent) resources.
    - Generalized: on resource lifecycle management:
        - if you created a resource, you’re responsible for its lifecycle
        - if you were given a reference to a resource, you are *not* responsible for its lifecycle

## For Go:

- [A Tour of Go](https://go.dev/tour/) for an introduction to the Language
- [Effective Go](https://go.dev/doc/effective_go) covers some good idioms/conventions
- [Go by Example](https://gobyexample.com/) is useful for little snippets of common needs
- [Rethinking Classical Concurrency Patterns](https://www.youtube.com/watch?v=5zXAHh5tJqQ) really, really good talk about how to think about concurrency in Go
- [Don't just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully): errors are part of your API surface; consider not having, or relying upon, sentinel errors.
    - for wrapping them, consider the following from The Go Programming Language:
        
        > *When the error is ultimately handled by the program’s main function, it should provide a clear causal chain from the root problem to the overall failure, reminiscent of a NASA accident investigation:*
        > 
        > 
        >     *genesis: crashed: no parachute: G-switch failed: bad relay orientation*
        > 
    - Simpler: https://github.com/uber-go/guide/blob/master/style.md#error-wrapping
- [Align the happy path to the left](https://medium.com/@matryer/line-of-sight-in-code-186dd7cdea88)—this will also inevitably show up as a PR comment to encourage reducing how indented code may be
    - Simpler: https://github.com/uber-go/guide/blob/master/style.md#reduce-nesting
- [Functional Options for Friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis): how to create variadic functions whose use is obvious at the call-site
- [Package Naming](https://dave.cheney.net/2019/01/08/avoid-package-names-like-base-util-or-common): package names should be descriptive. Utils is not.

## Technically neat

- [How NAT Traversal Works](https://tailscale.com/blog/how-nat-traversal-works)
- [The Error Model](https://joeduffyblog.com/2016/02/07/the-error-model/)

## Brain Breaking:

- [Reflections on Trusting Trust](https://dl.acm.org/doi/pdf/10.1145/358198.358210)
- [PostgreSQL used fsync incorrectly for 20 years](https://archive.fosdem.org/2019/schedule/event/postgresql_fsync/)
