load("@//rules_lambda/private:tar2zip.bzl", _tar2zip = "tar2zip")
load("@//rules_lambda/private:lambda.bzl", _update_function_code = "update_function_code")

tar2zip = _tar2zip
update_function_code = _update_function_code
