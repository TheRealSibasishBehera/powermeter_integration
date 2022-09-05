<h1>PowerMeter Monitoring</h1>
<h3>Using Prometheus and Graphana</h3>

Using OS- tools can more or less give accurate information about the energy consumption by the operating system and the CPU. In any general device, this could be sufficient as in such devices as the most energy supply Is used by CPU computation.

But in cases of SBCs like Raspberry Pi, Nvidia Jetson Boards this is not sufficient as there are other factors as well which consume energy comparable to that of CPU computation. Eg.Energy consumption by networking components like Ethernet, GPIO, UART, etc.

To overcome this problem and to get the complete power consumption, an external power meter needs to be integrated into the SBC
![f83f1a_e2d765a97e914125b200aeefdd810740_mv2](https://user-images.githubusercontent.com/95071627/188454874-a0691223-b8ec-4ee5-b1ac-c28d69ba6306.jpg)


For this we  would be using Tasmota EU plug V2 by Athom . This is based on tasmota-HLW8032 , providing control using MQTT,Web UI , HTTP.


<h3>Local SetUp</h3>

<h3>Dev SetUp</h3>
<h4>Pre-requisite</h4>
<ol>
   <li>Go compiler <ul>
</ol>
<ol>
   <li> Configured PowerMeter (Tasmotta-HLW8032)  <ul>
</ol>



open up a terminal

1. Clone the repo ,go in the folder

2. Run using go compiler <code>sudo go run cmd/main.go</code>  
   powertop requires sudo permission to access the system stats

3. Bare prometheus metrics can be seen using <code>curl 0.0.0.0:8881/metrics</code>

<h3>Run Using Docker</h3>

1. <code> docker run -d -p 8888:8881 sibseh/powermeter:v4</code>  
2. Bare prometheus metrics can be seen using <code> curl 0.0.0.0:8888/metrics |grep current_count\|voltage_count</code>

These can be run with graphana and prometheus easily with the docker compose file

<h3>Monitoring with Graphana and Prometheus using docker compose </h3>

1. Open up a terminal in the same directory <code>docker-compose up</code>

2. Open your favourite brower with localhost:3000 , it will open up Graphana, login with username and password both as <code>admin</code>

3. Go to configuration ->Data sources -> Add Prometheus -> set Http as <code>http://prometheus:9090</code>

4. Go to create -> Dashboard -> Select one

5. Add current_count , voltage_count

6 . Now you can see clearly the parameters of your system calculated !!


The final set up should look like this
![Screenshot from 2022-09-05 16-50-16](https://user-images.githubusercontent.com/95071627/188454799-a4099d31-dd7b-45fc-940b-89a4ba6c0aa8.png)



Viewing voltage_count variation
