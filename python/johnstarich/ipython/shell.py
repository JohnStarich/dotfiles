from IPython import get_ipython
from IPython.core.interactiveshell import InteractiveShell
from IPython.core.inputtransformer import StatelessInputTransformer
import keyword
import shutil
import os
import re
import subprocess
import shlex


def find_var(shell: InteractiveShell, varname: str) -> dict:
    var_locations = [
        lambda: shell.user_ns,
        lambda: shell.user_global_ns,
        lambda: shell.ns_table['builtin'],
        # lambda: shell.alias_manager.linemagics,
        lambda: shell.magics_manager.magics['line'],
        lambda: shell.magics_manager.magics['cell'],
        # ip.alias_manager.aliases
    ]
    for location_func in var_locations:
        var_location = location_func()
        if varname in var_location:
            return var_location
    return None


_command_regex = re.compile(r"""(?x)
(?:(?<=^) | (?<=;)) # don't match invalid commands
(?:
  (?=(?P<segment> # normal tokens or quote block
    [^`'"\\\n;]+ # normal tokens
    | (?:\\(?:.|\n))+ # escaped character or line continuation
    | (?P<quoteblock> # enclosed quotes
        (?P<quote>[`'"]) # quote types
        (?: \\(?:.|\n) | [^\\] )*?
        (?P=quote)
      )+
  ))
  (?P=segment) # emulate atomic group
)+
(?=$|;)
(?:$|;) # tell lazy quantifier when to stop
""")


def run_shell(line: str):
    ip = get_ipython()
    ip.system_raw(line)


def run_magic(line: str):
    ip = get_ipython()
    line_tokens = line.split(maxsplit=1)
    magic_name = line_tokens[0].lstrip('%')
    if len(line_tokens) > 1:
        return ip.run_line_magic(magic_name, line_tokens[1])
    else:
        return ip.run_line_magic(magic_name, '')


@StatelessInputTransformer.wrap
def shell_transformer(line: str):
    # TODO handle more statements in a single line
    line_tokens = line.split(maxsplit=1)
    if len(line_tokens) == 0:
        return line
    if len(line_tokens) > 1 and line_tokens[1].find('=') == 0:
        # If this is an assignment operation,
        # then don't run as a shell command
        return line
    if line_tokens[0].startswith('!'):
        return line

    ip = get_ipython()
    commands = []
    for match in _command_regex.finditer(line):
        command = match.group(0).strip().rstrip(';')
        if command == '':
            # TODO remove this possibility of whitespace chars only
            continue
        if command.startswith('!'):
            commands.append("__run_shell({})".format(
                repr(command.replace("!", "", 1))))
            continue
        command_tokens = command.split(maxsplit=1)
        command_tokens[0] = command_tokens[0].rstrip(';')
        var_location = find_var(ip, command_tokens[0])
        if var_location is None \
                and not keyword.iskeyword(command_tokens[0]) \
                and shutil.which(command_tokens[0]) is not None:
            commands.append("__run_shell({})".format(repr(command)))
        elif var_location == ip.magics_manager.magics['line']:
            # TODO check config for autocall?
            # TODO cell magics?
            commands.append("__run_magic({})".format(repr(command)))
        else:
            commands.append(command)
    return '; '.join(commands)


ip = get_ipython()
ip.input_splitter.logical_line_transforms.insert(0, shell_transformer())
ip.input_transformer_manager.logical_line_transforms\
    .insert(0, shell_transformer())
ip.user_global_ns['__run_shell'] = run_shell
ip.user_global_ns['__run_magic'] = run_magic
del ip

if 'SHELL' not in os.environ:
    os.environ['SHELL'] = 'starsh'
