# Modules

This project has been split into multiple modules. Each modules has its own
repository and provides a golang module. The main repository `/yagcl` is what
all modules have to depend on.

It's the repository that defines which interface has to be implemented, which
errors should be returned in certain situations and even how certain tags
should be treated. While not all of these things are caught by the compiler, it
is still important the rules are followed, as a consistent API usage can't be
guaranteed otherwise. More on this topic in
[Writing your own module](./own-modules.md).
