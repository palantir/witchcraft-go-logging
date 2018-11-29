witchcraft-error-go
===================
`witchcraft-error-go` defines the `werror` package, which provides an implementation of the `error` interface that
stores safe and unsafe parameters and has the ability to specify another error as a cause.

Associating structured safe and unsafe parameters with an error allows other infrastructure such as logging to make
decisions about what parameters should or should not be extricated.

TODO:
* Provide example usage and output in README
