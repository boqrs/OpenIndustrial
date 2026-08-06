### OpenIndustrial

OpenIndustrial is an open-source industrial IoT platform engineered for scalability, flexibility, and resilience. It provides a comprehensive solution for collecting data from industrial devices, processing it in the cloud, and enabling powerful business applications.

The platform's architecture is cleanly divided into three core, decoupled components that work in concert:

1. Gateway: The Edge Runtime
   The gateway is a lightweight, high-performance application designed to run on edge hardware, close to the physical devices (PLCs, sensors, CNCs, etc.). Its sole responsibilities are:

Data Acquisition: Connects to a wide range of industrial devices using various protocols (Modbus, OPC-UA, EtherNet/IP, etc.) via a pluggable driver system.
Protocol Translation: Normalizes disparate industrial protocols into a unified data model.
Secure Communication: Publishes data reliably to the cloud message bus and receives commands from it.
The Gateway is designed for high reliability and low resource consumption, ensuring robust operation in demanding industrial environments.

2. Cloud: The Business Platform
   The cloud is the central nervous system of the platform. It is a multi-tenant, service-oriented application responsible for all business logic and data processing. Key modules include:

Identity & Access Management: Manages organizations, users, roles, and permissions (org, user, role, auth).
Digital Twin Core: Models the physical world with entities like product (blueprints), asset (unique instances/SNs), gateway (edge points), and device (equipment).
MES & Operations: Orchestrates factory operations with modules like workorder.
Data Processing & Analytics: Ingests data from the message bus, stores it, and prepares it for analytics and visualization.
The Cloud is built on a modern, domain-driven design (DDD) architecture, making it highly extensible and easy to maintain.

3. Message Bus: The Communication Backbone
   Underpinning the entire platform is a message broker (like Mosquitto MQTT), which acts as the communication backbone. This component decouples the gateway from the cloud:

Decoupling: The Gateway and Cloud do not communicate directly. They only publish and subscribe to topics on the message bus. This allows them to be developed, deployed, and scaled independently.
Reliability: Ensures that data from the edge is not lost, even with intermittent network connectivity.
Bidirectional: Enables both data telemetry from the edge to the cloud and command-and-control messages from the cloud to the edge.
This robust, three-part architecture makes OpenIndustrial a powerful foundation for building sophisticated industrial applications, from real-time monitoring and asset tracking to predictive maintenance and full-scale Manufacturing Execution Systems (MES).
