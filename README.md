# Cato

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/NCUHOME/cato/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/NCUHOME/cato)](https://github.com/NCUHOME/cato/issues)
[![GitHub stars](https://img.shields.io/github/stars/NCUHOME/cato)](https://github.com/NCUHOME/cato/stargazers)

> A powerful code generation toolchain for building stable and maintainable Golang server projects.

**⚠️ Important Notice: Under Active Development**
This project is currently under heavy construction. Breaking changes may be introduced at any time. Please do not use it in production environments.

## ✨ What Can Cato Do?

Cato's capabilities are currently provided through `proto-gen-*` plugins.

### proto-gen-cato

#### Note ⚠️
`proto-gen-cato` requires that your `.proto` file content adheres to the [Protocol Buffers Style Guide](https://protobuf.dev/programming-guides/style/). Non-compliance may result in generated struct field names that are difficult to understand.

This plugin provides different capabilities by using the custom options defined in the [`proto`](https://github.com/NCUHOME/cato/tree/main/proto) directory within various scopes. For usage examples, refer to the implementation in the [cato-example-bms (Simple Book Management System)](https://github.com/NCUHOME/cato-example-bms/tree/main/proto) repository.

The responsibilities of the files in the `proto` directory are as follows:
+ **`extension.proto`**: Entry point for using plugin options.
+ **`db.proto`**: Options for use with relational databases.
+ **`defines.proto`**: Options for controlling file generation.
+ **`http.proto`**: Options for generating HTTP services.
+ **`struct.proto`**: Options for common structures.

**Key Features:**
+ **`cato_package`**: Generates corresponding Go structs from Protobuf messages.
+ **`repo_package`**: If a message represents a Data Object, this generates corresponding repository operation interfaces (e.g., `FindBy`, `UpdateBy`, `DeleteBy`, `Insert`) based on defined keys.
+ **`rdb_repo_package`**: For messages with a repository interface, this generates the concrete implementation for relational databases (Relational Database Repository Instance).

## 🎯 Project Vision

Cato is dedicated to the long-term maintenance of stability and clarity in Golang server projects. The name "Cato" is inspired by the game [CATO: Butterered Cat](https://store.steampowered.com/app/1999520/CATO_Buttered_Cat/), embodying the principle that *"Cats always land on their feet."*

Similarly, for Golang server projects, Cato aims to ensure your projects always "land on their feet" – remaining stable and manageable – while the internal code maintains the fluid grace and flexibility of a cat.

## 🛣️ Roadmap

There's still a lot to be done. Many features and improvements are planned for future releases.
