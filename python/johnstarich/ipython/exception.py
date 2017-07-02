from IPython import get_ipython
from IPython.terminal.prompts import Token
from prompt_toolkit.shortcuts import print_tokens
import os
import re
import shutil
import traceback


def exception_handler(self, exc_type, exc, tb, tb_offset=None):
    verbose = 'STARSH_DEBUG' in os.environ
    self.showtraceback(exception_only=not verbose)
    self.user_ns['_exit_code'] = 1
    return traceback.format_list(traceback.extract_tb(tb))
    # self.showtraceback(exception_only=True)
    # print_tokens([
    #     (Token.Name.Exception, exc_type.__name__ + ": "),
    #     (Token.Error, str(exc)),
    #     (Token.ZeroWidthEscape, "\n"),
    # ])
    # return ['{}: {}\n'.format(exc_type.__name__, exc)]


ip = get_ipython()
ip.set_custom_exc((Exception,), exception_handler)
