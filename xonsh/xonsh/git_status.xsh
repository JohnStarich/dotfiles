#!/usr/bin/env xonsh

def git_status() -> str:
    """
    Create git status prompt similar to oh-my-zsh's robbyrussel prompt
    """
    branch = $PROMPT_FIELDS['curr_branch']()
    if branch is not None:
        branch_color = $PROMPT_FIELDS['branch_color']()
        return "{BLUE}git:(" + branch_color + branch + "{BLUE}){NO_COLOR}"
    else:
        return None
