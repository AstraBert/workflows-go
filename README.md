# workflows-go

## What is workflows-go?

`workflows-go` is a package that implements event-driven workflows in Go: it was designed to build AI applications and intelligent automation leveraging the scalability and the velocity of Go.

## What are event-driven workflows?

When explaining concepts, it is often a good idea to start from the words - in this case:

- **Workflow**: A workflow can be seen as a set of steps, each of which implements a specific logic and plays a unique role: programmatically speaking, we can see the steps as methods of a class, the workflow. A workflow can be linear (all the steps get executed only once per run) or cyclic (steps can be executed more than once per run).
- **Event:** Event is a very broad concept - you can think about it as a trigger that sets off one specific step. Ideally, you have one input event, many intermediate events and one output event.

As you can see, the very own description of event embeds the idea of *event driven workflow*: since an event is a trigger for a step, the workflow will run to completion only if the correct set of events are emitted by the steps - this also means that, by design, steps need to take an event as input and produce an event as output.

Steps and events are the backbone of a workflow, but there might be other things that can help with designing the workflow itself: one of them is **statefulness**.

Statefulness means that your workflow maintains an internal knowledge of what has been going on in previous steps, and maybe even in previous runs: to put it in just a few words, the workflow has an internal representation (state) of the external context.

## Why Go?

Go is a fascinating language, mainly for four reason:

- **Essential design**: Go has only 25 keywords, as opposed to other backend programming languages like Python 3 (38 keywords), JavaScript (45) Rust (53) or C++ (74).
- **Statical typing**: unlike Python or JS, in Go types stay static, but you don't need to declare them explicitly (you can use the walrus, `:=` , operator). This is extremely handy: if you variable changes types or doesn't behaved as expected, it gets flagged immediately at build time.
- **Compiled**: Go is a compiled language, which means that unlike Python (where building and running are combined when you launch the script), build time and run time are separated. The build process spits out portable executables that people can just launch on their machines. Also, Go was designed so that its compile time would be drastically lower than the one of C, and it is also lower than more modern languages like Rust
- **Ecosystem**: Go has an all-in-one tooling ecosystem, that uses one command (`go`) to do everything, from building to installing packages to testing - unlike Python or JavaScript, where we have many tools like uv, poetry and pip or pnpm, yarn and npm.

To cap it all off, Go was designed for scalable and fast large-scale applications, as well as for distributed systems and high-performance microservices: it is thus the perfect language to build the future of reliable AI automation.

## Installation

> _Installing the package requires Go 1.21+_

To install the package, simply run:

```bash
go get https://github.com/AstraBert/workflows-go
```

Or, if you want to get the package from source, clone the repository:

```bash
git clone https://github.com/AstraBert/workflows-go
cd workflows-go
```

And include the code as a base for your project's workflows!

## Import Patterns

After you installed the package, importing is very easy:

```go
import "github.com/AstraBert/workflows-go"
```

You can then use the package as `workflowsgo`

## Contributing

We welcome contributions! Please read our [Contributing Guide](./CONTRIBUTING.md) to get started.

## License

This project is licensed under the [MIT License](./LICENSE)
