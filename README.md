# Scry

Scry serves a server health dashboard over ssh. 

It is inspired by the [Talos dashboard](https://docs.siderolabs.com/talos/v1.13/deploy-and-manage-workloads/interactive-dashboard) 
but is designed to run on any generic linux server, without requirng a custom GRPC client 
application. 

Scry also contains a plugin system that allows for users to create their own data collection 
functions and dashboard pages that will get automaticlly loaded into the scry application.

## Installation

Scry is intended to be run as a container that mounts the host namespace, so that it can 
collect information about the server, and not just the container sandbox that it runs in.

This is an inherently priveleged set of capabilities, and it should be secured appropriately.
Scry is intended to be run as an alternative to an openssh server on an immutable os, with 
no ssh server or shell. 

## Scry Systems

Scry consists of two main components. The "facter" engine, which collects system information 
by running collection functions in separate threads, and storing that information in a 
key value cache that is passed to the main TUI interface. The TUI interface is the second 
main component of Scry. It's a dynamically generated TUI interface created with bubbletea 
and served over ssh using wish.

### Fact Engine

### Scry TUI

## Plugin System
