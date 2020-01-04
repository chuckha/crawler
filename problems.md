Problems with this version:

We've now arrived at a very nice solution that lets us test virtually every line
of code we have written. Since we can test every line we simultaneously do not
need to. However, there is at least one line of code we cannot easily test. It's
inside the CrawlPage function. The return condition if `ioutil.ReadAll` fails.
However, if we think about this, we only need to test that line of code if for
some reason we really care about the error it's returning. This might actually
be the case. For instance, we might say explicitly if we can't parse the body we
return a specific error. If that's part of our public API we should definitely
test the line of code so that we can't introduce a regression that breaks our
public contract.

If this is the case we can do exactly what we did with the function
`html.Parse`. We put it inside a struct and inject it at the highest level,
`func main()`.

Caveat: main.go is excluded from the testing case. It should be simple enough
that we don't need to test it. But if it's not, we would have to fight Go to
make it testable. One trick here is to move the entire contents of `main` into a
function that returns an int. Instead of panicking, return a non 0 value and
only call that function in `main` and check the return value. This avoids us
having to test `os.Exit` conditions which is annoying to do with Go.