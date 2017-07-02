from johnstarich.power.parse import raw_power_info
from johnstarich.interval import Interval
from johnstarich.segment import segment, segment_default


power_status_mappings = {
    'discharging': '🔥',
    'charging': '⚡️',
    'finishing charge': '🔋',
    'charged': '🔋',
    'AC attached': '🔌',
}


power_highlight_groups = ['battery_gradient', 'battery']


update_interval = Interval(10)
last_status = ''
last_gradient = 0
last_highlight_groups = power_highlight_groups


def power(pl, **kwargs):
    global last_status, last_gradient, last_highlight_groups
    if not update_interval.should_run():
        return segment(last_status, gradient_level=last_gradient,
                       highlight_groups=last_highlight_groups)

    power = raw_power_info()
    percentage = int(power['percentage'])
    status = power['status']
    time = power['time']
    time = time.replace(" remaining", "")
    if status == 'charged':
        time = 'full'
    elif time == '0:00' or time == '(no estimate)' or time == 'not charging':
        time = '-:--'
    if status in power_status_mappings:
        status = power_status_mappings[status]

    contents = '{s} {p}% ({t})'.format(
        p=percentage,
        s=status,
        t=time,
    )
    update_interval.start()
    last_status = contents
    last_gradient = 100 - percentage
    last_highlight_groups = power_highlight_groups
    return segment(last_status, gradient_level=last_gradient,
                   highlight_groups=last_highlight_groups)
