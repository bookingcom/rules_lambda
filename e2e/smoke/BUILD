"""Provides a simple way to test your rules as an external workspace.
Add a basic smoke-test target below.
"""

load("@bazel_skylib//rules:build_test.bzl", "build_test")
load("@com_booking_rules_lambda//rules_lambda:defs.bzl", "tar2zip", "targz2tar", "update_function_code")

targz2tar(
    name = "test_tar",
    input = "@test_tar_gz//file",
)

tar2zip(
    name = "test_zip",
    compress = False,
    input = ":test_tar",
)

update_function_code(
    name = "update-lambda",
    dry_run = True,
    function_name = "some_lambda",
    region = "eu-central-1",
    zip_file = ":test_zip",
)


build_test(
    name = "smoke_test",
    targets = [
        ":test_tar",
        ":test_zip",
        ":update-lambda"
    ],
)
