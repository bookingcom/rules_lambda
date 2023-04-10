# Bazel rules for lambda

## Installation

In your MODULE.bazel you add a line that says

bazel_dep(name = "com_booking_rules_lambda", version = "0.0.1")

We only support bzlmod, we could add non bzlmod dependency if someone is willing
to maintain it, but so far we don't have a need for it. 
