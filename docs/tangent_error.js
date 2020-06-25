/*
DOWNLOADED FROM https://gist.githubusercontent.com/indstronomy/ecea64a8d838d29312092903324f8acd/raw/66e73a40eb5bd2a2ed6162485e6dc14c0ba51caa/tangent_error.js

Copyright 2018 Arun Venkataswamy

Permission is hereby granted, free of charge, to any person obtaining 
a copy of this software and associated documentation files 
(the "Software"), to deal in the Software without restriction, 
including without limitation the rights to use, copy, modify, merge, 
publish, distribute, sublicense, and/or sell copies of the Software, 
and to permit persons to whom the Software is furnished to do so, 
subject to the following conditions:

The above copyright notice and this permission notice shall be 
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, 
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF 
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND 
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS 
BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN 
ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN 
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE 
SOFTWARE.

*/

var reqRisePerMinute_deg = 0.25068448;
var reqRisePerSec_deg = reqRisePerMinute_deg / 60.0;

//Modify these parameters to suit your design
var armLength = 300;
var pitch = 8;
var microSteps = 32;
var fullStepsPerRotation = 200;

var microStepsPerRotation = microSteps * fullStepsPerRotation;
var deviceReqRisePerMinute = 2 * armLength * Math.sin(reqRisePerMinute_deg / 2 * Math.PI / 180.0);
var deviceReqRisePerSecond = deviceReqRisePerMinute / 60.0;
var deviceReqRPM = deviceReqRisePerMinute / pitch;
var devicePulseFreq = deviceReqRPM * microStepsPerRotation / 60.0;

//console.log(deviceReqRisePerMinute,deviceReqRPM,devicePulseFreq);

var oDeviationInDistanceTravelled = 0;
var ot = t = 0;
var correctedDistanceAccumulator = 0;

var arrayString = '';

/*
This is the "bias" time in seconds. The barn door at 
it's home position will have an offset (around 20-30mm). 
That is, it would already be open. It would not be at 0°

x = offset in mm
L = arm length in mm
offset time = (x * 60) / rise per minute
			= (x * 60) / (2 * L * sin(angular rise per min / 2))
			= (x * 60) / (2 * L * sin(0.00437/2)) Note:0.00437 is in radians
*/

var offsetTime = 1371; 

for (t = offsetTime; t <= offsetTime+(3 * 60 * 60); ot = t, t += 5) {

	oDeviationInDistanceTravelled = deviationInDistanceTravelled;

	var expectedAngleDeg = t * reqRisePerMinute_deg/60.0;
	var distanceTravelled = devicePulseFreq * t / microStepsPerRotation * pitch;
	var actualAngleDeg = 2 * Math.asin(distanceTravelled / 2 / armLength) * 180 / Math.PI;
	var deviationDeg = actualAngleDeg - expectedAngleDeg;
	var deviationInDistanceTravelled = 2 * armLength * Math.sin(deviationDeg / 2 * Math.PI / 180);
	var correctedDistanceToTravel = distanceTravelled - deviationInDistanceTravelled;
	var projectedAngleDeg = 2 * Math.asin(correctedDistanceToTravel / 2 / armLength) * 180 / Math.PI;
	var projectedDeviationDeg = projectedAngleDeg - expectedAngleDeg;

	var deltaDeviationInDistanceTravelled = deviationInDistanceTravelled - oDeviationInDistanceTravelled;
	var pulsesForDelta = deltaDeviationInDistanceTravelled / pitch * microStepsPerRotation;
	var pulsesInTimeSegment = ( t - ot ) * devicePulseFreq;
	var correctedPulsesInTimeSegment = pulsesInTimeSegment - pulsesForDelta;
	var correctedDevicePulseFreq = correctedPulsesInTimeSegment / ( t - ot );

	if (t > 0)
		correctedDistanceAccumulator += correctedDevicePulseFreq * ( t - ot ) / microStepsPerRotation * pitch;

	var correctedDevicePulseTimePeriodMicros = 1000000 / correctedDevicePulseFreq;
	correctedDevicePulseTimePeriodMicros /= 2; //Square wave inversions require double speed

	// console.log(
	// 	t,
	// 	distanceTravelled.toFixed(2),
	// 	expectedAngleDeg.toFixed(2),		
	// 	actualAngleDeg.toFixed(2),
	// 	deviationDeg.toFixed(2),
	// 	deviationInDistanceTravelled.toFixed(2),
	// 	correctedDistanceToTravel.toFixed(2),
	// 	correctedDistanceAccumulator.toFixed(2),
	// 	deltaDeviationInDistanceTravelled.toFixed(4),'*',
	// 	pulsesForDelta.toFixed(2),
	// 	correctedDevicePulseFreq.toFixed(2),
	// 	correctedDevicePulseTimePeriodMicros.toFixed(0),'*',
	// 	projectedAngleDeg.toFixed(2),
	// 	projectedDeviationDeg.toFixed(4)
	// );

	if (t != 0) arrayString += ',';
	arrayString += correctedDevicePulseTimePeriodMicros.toFixed(0);
}

console.log(arrayString);

