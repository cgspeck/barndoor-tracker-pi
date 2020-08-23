EESchema Schematic File Version 4
EELAYER 30 0
EELAYER END
$Descr A4 11693 8268
encoding utf-8
Sheet 1 1
Title ""
Date ""
Rev ""
Comp ""
Comment1 ""
Comment2 ""
Comment3 ""
Comment4 ""
$EndDescr
$Comp
L Connector:Raspberry_Pi_2_3 J1
U 1 1 5F41CD1F
P 1500 2200
F 0 "J1" H 1500 3681 50  0000 C CNN
F 1 "Raspberry_Pi_2_3" H 1500 3590 50  0000 C CNN
F 2 "" H 1500 2200 50  0001 C CNN
F 3 "https://www.raspberrypi.org/documentation/hardware/raspberrypi/schematics/rpi_SCH_3bplus_1p0_reduced.pdf" H 1500 2200 50  0001 C CNN
	1    1500 2200
	1    0    0    -1  
$EndComp
$Comp
L Connector_Generic:Conn_02x20_Odd_Even J6
U 1 1 5F41E96B
P 5200 2100
F 0 "J6" H 5250 3217 50  0000 C CNN
F 1 "Conn_02x20_Odd_Even" H 5250 3126 50  0000 C CNN
F 2 "" H 5200 2100 50  0001 C CNN
F 3 "~" H 5200 2100 50  0001 C CNN
	1    5200 2100
	1    0    0    -1  
$EndComp
$Comp
L Sensor_Motion:LSM9DS1 U1
U 1 1 5F41FB8E
P 1300 4600
F 0 "U1" H 1300 3711 50  0000 C CNN
F 1 "LSM9DS1" H 1300 3620 50  0000 C CNN
F 2 "Package_LGA:LGA-24L_3x3.5mm_P0.43mm" H 2800 5350 50  0001 C CNN
F 3 "http://www.st.com/content/ccc/resource/technical/document/datasheet/1e/3f/2a/d6/25/eb/48/46/DM00103319.pdf/files/DM00103319.pdf/jcr:content/translations/en.DM00103319.pdf" H 1300 4700 50  0001 C CNN
	1    1300 4600
	1    0    0    -1  
$EndComp
$Comp
L Connector:Conn_01x04_Female J12
U 1 1 5F42078D
P 8850 1400
F 0 "J12" H 8878 1376 50  0000 L CNN
F 1 "Conn_01x04_Female" H 8878 1285 50  0000 L CNN
F 2 "" H 8850 1400 50  0001 C CNN
F 3 "~" H 8850 1400 50  0001 C CNN
	1    8850 1400
	1    0    0    -1  
$EndComp
$Comp
L Connector:Conn_01x03_Male J11
U 1 1 5F420F88
P 8650 3550
F 0 "J11" H 8758 3831 50  0000 C CNN
F 1 "Conn_01x03_Male" H 8758 3740 50  0000 C CNN
F 2 "" H 8650 3550 50  0001 C CNN
F 3 "~" H 8650 3550 50  0001 C CNN
	1    8650 3550
	1    0    0    -1  
$EndComp
$Comp
L PS2501-4-A:PS2501-4-A IC1
U 1 1 5F423967
P 6500 4400
F 0 "IC1" H 7000 4665 50  0000 C CNN
F 1 "PS2501-4-A" H 7000 4574 50  0000 C CNN
F 2 "DIP762W60P254L1980H455Q16N" H 7350 4500 50  0001 L CNN
F 3 "https://static6.arrow.com/aropdfconversion/f7a3ab99120737b5601e335c604d7372a3bcd72c/1803056367660517pn10225ej05v0ds.pdf" H 7350 4400 50  0001 L CNN
F 4 "Transistor Output Optocouplers Optocoupler 16-pin DIP" H 7350 4300 50  0001 L CNN "Description"
F 5 "4.55" H 7350 4200 50  0001 L CNN "Height"
F 6 "Renesas Electronics" H 7350 4100 50  0001 L CNN "Manufacturer_Name"
F 7 "PS2501-4-A" H 7350 4000 50  0001 L CNN "Manufacturer_Part_Number"
F 8 "PS2501-4-A" H 7350 3900 50  0001 L CNN "Arrow Part Number"
F 9 "https://www.arrow.com/en/products/ps2501-4-a/renesas-electronics" H 7350 3800 50  0001 L CNN "Arrow Price/Stock"
F 10 "551-PS2501-4-A" H 7350 3700 50  0001 L CNN "Mouser Part Number"
F 11 "https://www.mouser.co.uk/ProductDetail/Renesas-Electronics/PS2501-4-A?qs=qSfuJ%252Bfl%2Fd73DT6GctNDSg%3D%3D" H 7350 3600 50  0001 L CNN "Mouser Price/Stock"
	1    6500 4400
	1    0    0    -1  
$EndComp
$Comp
L Sensor_Temperature:DS18B20 U2
U 1 1 5F424393
P 9200 5500
F 0 "U2" H 8970 5546 50  0000 R CNN
F 1 "DS18B20" H 8970 5455 50  0000 R CNN
F 2 "Package_TO_SOT_THT:TO-92_Inline" H 8200 5250 50  0001 C CNN
F 3 "http://datasheets.maximintegrated.com/en/ds/DS18B20.pdf" H 9050 5750 50  0001 C CNN
	1    9200 5500
	1    0    0    -1  
$EndComp
$Comp
L Connector:AudioJack3 J10
U 1 1 5F426E8C
P 8050 2550
F 0 "J10" H 8032 2875 50  0000 C CNN
F 1 "AudioJack3" H 8032 2784 50  0000 C CNN
F 2 "" H 8050 2550 50  0001 C CNN
F 3 "~" H 8050 2550 50  0001 C CNN
	1    8050 2550
	1    0    0    -1  
$EndComp
$Comp
L Connector:Conn_01x03_Female J9
U 1 1 5F4278A5
P 6800 2700
F 0 "J9" H 6828 2726 50  0000 L CNN
F 1 "Conn_01x03_Female" H 6828 2635 50  0000 L CNN
F 2 "" H 6800 2700 50  0001 C CNN
F 3 "~" H 6800 2700 50  0001 C CNN
	1    6800 2700
	1    0    0    -1  
$EndComp
$Comp
L Device:R R1
U 1 1 5F427E47
P 4000 3600
F 0 "R1" H 4070 3646 50  0000 L CNN
F 1 "4K7" H 4070 3555 50  0000 L CNN
F 2 "" V 3930 3600 50  0001 C CNN
F 3 "~" H 4000 3600 50  0001 C CNN
	1    4000 3600
	1    0    0    -1  
$EndComp
$Comp
L Device:R R2
U 1 1 5F429448
P 4500 3700
F 0 "R2" H 4570 3746 50  0000 L CNN
F 1 "1K" H 4570 3655 50  0000 L CNN
F 2 "" V 4430 3700 50  0001 C CNN
F 3 "~" H 4500 3700 50  0001 C CNN
	1    4500 3700
	1    0    0    -1  
$EndComp
$Comp
L Device:R R4
U 1 1 5F42990F
P 4850 3750
F 0 "R4" H 4920 3796 50  0000 L CNN
F 1 "1K" H 4920 3705 50  0000 L CNN
F 2 "" V 4780 3750 50  0001 C CNN
F 3 "~" H 4850 3750 50  0001 C CNN
	1    4850 3750
	1    0    0    -1  
$EndComp
$Comp
L FQP30N06L:FQP30N06L Q1
U 1 1 5F41F6DD
P 5000 5300
F 0 "Q1" H 5430 5446 50  0000 L CNN
F 1 "FQP30N06L" H 5430 5355 50  0000 L CNN
F 2 "TO254P470X1016X2019-3P" H 5450 5250 50  0001 L CNN
F 3 "https://www.mouser.com/datasheet/2/308/FQP30N06L-1306227.pdf" H 5450 5150 50  0001 L CNN
F 4 "FQP30N06L N-Channel MOSFET, 32 A, 60 V QFET, 3-Pin TO-220 ON Semiconductor" H 5450 5050 50  0001 L CNN "Description"
F 5 "4.7" H 5450 4950 50  0001 L CNN "Height"
F 6 "ON Semiconductor" H 5450 4850 50  0001 L CNN "Manufacturer_Name"
F 7 "FQP30N06L" H 5450 4750 50  0001 L CNN "Manufacturer_Part_Number"
F 8 "FQP30N06L" H 5450 4650 50  0001 L CNN "Arrow Part Number"
F 9 "https://www.arrow.com/en/products/fqp30n06l/on-semiconductor" H 5450 4550 50  0001 L CNN "Arrow Price/Stock"
F 10 "512-FQP30N06L" H 5450 4450 50  0001 L CNN "Mouser Part Number"
F 11 "https://www.mouser.co.uk/ProductDetail/ON-Semiconductor-Fairchild/FQP30N06L?qs=%252Beu4QTTVkDqugBwFplJ%252Bug%3D%3D" H 5450 4350 50  0001 L CNN "Mouser Price/Stock"
	1    5000 5300
	1    0    0    -1  
$EndComp
$Comp
L Connector:Screw_Terminal_01x02 J8
U 1 1 5F420F7F
P 6200 6300
F 0 "J8" H 6280 6292 50  0000 L CNN
F 1 "Screw_Terminal_01x02" H 6280 6201 50  0000 L CNN
F 2 "" H 6200 6300 50  0001 C CNN
F 3 "~" H 6200 6300 50  0001 C CNN
	1    6200 6300
	1    0    0    -1  
$EndComp
$Comp
L Connector:USB_A J3
U 1 1 5F421A3C
P 3650 6550
F 0 "J3" H 3707 7017 50  0000 C CNN
F 1 "USB_A" H 3707 6926 50  0000 C CNN
F 2 "" H 3800 6500 50  0001 C CNN
F 3 " ~" H 3800 6500 50  0001 C CNN
	1    3650 6550
	1    0    0    -1  
$EndComp
$Comp
L Switch:SW_DPDT_x2 SW1
U 2 1 5F4251C6
P 3800 4650
F 0 "SW1" H 3800 4935 50  0000 C CNN
F 1 "SW_DPDT_x2" H 3800 4844 50  0000 C CNN
F 2 "" H 3800 4650 50  0001 C CNN
F 3 "~" H 3800 4650 50  0001 C CNN
	2    3800 4650
	1    0    0    -1  
$EndComp
$Comp
L Switch:SW_DPDT_x2 SW1
U 1 1 5F425BF6
P 3750 5150
F 0 "SW1" H 3750 5435 50  0000 C CNN
F 1 "SW_DPDT_x2" H 3750 5344 50  0000 C CNN
F 2 "" H 3750 5150 50  0001 C CNN
F 3 "~" H 3750 5150 50  0001 C CNN
	1    3750 5150
	1    0    0    -1  
$EndComp
$Comp
L Connector:Screw_Terminal_01x03 J5
U 1 1 5F4268C0
P 5050 7100
F 0 "J5" H 5130 7142 50  0000 L CNN
F 1 "Screw_Terminal_01x03" H 5130 7051 50  0000 L CNN
F 2 "" H 5050 7100 50  0001 C CNN
F 3 "~" H 5050 7100 50  0001 C CNN
	1    5050 7100
	1    0    0    -1  
$EndComp
$Comp
L MCU_Module:Arduino_Nano_v3.x A1
U 1 1 5F427678
P 1200 6750
F 0 "A1" H 1200 5661 50  0000 C CNN
F 1 "Arduino_Nano_v3.x" H 1200 5570 50  0000 C CNN
F 2 "Module:Arduino_Nano" H 1200 6750 50  0001 C CIN
F 3 "http://www.mouser.com/pdfdocs/Gravitech_Arduino_Nano3_0.pdf" H 1200 6750 50  0001 C CNN
	1    1200 6750
	1    0    0    -1  
$EndComp
$Comp
L Driver_Motor:Pololu_Breakout_A4988 A2
U 1 1 5F42F060
P 2500 6500
F 0 "A2" H 2550 7381 50  0000 C CNN
F 1 "Pololu_Breakout_A4988" H 2550 7290 50  0000 C CNN
F 2 "Module:Pololu_Breakout-16_15.2x20.3mm" H 2775 5750 50  0001 L CNN
F 3 "https://www.pololu.com/product/2980/pictures" H 2600 6200 50  0001 C CNN
	1    2500 6500
	1    0    0    -1  
$EndComp
$Comp
L Switch:SW_SPST SW2
U 1 1 5F4300A8
P 5200 6200
F 0 "SW2" H 5200 6435 50  0000 C CNN
F 1 "SW_SPST" H 5200 6344 50  0000 C CNN
F 2 "" H 5200 6200 50  0001 C CNN
F 3 "~" H 5200 6200 50  0001 C CNN
	1    5200 6200
	1    0    0    -1  
$EndComp
$Comp
L Device:LED D1
U 1 1 5F430A93
P 4350 7550
F 0 "D1" H 4343 7766 50  0000 C CNN
F 1 "LED" H 4343 7675 50  0000 C CNN
F 2 "" H 4350 7550 50  0001 C CNN
F 3 "~" H 4350 7550 50  0001 C CNN
	1    4350 7550
	1    0    0    -1  
$EndComp
$Comp
L Device:R R3
U 1 1 5F433662
P 4500 6600
F 0 "R3" H 4570 6646 50  0000 L CNN
F 1 "220R" H 4570 6555 50  0000 L CNN
F 2 "" V 4430 6600 50  0001 C CNN
F 3 "~" H 4500 6600 50  0001 C CNN
	1    4500 6600
	1    0    0    -1  
$EndComp
$Comp
L Connector:Conn_01x02_Female J4
U 1 1 5F434875
P 4850 7550
F 0 "J4" H 4878 7526 50  0000 L CNN
F 1 "Conn_01x02_Female" H 4878 7435 50  0000 L CNN
F 2 "" H 4850 7550 50  0001 C CNN
F 3 "~" H 4850 7550 50  0001 C CNN
	1    4850 7550
	1    0    0    -1  
$EndComp
$Comp
L Switch:SW_SPST SW3
U 1 1 5F435372
P 6650 7450
F 0 "SW3" H 6650 7685 50  0000 C CNN
F 1 "SW_SPST" H 6650 7594 50  0000 C CNN
F 2 "" H 6650 7450 50  0001 C CNN
F 3 "~" H 6650 7450 50  0001 C CNN
	1    6650 7450
	1    0    0    -1  
$EndComp
$Comp
L Connector:Conn_01x02_Female J7
U 1 1 5F435E50
P 5950 7400
F 0 "J7" H 5978 7376 50  0000 L CNN
F 1 "Conn_01x02_Female" H 5978 7285 50  0000 L CNN
F 2 "" H 5950 7400 50  0001 C CNN
F 3 "~" H 5950 7400 50  0001 C CNN
	1    5950 7400
	1    0    0    -1  
$EndComp
$Comp
L Device:CP C1
U 1 1 5F4368C5
P 7650 5700
F 0 "C1" H 7768 5746 50  0000 L CNN
F 1 "100uf" H 7768 5655 50  0000 L CNN
F 2 "" H 7688 5550 50  0001 C CNN
F 3 "~" H 7650 5700 50  0001 C CNN
	1    7650 5700
	1    0    0    -1  
$EndComp
$Comp
L Device:CP C2
U 1 1 5F43713F
P 8150 6100
F 0 "C2" H 8268 6146 50  0000 L CNN
F 1 "100uf" H 8268 6055 50  0000 L CNN
F 2 "" H 8188 5950 50  0001 C CNN
F 3 "~" H 8150 6100 50  0001 C CNN
	1    8150 6100
	1    0    0    -1  
$EndComp
$Comp
L Connector:Barrel_Jack J13
U 1 1 5F438571
P 9950 6300
F 0 "J13" H 10007 6625 50  0000 C CNN
F 1 "Barrel_Jack" H 10007 6534 50  0000 C CNN
F 2 "" H 10000 6260 50  0001 C CNN
F 3 "~" H 10000 6260 50  0001 C CNN
	1    9950 6300
	1    0    0    -1  
$EndComp
$Comp
L Device:Transformer_1P_1S T1
U 1 1 5F43A223
P 10300 4700
F 0 "T1" H 10300 5081 50  0000 C CNN
F 1 "Transformer_1P_1S" H 10300 4990 50  0000 C CNN
F 2 "" H 10300 4700 50  0001 C CNN
F 3 "~" H 10300 4700 50  0001 C CNN
	1    10300 4700
	1    0    0    -1  
$EndComp
$Comp
L Motor:Stepper_Motor_bipolar M1
U 1 1 5F43B30E
P 3000 1200
F 0 "M1" H 3188 1324 50  0000 L CNN
F 1 "Stepper_Motor_bipolar" H 3188 1233 50  0000 L CNN
F 2 "" H 3010 1190 50  0001 C CNN
F 3 "http://www.infineon.com/dgdl/Application-Note-TLE8110EE_driving_UniPolarStepperMotor_V1.1.pdf?fileId=db3a30431be39b97011be5d0aa0a00b0" H 3010 1190 50  0001 C CNN
	1    3000 1200
	1    0    0    -1  
$EndComp
$Comp
L Connector:Screw_Terminal_01x04 J2
U 1 1 5F4422BB
P 3350 2050
F 0 "J2" H 3430 2042 50  0000 L CNN
F 1 "Screw_Terminal_01x04" H 3430 1951 50  0000 L CNN
F 2 "" H 3350 2050 50  0001 C CNN
F 3 "~" H 3350 2050 50  0001 C CNN
	1    3350 2050
	1    0    0    -1  
$EndComp
$EndSCHEMATC
