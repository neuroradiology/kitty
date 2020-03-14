#!/usr/bin/env python
# vim:fileencoding=utf-8
# License: GPLv3 Copyright: 2020, Kovid Goyal <kovid at kovidgoyal.net>


from typing import Sequence


class CLIOptions:
    pass


LaunchCLIOptions = AskCLIOptions = ClipboardCLIOptions = DiffCLIOptions = CLIOptions
HintsCLIOptions = IcatCLIOptions = PanelCLIOptions = ResizeCLIOptions = CLIOptions
ErrorCLIOptions = UnicodeCLIOptions = RCOptions = CLIOptions


def generate_stub() -> None:
    from .cli import parse_option_spec, as_type_stub
    from .conf.definition import save_type_stub
    text = 'import typing\n\n\n'

    def do(otext=None, cls: str = 'CLIOptions', extra_fields: Sequence[str] = ()):
        nonlocal text
        text += as_type_stub(*parse_option_spec(otext), class_name=cls, extra_fields=extra_fields)

    do(extra_fields=('args: typing.List[str]',))

    from .launch import options_spec
    do(options_spec(), 'LaunchCLIOptions')

    from .remote_control import global_options_spec
    do(global_options_spec(), 'RCOptions', extra_fields=['no_command_response: typing.Optional[bool]'])

    from kittens.ask.main import option_text
    do(option_text(), 'AskCLIOptions')

    from kittens.clipboard.main import OPTIONS
    do(OPTIONS(), 'ClipboardCLIOptions')

    from kittens.diff.main import OPTIONS
    do(OPTIONS(), 'DiffCLIOptions')

    from kittens.hints.main import OPTIONS
    do(OPTIONS(), 'HintsCLIOptions')

    from kittens.icat.main import options_spec
    do(options_spec(), 'IcatCLIOptions')

    from kittens.panel.main import OPTIONS
    do(OPTIONS(), 'PanelCLIOptions')

    from kittens.resize_window.main import OPTIONS
    do(OPTIONS(), 'ResizeCLIOptions')

    from kittens.show_error.main import OPTIONS
    do(OPTIONS(), 'ErrorCLIOptions')

    from kittens.unicode_input.main import OPTIONS
    do(OPTIONS(), 'UnicodeCLIOptions')

    from kitty.rc.base import all_command_names, command_for_name
    for cmd_name in all_command_names():
        cmd = command_for_name(cmd_name)
        if cmd.options_spec:
            do(cmd.options_spec, cmd.__class__.__name__ + 'RCOptions')

    save_type_stub(text, __file__)


if __name__ == '__main__':
    import subprocess
    subprocess.Popen([
        'kitty', '+runpy',
        'from kitty.cli_stub import generate_stub; generate_stub()'
    ])