load("@//internal:tar2zip.bzl", _tar2zip = "tar2zip")
load("@//internal:targz2zip.bzl", _targz2zip = "targz2zip")
load("@//internal:lambda.bzl", _update_function_code = "update_function_code")

tar2zip = _tar2zip
targz2zip = _targz2zip
update_function_code = _update_function_code
