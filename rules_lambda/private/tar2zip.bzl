def _tar2zip_impl(ctx):
    output = ctx.label.name + ".zip"
    output_file = ctx.actions.declare_file(output)

    args = ctx.actions.args()
    args.add("-output", output_file.path)
    args.add("-input", ctx.file.input.path)

    if ctx.attr.compress:
        args.add("-compress")

    ctx.actions.run(
        mnemonic = "Tar2Zip",
        inputs = [ctx.file.input],
        executable = ctx.executable._tar2zip_binary,
        arguments = [args],
        outputs = [output_file],
        use_default_shell_env = True,
    )

    return [
        DefaultInfo(
            files = depset([output_file]),
        ),
    ]

_tar2zip_attr = {
    "input": attr.label(
        doc = "Input file",
        mandatory = True,
        allow_single_file = True,
    ),
    "compress": attr.bool(
        doc = "Enable ZIP compression, this has an impact in performance",
        default = False,
    ),
    "_tar2zip_binary": attr.label(
        default = "//rules_lambda/private/cmd/tar2zip",
        executable = True,
        cfg = "host",
    ),
}

tar2zip = rule(
    implementation = _tar2zip_impl,
    attrs = _tar2zip_attr,
)
