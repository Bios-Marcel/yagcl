# Concept

The idea of this project is to provide a simple and extensible API for
loading application configuration.

It shouldn't matter whether you want to use JSON, TOML, .env files or other
formats and sources. Each source can be implemented by implementing a source
interface, meaning that anyone can implement a custom source and make it
available to the community without having to get it accepted into the official
yagcl repositories, which is an option too if the standards are met.

By default, this library provides some source implementations. More on that
in chapter [Modules](./modules.md).

Furthermore, the configuration of your available settings is done via structs
and tags on its fields. All of this is as generic as possible, which results
in some rules that modules have to follow.
