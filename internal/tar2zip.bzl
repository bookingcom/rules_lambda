def _tar2zip_impl(ctx):
    pass

_tar2zip_args = {
    "input": attr.label(doc = "Input file", mandatory = True),
    "compress": attr.bool(doc = "Enable ZIP compression, this has an impact in performance", default = False),
    "_tar2zip_binary": attr.label(default = "@//cmd/tar2zip", executable = True),
}

tar2zip = rule(
    implementation = _tar2zip_impl,
    attrs = _tar2zip_args,
)
