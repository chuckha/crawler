Problems with this version:

We now have two low-level functions and one high-level function.

Low-level functions:

- LinkExtractor
- PageReader

High-level function:

- CrawlPage

We might be following the single responsibility principle a bit better now but
we are failing the Dependency Inversion Principle:

> High-level modules should not depend on low-level modules. Both should depend on abstractions.
> Abstractions should not depend on details. Details should depend on abstractions.
> –Robert C. Martin

Our high-level function depends on our low-level functions.

But the way forward is clear, both should depend on abstractions. Details should
depend on abstractions.

Note:

This is still not a very easy function to test.