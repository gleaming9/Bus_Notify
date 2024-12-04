<h1 align="center">🚍 Bus Notify Service</h1>
<p align="center">
    <img src="https://img.shields.io/badge/Language-Go-blue" alt="Language" />
    <img src="https://img.shields.io/badge/Framework-Gin-lightblue" alt="Framework" />
    <img src="https://img.shields.io/badge/Queue-RabbitMQ-orange" alt="Queue" />
    <img src="https://img.shields.io/badge/Deployment-Docker-yellow" alt="Deployment" />
    <img src="https://img.shields.io/badge/Status-Complete-green" alt="Status" />
</p>

<p align="center"> A real-time bus notification service using RabbitMQ and Gin Framework.<br> Designed to provide accurate and timely bus arrival information for Gyeonggi-do bus stops. </p>
<!--📢🕒-->
<hr/>

<h2>📋 Features</h2>

<ul>
    <li><b>Real-time Bus Info</b>: Retrieve the arrival information of buses at a specified bus stop.</li>
    <li><b>Custom Notifications</b>:
        <ul>
            <li>Set bus stop, target time, and email for personalized notifications.</li>
            <li>Receive notifications when the nearest bus is arriving within 15 minutes of the target time.</li>
        </ul>
    </li>
    <li><b>Server Health Check</b>: Use `/health` endpoint to monitor server status.</li>
</ul>

<hr/>

<h2>📂 Project Structure</h2>

<pre>
.
├── api              # External API integrations and data handling
├── consume          # RabbitMQ consumers and message processing
├── handlers         # Request and response handling
├── model            # Data models and structures
├── produce          # RabbitMQ publishers
├── routes           # API routing
├── service          # Business logic and core functionality
├── Dockerfile       # Dockerfile for container setup
├── docker-compose.yml # Docker Compose configuration
└── README.md        # Project documentation
</pre>

<hr/>

<h2>🚀 How to Run</h2>

<ol>
    <li>Clone the repository:</li>
    <pre><code>git clone https://github.com/gleaming9/Bus_Notify.git</code></pre>
    <li>Build and run the project using Docker Compose:</li>
    <pre><code>docker-compose up --build</code></pre>
    <li>Access the API endpoints:</li>
    <ul>
        <li><b>Health Check</b>: <code>curl http://localhost:9090/health</code></li>
        <li><b>Bus Info</b>: <code>curl http://localhost:9090/bus-info?stationName=StationName</code></li>
        <li><b>Alert</b>:
            <pre>
curl -X POST http://localhost:9090/alert \
-H "Content-Type: application/json" \
-d '{
    "stationName": "StationName",
    "email": "example@gmail.com",
    "targetTime": "15:30"
}'
            </pre>
        </li>
    </ul>
</ol>

<hr>

<h2>📖 About</h2>
<p>The Bus Notify Service was developed to provide a real-time bus arrival notification system. Utilizing RabbitMQ, the service processes notification requests and efficiently delivers email alerts for buses arriving near the target time. The project is containerized for easy deployment and management.</p>

