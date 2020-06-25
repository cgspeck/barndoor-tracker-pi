#! /usr/bin/env python
import argparse
import math

from dataclasses import dataclass

output_choices=['verbose', 'cpp']
parser = argparse.ArgumentParser(description='Calculate step rate for astro-tracker, correcting for tangent error.')
parser.add_argument('-m', '--output_mode', choices=output_choices, default=output_choices[0])

args = parser.parse_args()

sidereal_time=23*60*60+56*60+4.0916
degrees_per_second=360 / sidereal_time
degrees_per_minute=degrees_per_second * 60

start_moment=0
end_moment=3*60*60
moment_duration=60

steps_per_rotation=200
rotations_per_cm=10  # this is threads per cm

# describing an isoseleces triangle,
# this is the lenth of the two equal sides 
leg_length_mm=300

@dataclass
class Moment:
    moment: int
    minutes: float = None
    target_angle: float = None
    target_base_len_mm: float = None
    target_distance_moved_mm: float = None
    target_rotations: float = None
    target_steps: float = None
    steps_per_second: float = None

    def do_calculations(self):
        self.minutes = self.moment / 60
        self._calculate_target_angle()
        self._calculate_target_base_len_mm()
        self._calculate_target_distance_moved_mm()
        self._calculate_target_rotations()
        self._calculate_target_steps()
        self._calculate_steps_per_second()
    
    def to_csv(self):
        return f'{self.moment}, {self.target_angle}, {self.target_base_len_mm}, {self.target_steps}, {self.steps_per_second}'

    def _calculate_target_angle(self):
        self.target_angle = degrees_per_second * self.moment
        self._previous_angle = degrees_per_second * (self.moment - moment_duration)
    
    def _calculate_target_base_len_mm(self):
        half_angle_radians = (self.target_angle / 2) * math.pi / 180
        self.target_base_len_mm = math.sin(half_angle_radians) * leg_length_mm * 2
        previous_half_angle_radians = (self._previous_angle / 2) * math.pi / 180
        self._previous_base_len = math.sin(previous_half_angle_radians) * leg_length_mm * 2
    
    def _calculate_target_distance_moved_mm(self):
        self.target_distance_moved_mm = self.target_base_len_mm - self._previous_base_len
    
    def _calculate_target_rotations(self):
        self.target_rotations = (self.target_distance_moved_mm / 10) * rotations_per_cm
    
    def _calculate_target_steps(self):
        self.target_steps = self.target_rotations * 200
    
    def _calculate_steps_per_second(self):
        self.steps_per_second = self.target_steps / moment_duration

res = []

for moment in range(start_moment, (end_moment + moment_duration), moment_duration):
    moment_ = Moment(moment=moment)
    moment_.do_calculations()

    res.append(moment_)

def output_csv():
    print('moment, target_angle, target_base_len_mm, target_steps, steps_per_second')
    for r in res:
        print(r.to_csv())


def _calculate_minute_to_steps_per_second():
    minute_to_steps_per_second = []
    start_minute = 1
    end_minute = (end_moment + moment_duration) / 60

    for minute in range(start_minute, int(end_minute)):
        current = None
        previousVal = None
        previousDiff = None

        for r in res:
            diff = abs(minute - r.minutes)
            if diff == 0:
                current = r.steps_per_second
                break
            
            if previousDiff is not None:
                if previousDiff < diff:
                    current = previousVal
                    break

            previousDiff = diff
            previousVal = r.steps_per_second
        
        minute_to_steps_per_second.append(current)

    return minute_to_steps_per_second

def output_cpp():
    print('''
#ifndef lookupTable_h
#define lookupTable_h
''')
    # minute_to_steps_per_second
    minute_to_steps_per_second = _calculate_minute_to_steps_per_second()
    var1_array_str = ', '.join(["%s" % x for x in minute_to_steps_per_second])
    var1_str = f'''
float MINUTE_TO_STEPS_PER_SECOND[{len(minute_to_steps_per_second)}] = {{{var1_array_str}}};
    '''
    print(var1_str)
    print('#endif')

if args.output_mode == 'verbose':
    output_csv()

if args.output_mode == 'cpp':
    output_cpp()