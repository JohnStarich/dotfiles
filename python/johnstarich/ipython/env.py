from IPython import get_ipython
from johnstarich.ipython.shell import find_var
import ast
import collections
import os


class CheckedSequence(collections.MutableSequence):
    # TODO use collections.UserList

    def __init__(self, *args, allowed_types: tuple=(object,)):
        self.allowed_types = allowed_types
        self.list = list()
        print(repr(args))
        self.extend(list(args))

    def check(self, value):
        if not isinstance(value, self.allowed_types):
            raise TypeError("Value is not one of the allowed types {}: {}"
                            .format(self.allowed_types, value))

    def changed(self):
        raise NotImplementedError()

    def __len__(self):
        return len(self.list)

    def __getitem__(self, i):
        return self.list[i]

    def __delitem__(self, i):
        del self.list[i]
        self.changed()

    def __setitem__(self, i, value):
        self.check(value)
        self.list[i] = value
        self.changed()

    def insert(self, i, value):
        if not isinstance(i, int):
            raise TypeError("Insert index must be an integer")
        self.check(value)
        self.list.insert(i, value)
        self.changed()

    def __str__(self):
        print('stringing!')
        return str(self.list)

    def __repr__(self):
        return repr(self.list)


class EnvDirVar(CheckedSequence):
    def __init__(self, name, *args):
        self.name = name
        if len(args) == 1:
            if isinstance(args[0], str):
                args = args[0].split(':')
                print('  is string: ' + str(repr(args)))
            elif isinstance(args[0], list):
                args = args[0]
        print('init: ' + str(repr(args)))
        super(EnvDirVar, self).__init__(*args, allowed_types=(str,))

    def changed(self):
        os.environ[self.name] = ':'.join(self.list)


class EnvVar(object):
    def __init__(self, name: str, value=None):
        self.name = name
        if value is None:
            value = self._prepare_var(name)
        self.value = value

    @staticmethod
    def _is_dir_var(name: str):
        return name.endswith('_DIRS') or name == 'PATH'

    @staticmethod
    def _prepare_var(name: str):
        value = os.environ.get(name)
        if EnvVar._is_dir_var(name):
            print('_prepare_var: ' + str(repr(value)))
            if value is None:
                value = []
            value = EnvDirVar(name, value)
        return value

    def _validate_dirs_var(self, value):
        if isinstance(value, EnvDirVar):
            print("_validate a")
            return value
        elif isinstance(value, (str, list)):
            print("_validate 1: " + str(repr(value)))
            print("_validate 2: " + str(repr(EnvDirVar(value))))
            return EnvDirVar(value)
        else:
            raise TypeError("EnvVars ending in '_DIRS' or equal to 'PATH' "
                            "must only be assigned values of type string, "
                            "list, or EnvDirVar")

    def __setattr__(self, name, value):
        if name in ('name', 'value') and not hasattr(self, name):
            if name == 'value' and self._is_dir_var(self.name):
                value = self._validate_dirs_var(value)
            super(EnvVar, self).__setattr__(name, value)
        elif name == 'name':
            raise NotImplementedError("EnvVar name field is immutable.")
        elif name == 'value':
            if self._is_dir_var(self.name):
                print("1) this is what it is: " + str(value))
                value = self._validate_dirs_var(value)
                print("2) this is what it is: " + str(value))
            super(EnvVar, self).__setattr__(name, value)
            os.environ[self.name] = str(self)
        else:
            raise AttributeError("EnvVar has no attribute '{}'".format(name))

    def append(self, value):
        return self.value.append(value)

    def extend(self, iterable):
        return self.value.extend(iterable)

    def insert(self, i, value):
        self.value.insert(i, value)

    def __contains__(self, item):
        return self.value.__contains__(item)

    def __delitem__(self, key):
        del self.value[key]

    def __eq__(self, other):
        return isinstance(other, EnvVar) \
            and other.name == self.name \
            and other.value == self.value

    def __getitem__(self, key):
        return self.value[key]

    def __iadd__(self, other):
        self.value = self.value.__iadd__(other)
        return self

    # TODO iconcat operator?

    def __len__(self):
        return len(self.value)

    def __setitem__(self, i, value):
        self.value[i] = value

    def __str__(self):
        if isinstance(self.value, EnvDirVar):
            return ':'.join(self.value)
        return str(self.value)

    def __repr__(self):
        return repr(self.value)


def env(var: str):
    return EnvVar(var)


class EnvironmentTransformer(ast.NodeTransformer):
    def _is_valid_env_var(self, var: str):
        return var.isidentifier() and var.isupper() and var.startswith('E_')

    def _is_actual_env_var(self, var: str):
        shell = get_ipython()
        return var in shell.user_ns['__env'] or var in os.environ

    def visit_Name(self, node: ast.AST):
        try:
            if self._is_valid_env_var(node.id) or self._is_actual_env_var(node.id):
                return ast.copy_location(ast.Subscript(
                    value=ast.Name(id='__env', ctx=ast.Load()),
                    slice=ast.Index(value=ast.Str(s=node.id)),
                    ctx=node.ctx,
                ), node)
        except Exception as e:
            # TODO remove this exception handler
            print(e)
            return node
        return node


class EnvironmentManager(collections.defaultdict):
    def __missing__(self, key: str):
        if EnvVar._is_dir_var(key):
            return EnvDirVar(key, os.environ.get(key))
        return os.environ.get(key)

    def __setitem__(self, key: str, value):
        if not isinstance(value, (str, list, EnvDirVar)):
            raise TypeError("EnvironmentManager: str expected, not {}".format(type(value).__name__))
        if isinstance(value, list):
            value = EnvDirVar(value)
        os.environ[key] = str(value)
        super(EnvironmentManager, self).__setitem__(key, value)


ip = get_ipython()
ip.ast_transformers.append(EnvironmentTransformer())
ip.push({
    'env': env,
    '__env': EnvironmentManager(),
})
