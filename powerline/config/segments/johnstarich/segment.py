

def segment(contents: str,
            highlight_groups: list=None,
            divider_highlight_group: str='background:divider',
            **kwargs) -> list:
    if highlight_groups is None:
        highlight_groups = []
    return [{
        'contents': contents,
        'highlight_groups': highlight_groups + ['date'],
        'divider_highlight_group': divider_highlight_group,
        **kwargs
    }]


def segment_default(contents: str, default_contents: str,
                    **kwargs) -> (str, list):
    if contents is None:
        return (default_contents, segment(default_contents, **kwargs))
    else:
        return (contents, segment(contents, **kwargs))
