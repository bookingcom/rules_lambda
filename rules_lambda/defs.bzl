load("//rules_lambda/private:tar2zip.bzl", _tar2zip = "tar2zip")
load("//rules_lambda/private:lambda.bzl", _update_function_code = "update_function_code")
load("//rules_lambda/private:targz2tar.bzl", _targz2tar = "targz2tar")

tar2zip = _tar2zip
targz2tar = _targz2tar
update_function_code = _update_function_code
