def _update_function_code_impl(ctx):
    output = ctx.label.name + ".deployed"
    output_file = ctx.actions.declare_file(output)

    args = []
    if ctx.attr.function_name:
        args.extend(["-function-arn", "'%s'" % ctx.attr.function_name])
    if ctx.attr.function_name_prefix:
        args.extend(["-function-name-prefix", "'%s'" % ctx.attr.function_name_prefix])
    args.extend(["-region", "'%s'" % ctx.attr.region])

    if ctx.attr.publish:
        args.append("-publish")

    if ctx.attr.dry_run:
        args.append("-dry-run")

    args.extend(["-zip-file", "'%s'" % ctx.file.zip_file.short_path])

    ctx.actions.expand_template(
        template = ctx.file._lambda_update_wrapper,
        output = output_file,
        substitutions = {
            "${args}": " ".join(args),
            "${update_function_code}": ctx.executable._update_function_code.short_path,
        },
        is_executable = True,
    )
    return [
        DefaultInfo(
            executable = output_file,
            runfiles = ctx.runfiles(
                files = [
                    ctx.executable._update_function_code,
                    ctx.file.zip_file,
                ],
            ),
        ),
    ]

_update_function_code_attr = {
    "function_name": attr.string(
        doc = "Function name or Arn",
        mandatory = False,
    ),
    "function_name_prefix": attr.string(
        doc = "Function name Prefix",
        mandatory = False,
    ),
    "region": attr.string(
        doc = "AWS Region to deploy the code",
        mandatory = True,
    ),
    "zip_file": attr.label(
        doc = "The zip file containing the code",
        mandatory = True,
        allow_single_file = True,
    ),
    "publish": attr.bool(
        doc = "Publish new version right away",
        default = False,
    ),
    "dry_run": attr.bool(
        doc = "Dry-run and not apply the change",
        default = False,
    ),
    "_update_function_code": attr.label(
        default = "//rules_lambda/private/cmd/update-function-code",
        executable = True,
        cfg = "host",
    ),
    "_lambda_update_wrapper": attr.label(
        default = "//rules_lambda/private:lambda_update_wrapper.sh",
        allow_single_file = True,
    ),
}

update_function_code = rule(
    implementation = _update_function_code_impl,
    attrs = _update_function_code_attr,
    executable = True,
)
