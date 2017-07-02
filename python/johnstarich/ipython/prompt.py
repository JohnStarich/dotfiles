from IPython import get_ipython
from IPython.terminal.prompts import Prompts, Token
from lazy import lazy
from prompt_toolkit.auto_suggest import AutoSuggestFromHistory
import os
import platform


class StarshPrompts(Prompts):
    @lazy
    def home(self) -> str:
        return os.environ.get('HOME')

    @lazy
    def username(self) -> str:
        return os.getlogin()

    @lazy
    def hostname(self) -> str:
        return platform.node()

    @lazy
    def prompt_char(self) -> str:
        if os.geteuid() == 0:
            return "#"
        else:
            return "$"

    def cwd(self) -> str:
        """
        Get current working directory and
          collapses home directory to ~.
        Currently gets absolute path.
        TODO Look into changing this to use
          os.path.normpath(os.path.join(cwd, newpath))
        """
        cwd = os.getcwd()
        if cwd.startswith(self.home):
            return cwd.replace(self.home, "~", 1)
        return cwd

    def cwd_basename(self) -> str:
        cwd = self.cwd()
        base = os.path.basename(cwd)

        if base == '':
            return cwd
        else:
            return base

    @property
    def exit_status(self) -> int:
        return self.shell.user_ns.get('_exit_code', 0)

    def git_status() -> str:
        pass

    @lazy
    def _space_token(self):
        return (Token, " ")

    @lazy
    def _auto_suggest(self):
        return AutoSuggestFromHistory()

    def in_prompt_tokens(self, cli=None):
        if cli is not None:
            cli.application.buffer.auto_suggest = self._auto_suggest
        return [
            (Token.Prompt if self.exit_status == 0
                else Token.Generic.Error, '➜'),
            self._space_token,
            self._space_token,
            (Token.Name.Namespace, self.cwd_basename()),
            self._space_token,
        ]

ip = get_ipython()
ip.prompts = StarshPrompts(ip)
del ip
