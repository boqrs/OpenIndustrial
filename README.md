OpenIndustrial is an industrial platform that can be clearly divided into three core parts, which together form a cohesive whole:

1. The Edge: OpenIndustrial Gateway

Role: Data collector, translator, and frontline sentry.

What we've completed:

Core Framework: Gateway, Driver, Registry, Collector.

Data Unification: Defined a standard Event format.

Driver Integration: Successfully integrated BACnet (real protocol) and Simulator (testing tool).

Status: The foundation is solid and ready to stream data to the cloud.

2. The Cloud: OpenIndustrial Cloud Platform

Role: Data aggregation center, processing brain, and storage warehouse.

Core components (to be built):

Data Ingestion Service: A high-availability MQTT Broker as the entry point for all gateway data reports.

Data Processing Service: Subscribes to MQTT data for real-time computation, rule engine evaluation, and alert generation.

Data Storage Service: Stores processed data in a Time-Series Database (e.g., InfluxDB or TimescaleDB) for historical traceability and analysis.

Device Management Service: Maintains the status (online/offline), configuration versions, and health of all gateways, enabling Digital Twin capabilities.

API Service: Provides RESTful or GraphQL APIs for the management console to query data, control devices, and push configurations.

3. The Application: OpenIndustrial Management Console

Role: The unified management interface and data visualization window for the platform.

Core features (to be built):

Data Visualization: Real-time dashboards and charts displaying device data queried from the cloud.

Gateway Management: View the list and status of all registered gateways and perform operations on them.

Remote Configuration (Core Feature): Modify gateway configurations via the web interface (e.g., adding/removing a driver, modifying a data point). Upon saving, the cloud platform pushes the new configuration to the corresponding edge gateway through a command channel, and the gateway automatically applies the updates.

Alert Center: View and manage system-generated alerts.
