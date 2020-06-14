/*
DSLR Intervalometer
 */

//setup variables
int sensorPin = A0;
int sensorValue = 0;
int lapsePeriod = 0;
int x = 0;

void setup() {
  pinMode(2, OUTPUT);
  pinMode(3, OUTPUT);
  pinMode(7, OUTPUT);
  pinMode(8, OUTPUT);
  
  
    //determines how many seconds it should wait between pictures.
    sensorValue = analogRead(sensorPin); 
    if (sensorValue >= 20) {
      lapsePeriod = 10 + (sensorValue / 20);
    }
    else {
       lapsePeriod = 10;
    }


  //blink 3 times to indicate it has started
  digitalWrite(7, HIGH);
  delay(250);
  digitalWrite(7, LOW);
  delay(250);
  digitalWrite(7, HIGH);
  delay(250);
  digitalWrite(7, LOW);
  delay(250);
  digitalWrite(7, HIGH);
  delay(250);
  digitalWrite(7, LOW);
  delay(250);
}


void loop() {
  
  //waits however many seconds you have chosen
  while (x < lapsePeriod) {
    delay(1000);
    x++;
  }
  x = 0;
  
  //focus and blink focus LED for one second
  digitalWrite(3, HIGH);
  digitalWrite(8, HIGH);
  delay(1000);
  digitalWrite(3, LOW);
  digitalWrite(8, LOW);
  delay(100);
  
  //shoot and blink shoot LED for one second 
  digitalWrite(2, HIGH);
  digitalWrite(7, HIGH);
  delay(1000);
  digitalWrite(2, LOW);
  digitalWrite(7, LOW);

}
