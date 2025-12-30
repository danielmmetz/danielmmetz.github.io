+++
title = "A Guide to Code Review"
date = "2025-12-24"
description = "How I think about code reviews"
+++

# A Guide to Code Review

**Document Goals:**

This is generally documents how I do, and think about, code reviews. Please consider it to be *a* resource you can use as you think about, and conduct, your own code reviews. I would be delighted if it helps you up-level the quality of reviews you’re providing.

That said, please do not take this as a “Law of Code Review”, nor even a comprehensive guide. There are very well other ways to go about code reviews and other things to consider.

## On What to Look For {#on-mystery}

1. Start by understanding **the goal**. Consider the following:
    - **What problem is the pull request intending to address?** **Is that the right problem?** If it’s not, spend your energy aligning onto the right problem(s).
    - Now that we’re aligned this is the right problem to solve, **is this the right solution?** Take a moment to **consider how else the problem can be addressed** (or even just mitigated). If you’re not convinced this is the right way to address the problem in focus, spend your energy aligning onto the right solution(s).
2. Now that we’re aligned on both this being the right problem to address and the right solution to that problem, **verify the change does as intended**. It’s not *that* uncommon for a PR to either solve a slightly different problem than claimed, or push forward an approach that is (subtlety, but significantly) different than described. You can consider **this as a pass checking for *high-level correctness*** (aka claim consistency).
3. **Evaluate the code changes for correctness**. This is now a *lower-level pass*. Consider the edge-cases in the flow. **What data will flow through the system?** **What assumptions is the code making, implicitly or explicitly?** Are those assumptions reasonable? **Are errors (or exceptions) properly considered and handled?**
    - Ideally, this correctness is sufficiently obvious. While testing is generally important, **ask for specific, additional tests to provide confidence for that which isn’t obvious**. And consider asking for them to **cement-in important behavior to prevent regressions**.
4. If concurrent reads/writes are at play, give another pass specifically focused on potential concurrency bugs. Rather than expand on specific “look for this or that”s here, I’ll just leave the following two statements:
    - a) channels are an incredibly low-level primitive whose direct use should usually be avoided—you’re often better off using a `sync.WaitGroup`, `sync.ErrGroup`, or `sync.Mutex`.
    - b) [Rethinking Classical Concurrency Patterns](https://www.youtube.com/watch?v=5zXAHh5tJqQ) is a really good talk that helps clarify that concurrency should be an implementation detail, and not something exposed to the caller.
5. **Consider the end user’s experience**. Will the experience feel appropriately snappy? Will they get feedback in a timely manner and will the messaging be sufficiently clear? Often, this means:
    - are the slow parts (network calls) appropriately parallelized?
    - are UIs and APIs obvious and hard/impossible to misunderstand/misuse?
    - is user-facing text clear and concise?
6. **Consider readability**. Code is read far more often than it’s written. Simple code is less likely to hide a bug than complex code. Simple code is more likely to be correctly read and correctly interpreted by other readers later on.
    - Humans being humans, they’re also likely to copy the patterns and styles they see in the codebase. So know that the code that merges has a reasonable likelihood of proliferating. If it’s good, great! If it’s bad, the impact down the road might be a codebase with a lot more of that.
    - Okay great, but how does this get applied in practice?
    - Before reading the code, **consider how you might approach the solution**. Maybe even consider multiple approaches.
    - **Does the code look as you expected**, based on the problem and proposed solution?
    - **Is the code obvious?** Ideally, it should both obviously do what it says (as implied by the function names, variable names, etc) and be obviously correct.
    - **Does the code adhere to relevant idioms** (language norms)?
7. **Consider debuggability and observability**. How will you know when things go wrong? When they do, will you have the necessary information to debug (and where appropriate, reproduce) the issue?

**The Approval Bar**

To be approved, a pull request should:

- address the right problem,
- with the right solution,
- with an implementation that correctly implements *that* solution to address *that* problem,
- that delivers a good user experience,
- via readable code.

The pull request does not need to be perfect. It does not need to be identical to what you’d have written yourself. However, **you should feel confident about the pull request**. After all, you’re very likely to be responsible for it as part of your team’s product(s) and on-call.

## On Giving Feedback

1. **Constructive suggestions are far more valuable than constructive critiques**. Where possible, try to include an answer with suggestions. Instead of “this is bad” or “take this alternative approach”, prefer to provide enough clarity so that the reader understands precisely what you’re looking for. This can look like:
    - providing a complete code snippet so it’s crystal clear how to implement the suggested change
    - providing several examples of satisfactory answers (like when suggesting a variable rename), both so it’s easy to satisfy the ask, and so that the underlying theme or pattern is clear
    - Very often, these can be teaching moments. With that in mind, cater the level of detail considering your audience.
2. **Motivate the suggestions**. Where possible, explain *why* you’re requesting a particular change. This can look like:
    - linking to an external resource that reinforces a concept
    - linking to a code demo to illustrate a bug
    - adding additional context for how the change fits into changes that may come later
3. Ideally, **be clear about what requests are blocking**. You may find yourself leaving many comments on a PR; some may identify bugs whereas some may be minor stylistic suggestions. In such cases, it can be useful to clarify which comments are those you intend to block approval or merging.
4. **Leave feedback that makes the codebase better**. Such feedback may be non-blocking, and that’s okay. Just because a PR is good enough to ship, doesn’t mean you shouldn’t still share feedback to make it better.

## On Timing

The ideal code review is timely. The best loop is one where there’s little delay between requesting the review of a pull request, the review coming in, and each subsequent iteration.

However, reviews can still be useful even if they’re not timely. **A review can still be valuable even after a pull request has been merged**.

## To Pull Request Authors

This guide applies to you too. **Review your own code before requesting a review** from others. Make your pull request easy to review by ensuring you provide enough context to understand the relevant problems. Consider adding comments (in code, on the pull request, or in its description) to clarify things for your readers, or bring attention to particular focus areas.

Additionally, as an author, your goal should not be “get a stamp”. Instead, consider the review process as one that can spread knowledge, help you (and the codebase) improve, and build wider confidence in the changes being considered.
