# ArbitrageClient

![Linting and Tests](https://github.com/robinmuhia/arbitrageClient/actions/workflows/ci.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github.com/robinmuhia/arbitrageClientactions/badge.svg?branch=main)](https://coveralls.io/github.com/robinmuhia/arbitrageClient?branch=main)

## Description

The project implements the `Clean Architecture` advocated by
Robert Martin ('Uncle Bob').

### Points to note

- Interfaces let Go programmers describe what their package provides–not how it does it. This is all just another way of saying “decoupling”, which is indeed the goal, because software that is loosely coupled is software that is easier to change.
- Design your public API/ports to keep secrets(Hide implementation details)
  abstract information that you present so that you can change your implementation behind your public API without changing the contract of exchanging information with other services.

For more information, see:

- [The Clean Architecture](https://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html) advocated by Robert Martin ('Uncle Bob')
- Ports & Adapters or [Hexagonal Architecture](http://alistair.cockburn.us/Hexagonal+architecture) by Alistair Cockburn
- [Onion Architecture](http://jeffreypalermo.com/blog/the-onion-architecture-part-1/) by Jeffrey Palermo
- [Implementing Domain-Driven Design](http://www.amazon.com/Implementing-Domain-Driven-Design-Vaughn-Vernon/dp/0321834577)
