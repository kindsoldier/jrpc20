### IoT-man, working ptototype of IoT device manager

It's is *weekend project* for design compact and fast working propotype of IoT device manager.

I used MQTT device for sample.

With my old laptop Lenovo X240 reached limit 
of processing **~5500 MQTT** messages per second.

For clarity of the prototype, I did not make device classes 
and limited the code to one instance of an device with one MQTT topic

Also for simplicity and clarity I have excluded: 
- authorization and authentication
- statistics and accounting
- saving states to database
- web socket control
and so on, so on

The code is 20 pages in total. 

Nevertheless, this is already a small framework for creating a manager 
of iot-like devices. The code can be adapted for any bus or network transport. 



