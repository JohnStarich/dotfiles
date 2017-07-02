#!/usr/bin/env python3

# Add custom python site dir
from os.path import dirname, realpath
import site
import os

dotfiles_dir = dirname(dirname(realpath(__file__)))
site.addsitedir(os.path.join(dotfiles_dir, "python"))

# c = the current IPython config dict
if 'c' not in vars() and 'c' not in globals():
    print("WARNING: Parsing IPython config script outside of IPython.")
    from traitlets.config import Config
    c = Config()

c.TerminalInteractiveShell.confirm_exit = False
c.TerminalInteractiveShell.separate_in = ''
c.TerminalIPythonApp.display_banner = False
c.InteractiveShell.editor = 'vim'
c.InteractiveShell.automagic = True
c.TerminalInteractiveShell.term_title = True
# autocall doesn't work with `from ... import ...'
#c.InteractiveShell.autocall = 1
c.InteractiveShellApp.extensions = [
    # 'autoreload',
    # 'Cython',
    'johnstarich.ipython.ast',
    'johnstarich.ipython.bashisms',
    'johnstarich.ipython.env',
    'johnstarich.ipython.exception',
    'johnstarich.ipython.prompt',
    'johnstarich.ipython.shell',
]
c.AliasManager.user_aliases = [
]
c.InteractiveShell.ast_transformers = [
]
#c.BaseIPythonApplication.verbose_crash = False
