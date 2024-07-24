# Contributing to lvm2go

Thank you for considering contributing to lvm2go! Your contributions are valuable and help improve the project for everyone. To get started, please read the following guidelines.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [How to Contribute](#how-to-contribute)
    - [Reporting Bugs](#reporting-bugs)
    - [Suggesting Features](#suggesting-features)
    - [Submitting Code Changes](#submitting-code-changes)
3. [Development Setup](#development-setup)
4. [Style Guide](#style-guide)
5. [Testing](#testing)
6. [License](#license)

## Code of Conduct

By participating in this project, you agree to abide by the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). Please report unacceptable behavior to [jakobmoellerdev@example.com](mailto:jakobmoellerdev@example.com).

## How to Contribute

### Reporting Bugs

If you find a bug, please report it by opening an issue on GitHub. Include the following details:
- A clear and descriptive title
- Steps to reproduce the bug
- Expected and actual behavior
- Any relevant screenshots or logs

### Suggesting Features

We welcome suggestions for new features! To suggest a feature, open an issue on GitHub and include:
- A clear and descriptive title
- A detailed description of the feature and its benefits
- Any relevant mockups or examples

### Submitting Code Changes

1. Fork the repository on GitHub.
2. Create a new branch with a descriptive name (`feature/new-feature`, `bugfix/issue-123`, etc.).
3. Make your changes in the new branch.
4. Ensure your code follows the project's style guidelines and passes all tests.
5. Commit your changes with clear and concise messages.
6. Push your branch to your forked repository.
7. Open a pull request on the main repository and describe your changes.

## Development Setup

To set up a development environment, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/jakobmoellerdev/lvm2go.git
    cd lvm2go
    ```

2. Install the necessary dependencies:
    ```sh
    # Example for Fedora
    dnf install lvm2
    ```

3. Run the library tests:
    ```sh
    # Example for a Node.js project
    sudo go test -v ./...
    ```

For detailed setup instructions, refer to the project's [README.md](README.md).

## Style Guide

Please adhere to the following style guidelines:
- Follow the existing code style and conventions.
- Write clear and concise comments.
- Use meaningful variable and function names.
- Format your code using a linter (e.g., ESLint for JavaScript projects).

## Testing

Ensure that your changes do not break existing functionality by running the tests. Add new tests for any new features or bug fixes.

To run tests:
```sh
sudo go test -v ./...
```

Features can only be accepted if rootful tests (Integration Tests on the CI) pass successfully. Only passing Unit Tests is not enough.

## License

By contributing to lvm2go, you agree that your contributions will be licensed under the [Project License](LICENSE).

---

Thank you for your contributions! If you have any questions or need further assistance, feel free to open an issue or contact the project maintainers.
