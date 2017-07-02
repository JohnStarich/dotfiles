from IPython import get_ipython
from IPython.core.magic import (magics_class, line_magic)
from IPython.core.magics.osm import OSMagics
from johnstarich.ipython.shell import find_var
import keyword
import shutil


@magics_class
class Bashisms(OSMagics):
    @property
    def _exit_code(self) -> int:
        return self.shell.user_ns['_exit_code']

    @_exit_code.setter
    def _exit_code(self, value: int):
        self.shell.user_ns['_exit_code'] = value

    @line_magic
    def echo(self, line: str):
        "Simply print out its received arguments."
        print(line.format(**vars(), **globals()))
        self._exit_code = 0
        return

    @line_magic
    def cd(self, parameter_s=''):
        super(Bashisms, self).cd('-q ' + parameter_s)

    @line_magic
    def which(self, line):
        var_location = find_var(self.shell, line)
        if var_location is not None:
            print(var_location.get(line))
            self._exit_code = 0
            return

        if keyword.iskeyword(line):
            help(line)
            self._exit_code = 0
            return

        ex = shutil.which(line)
        if ex is not None:
            print(ex)
            self._exit_code = 0
            return
        else:
            print('"{}" could not be found on $PATH'
                  .format(line))
            self._exit_code = 1
            return

ip = get_ipython()
ip.register_magics(Bashisms)
del ip
