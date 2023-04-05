def _targz2tar_impl(ctx):
    output = ctx.label.name + ".tar"
    output_file = ctx.actions.declare_file(output)

    args = ctx.actions.args()
    args.add("-output", output_file.path)
    args.add("-input", ctx.file.input.path)

    ctx.actions.run(
        mnemonic = "TarGz2Tar",
        inputs = [ctx.file.input],
        executable = ctx.executable._targz2tar_binary,
        arguments = [ args ],
        outputs = [output_file],
        use_default_shell_env = True,
    )

    return [
        DefaultInfo(
            files = depset([output_file])
        )
    ]

_targz2tar_attrs = {
    "input": attr.label(
        doc = "Input file",
        mandatory = True,
        allow_single_file=True,
    ),
    "_targz2tar_binary": attr.label(
        default = "@//cmd/targz2tar",
        executable = True,
        cfg = "host"
    ),
}

targz2tar = rule(
    implementation = _targz2tar_impl,
    attrs = _targz2tar_attrs
)
