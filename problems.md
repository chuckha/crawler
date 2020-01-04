Problems with this version:

This very small change makes the high-level code much less coupled. The
high-level function no longer depends on the low-level functions. The
dependencies are explicitly passed into CrawlPage.

The with the high-level code is fairly minor now. It's hard to read, but
honestly, since it's not the longest function it's not even that hard to read.
It's easily testable. We no longer need a network to test CrawlPage.

However, we still have a problem with the low-level functions. The low-level
functions don't depend on abstractions. They depend on other low-level modules.