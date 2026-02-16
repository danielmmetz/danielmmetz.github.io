
+++
title = "Package Advice"
date = "2026-02-16"
description = "Advice on developing packages in Go"
+++

# Advice on developing packages in Go

Before digging into organizing packages, let's start with a very brief primer
on [packages and modules](https://go.dev/ref/mod).
A **package** is a collection of source files that are grouped together.
A **module** is a collection of packages that are versioned together.
We're going to focus on packages, largely leaving modules out of scope.

## What does a package offer?

A package offers two main things: namespacing and an enforced boundary (by limiting access to exported members). To quickly define each:

* **namespace**: outside of a package, package members are referenced with a namespace prefix (e.g. `sync.Once` where `sync` is the name of the package and `Once` is an exported member of the package)
* **exports**: outside of a package, only exported (capitalized) identifiers are accessible; everything else is private (e.g. `sync.Once` is exported, but its field `m` is not)

A package _may_ be used for organization, but it's not the only organizational tool offered. Go also offers file-level organization within a package, as we'll see below.

## So when should one create a package?

I encourage treating package development as an organic process.

I typically start simple and let the package grow. When the code does only one thing, a single `main` is fine. A single big file is easy to scroll around and search through—even a simple `cmd+F` works well, let alone editor-assisted navigation. At this stage, jumping from callsite to definition is straightforward, and making changes across both is easy. Signatures are typically purpose-built for the package's particular use case. 

When `main` starts to have a lot to it like several core encapsulations of responsibilities, that can be a good time to move some of those bundles into their own files within the same package. You'll know it's time when you find yourself jumping around a lot within a single file just to find or change the bits you care about. Moving core types and their closely associated functions to dedicated files restores easy navigation within each file, while the package as a whole still provides a single coherent offering.

The [sync package](https://cs.opensource.google/go/go/+/master:src/sync/) is a good example at this stage. It coherently "provides basic synchronization primitives". Core exported types roughly get their own files. There are some additional files for further organization (like `runtime.go`). Not all types get their own files (e.g. `type noCopy`) and cross-file access is fair game.

Finally, I'll fork out a package when the boundary aspect of packages feels actively helpful. 
This is often, but not exclusively, when I want to reuse some functionality elsewhere. Splitting out a package forces consideration of what should be exported, triggering consideration of the exported API surface. Because a package hides/restricts non-exported members, it more easily allows its users to focus on only its exports and to treat the rest as opaque. Done well, introducing a package reduces the context load of its users.

On the flip side, development that routinely requires making coordinated changes across packages can be a sign that the boundary is not helpful. Similarly, many of the times I've run into circular dependency problems is when what could have been one package was cut into too many, tightly coupled pieces.

## Good Package Design

Unsurprisingly, Dave Cheney offers some good [advice on packages and package naming](https://dave.cheney.net/2019/01/08/avoid-package-names-like-base-util-or-common). He offers good examples for when duplication should be favored over creating a package, when code is better modeled as a single package rather than split, and leaves with the advice: "name your packages after what they provide, not what they contain".

A good package is a coherent offering. Its offerings make sense together and its exports are enough to be useful without unnecessarily revealing internals. It's often useful to anchor around the package name. As a heuristic, a positive signal is when a package name clearly communicates the offering(s) of the package. By contrast, difficulty naming the package can be a signal for poor coherence.

A couple of specific naming pitfalls to avoid:

* Avoid creating `base`, `common`, or `utils` packages—these are named for what they contain, not what they provide.
* Avoid stutter. Outside the package, users see the package name as a namespace, so prefer `foo.Client` over `foo.FooClient`.

Ideally, a package empowers its users, and is developed with mind to user experience. Concretely, there are a few things to emphasize:

Accept dependencies rather than creating them internally. Need access to S3? Accept an S3 client. Want to log? Accept a logger. By doing this, you give levers to your callers that let you more easily adapt to situations you didn't need to anticipate. For example, maybe your consumer actually wants to point against an alternative S3-compatible blob store. Or provide a cache-aware client variant. Or they want to force your logs to have an extra structured attribute. 

Further, these dependencies are best accepted as _interfaces_. Using an interface allows you to clearly communicate the functionality you need (e.g. maybe you only need a `GetObject`-capable client, not the whole of the S3 API) while providing the means for your callers to use alternative implementations. Likewise, `error`s aside, generally return concrete types. Callers can create their own interfaces as they need--there's not a need to force that indirection.

Additionally, reduce surprise by eliminating unnecessary side-effects. Avoid `init()` functions, which cause side-effects at import time that callers can't control, and avoid package-level mutable state including by minimizing package-level `var` declarations. For `const` incompatible types, their static values can be pushed behind a function like below:

```go
func ValidState(candidate State) bool {
    return map[State]bool{
        Foo: true,
        Bar: true,
    }[candidate]
}
```

To sum up: start simple, split into files before splitting into packages, and only introduce a package boundary when it actively helps. When you do, design for your users—name the package after what it provides, accept dependencies as interfaces, return concrete types, and minimize surprise.
