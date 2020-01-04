Problems with this version:

We've improved the code to clean up the function signature of CrawlPage and have
added a few types to help support the low-level functions as well as satisfy the
interfaces we have defined.

The problem is still that our low-level functions depend on lower-level
functions.

However, if we want to test literally any line of code in the CrawlPage function
we can easily do so. We must construct linkExtractors and or pageReaders that
will execute the exact line of code we wish to test. Since we *can* test any
code easily we don't need to write tests until something goes wrong. At least,
I'm not going to spend time writing tests to get to 100% test coverage when
everything is probably working any way.