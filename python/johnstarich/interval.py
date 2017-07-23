import time


class Interval(object):
    def __init__(self, delay_time: int):
        self.delay_time = delay_time
        self.current_time = 0

    @staticmethod
    def now():
        return time.gmtime().tm_sec

    def should_run(self) -> bool:
        if self.current_time == 0:
            self.current_time = Interval.now()
            return True
        return self.is_done()

    def is_done(self) -> bool:
        timestamp = Interval.now()
        return self.current_time + self.delay_time < timestamp or \
            self.current_time > timestamp

    def start(self) -> int:
        self.current_time = Interval.now()
        return self.current_time
