from IPython import get_ipython
from IPython.core.magic import register_line_magic
import astunparse
import ast
import sys


@register_line_magic
def parseprint(line: str):
    tree = None
    try:
        tree = ast.parse(line)
    except:
        print("ast's parse raised an exception:", file=sys.stderr)
        raise
    try:
        return print(astunparse.dump(ast.parse(line)))
    except:
        print("astunparse's dump raised an exception:", file=sys.stderr)
        raise
