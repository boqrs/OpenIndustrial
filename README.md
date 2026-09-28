# OpenIndustrial Industrial Cloud

> An Industrial Cloud platform connecting customers, orders, manufacturing, physical devices, and operational data through a unified industrial digital identity.

---

## 1. Project Overview

**OpenIndustrial** is an Industrial Cloud platform designed to connect the complete lifecycle of industrial products and physical devices.

It is not intended to be only a traditional MES system.

OpenIndustrial connects:

```text
Customer
   │
   ▼
Sales Order
   │
   ▼
Manufacturing
   │
   ▼
Physical Device
   │
   ▼
Warehouse
   │
   ▼
Shipment
   │
   ▼
Customer
   │
   ▼
Device Activation
   │
   ▼
IoT
   │
   ▼
Operational Data
```

The key concept connecting these stages is the **Device Digital Identity**.

A physical device is created during manufacturing, receives its digital identity, enters the warehouse and sales lifecycle, and eventually becomes an IoT-connected device after reaching the customer.

This creates a continuous digital lifecycle from:

**customer → order → factory → production → product → device → customer → runtime data**

---

# 2. Architecture

OpenIndustrial is organized around several business domains sharing a common resource and identity foundation.

```text
                         OpenIndustrial
                      Industrial Cloud
                             │
             ┌───────────────┼───────────────┐
             │               │               │
            CRM             MES             WMS
             │               │               │
       Customer / Order   Manufacturing   Inventory / Shipment
             │               │               │
             └───────────────┼───────────────┘
                             │
                             ▼
                 ┌─────────────────────┐
                 │ Device Digital      │
                 │ Identity            │
                 │                     │
                 │ Device              │
                 │ Serial Number       │
                 │ Resource Identity   │
                 │ Certificate         │
                 └──────────┬──────────┘
                            │
                 ┌──────────┴──────────┐
                 │                     │
                 ▼                     ▼
             Security                 IoT
                 │                     │
       Identity / Credential     MQTT / Runtime
                 │                     │
                 └──────────┬──────────┘
                            │
                            ▼
                    Device Runtime
```

The architecture can be understood as three layers:

### Business Layer

Handles the industrial business lifecycle:

- CRM
- Sales
- MES
- WMS

### Identity & Security Layer

Provides the common identity foundation:

- Resource
- Device identity
- Credentials
- Certificates
- Authorization

### Runtime Layer

Connects physical devices with the cloud:

- IoT
- MQTT
- Device status
- Commands
- Runtime events
- Operational data

---

# 3. Core Architecture Principle

The core architectural principle of OpenIndustrial is:

> **The physical device is created by manufacturing and connected to IoT later.**

This distinction is important.

```text
                    Manufacturing
                         │
                         ▼
                  Production Result
                         │
                         ▼
                  ┌──────────────┐
                  │    Device    │
                  │              │
                  │ Device ID    │
                  │ Serial No.   │
                  │ Resource ID  │
                  │ Certificate  │
                  └──────┬───────┘
                         │
                         │
              ┌──────────┴──────────┐
              │                     │
              ▼                     ▼
             WMS                   IoT
          Warehouse             Activation
              │                     │
              ▼                     ▼
          Shipment              Runtime
```

Therefore:

```text
Manufacturing
     │
     └──► creates Device

IoT
     │
     └──► connects Device
```

IoT does **not** create a Device.

IoT operates on an already existing physical device and its established digital identity.

---

# 4. Business Modules

OpenIndustrial is divided into several business domains.

The modules have clear responsibilities rather than forming independent isolated systems.

```text
┌──────────────────────────────────────────────────────────┐
│                    OpenIndustrial                        │
├──────────────┬──────────────┬──────────────┬─────────────┤
│     CRM      │     MES      │     WMS      │     IoT     │
├──────────────┼──────────────┼──────────────┼─────────────┤
│ Customer     │ Product      │ Warehouse    │ Device      │
│ Sales Order  │ BOM          │ Inventory    │ Connection  │
│ Business     │ Routing      │ Shipment     │ MQTT        │
│              │ Production   │              │ Commands    │
│              │ Execution    │              │ Status      │
└──────────────┴──────────────┴──────────────┴─────────────┘
                         │
                         ▼
                  Device Identity
                         │
                         ▼
                      Security
```

### CRM

Responsible for the commercial relationship between the enterprise and its customers.

The primary business objects include:

- Customer
- Sales Order

CRM starts the commercial lifecycle.

---

### MES

Responsible for manufacturing.

The core production flow is:

```text
Product
   │
   ▼
BOM
   │
   ▼
Routing
   │
   ▼
Production Plan
   │
   ▼
Work Order
   │
   ▼
Execution
   │
   ▼
Execution Operation
   │
   ▼
Production Result
   │
   ▼
Device
```

MES is where the physical product becomes a traceable device entity.

---

### WMS

Responsible for the physical logistics lifecycle after manufacturing.

```text
Device
  │
  ▼
Stock In
  │
  ▼
Inventory
  │
  ▼
Reservation
  │
  ▼
Shipment
  │
  ▼
Customer
```

WMS connects the manufacturing world with the customer delivery lifecycle.

---

### Security

Security provides the common identity and credential foundation.

It is responsible for concepts such as:

- Resource identity
- Device credentials
- Certificates
- Certificate authority
- Authentication
- Authorization

Security is a shared capability rather than a separate business lifecycle.

---

### IoT

IoT connects an already existing Device to its runtime environment.

Its primary responsibilities include:

- MQTT authentication
- MQTT authorization
- Device connection
- Device online/offline state
- Heartbeat
- Device status
- Cloud commands
- Command acknowledgements
- Runtime events

IoT does not own the manufacturing identity of a device.

---

# 5. Overall Business Flow

The complete OpenIndustrial business lifecycle can be represented as:

```text
                         ┌───────────┐
                         │ Customer  │
                         └─────┬─────┘
                               │
                               ▼
                         ┌───────────┐
                         │Sales Order│
                         └─────┬─────┘
                               │
                               ▼
                     ┌──────────────────┐
                     │ Production Plan  │
                     └────────┬─────────┘
                              │
                              ▼
                         ┌──────────┐
                         │Work Order│
                         └────┬─────┘
                              │
                              ▼
                     ┌─────────────────┐
                     │   Execution     │
                     │                 │
                     │ Operations      │
                     └────────┬────────┘
                              │
                              ▼
                    ┌───────────────────┐
                    │ Production Result │
                    └─────────┬─────────┘
                              │
                              │ Qualified Units
                              ▼
                    ╔═══════════════════╗
                    ║      DEVICE       ║
                    ║                   ║
                    ║ Device ID         ║
                    ║ Resource Identity ║
                    ║ Serial Number     ║
                    ║ Certificate       ║
                    ╚═════════╤═════════╝
                              │
                              ▼
                         ┌──────────┐
                         │   WMS    │
                         └────┬─────┘
                              │
                              ▼
                         ┌──────────┐
                         │ Shipment │
                         └────┬─────┘
                              │
                              ▼
                         ┌──────────┐
                         │ Customer │
                         └────┬─────┘
                              │
                              │ Activation
                              ▼
                       ┌──────────────┐
                       │     IoT      │
                       └──────┬───────┘
                              │
                              ▼
                      ┌───────────────┐
                      │ Device Runtime│
                      │               │
                      │ Status        │
                      │ Commands      │
                      │ Events        │
                      │ Telemetry     │
                      └───────────────┘
```

This flow is the central business model of OpenIndustrial.

---

# 6. Device Digital Identity

The most important concept in the OpenIndustrial architecture is the **Device Digital Identity**.

A product definition describes **what a product is**.

A Device represents **a specific physical instance of that product**.

```text
Product
   │
   │ Static Definition
   ▼
┌──────────────────┐
│ Product          │
│                  │
│ Name             │
│ Model            │
│ Attributes       │
│ Images           │
│ Specifications   │
└────────┬─────────┘
         │
         │ Manufacturing
         ▼
┌──────────────────────────┐
│ Device                   │
│                          │
│ Specific Physical Unit   │
│                          │
│ Device ID                │
│ Resource ID              │
│ Serial Number            │
│ Certificate              │
│ Runtime State            │
└──────────────────────────┘
```

The Device is therefore not merely another product table.

It represents a physical object with a persistent digital identity.

---

## 6.1 Device Identity Generation

Device identity is generated during manufacturing.

The key flow is:

```text
Product
   │
   ▼
BOM
   │
   ▼
Routing
   │
   ▼
Production Plan
   │
   ▼
Work Order
   │
   ▼
Execution
   │
   ▼
Execution Operation
   │
   ▼
Production Result
   │
   ▼
╔════════════════════════════╗
║      Device Identity       ║
║                            ║
║ Device ID                  ║
║ Resource Identity          ║
║ Serial Number              ║
║ Certificate / Credential   ║
╚════════════════════════════╝
```

A production result determines the actual qualified physical units produced by the manufacturing process.

Those qualified units become Devices.

This establishes the relationship:

```text
Manufacturing Result
        │
        ▼
Physical Device
        │
        ▼
Digital Identity
```

---

## 6.2 Identity Lifecycle

Once created, the Device identity continues through the rest of the product lifecycle.

```text
             Manufacturing
                   │
                   ▼
           Device Identity
                   │
        ┌──────────┼──────────┐
        ▼          ▼          ▼
       WMS      Shipment    Traceability
        │          │
        └────┬─────┘
             │
             ▼
          Customer
             │
             ▼
         Activation
             │
             ▼
            IoT
             │
             ▼
          Runtime
```

The same identity is therefore used across:

- Manufacturing
- Warehouse
- Shipment
- Customer ownership
- Activation
- IoT connectivity
- Runtime operations

This provides a continuous digital thread for the physical device.

---

# 7. Device Lifecycle

The Device lifecycle is different from the Product lifecycle.

```text
                    Product Definition
                           │
                           ▼
                    Manufacturing
                           │
                           ▼
                    Device Created
                           │
                           ▼
                     Warehouse
                           │
                           ▼
                       Shipped
                           │
                           ▼
                  Customer Possession
                           │
                           ▼
                       Activated
                           │
                           ▼
                    IoT Connected
                           │
                           ▼
                      In Operation
```

The lifecycle can therefore be divided into two major phases:

### Industrial Lifecycle

```text
Product
   ↓
Manufacturing
   ↓
Device Creation
   ↓
Warehouse
   ↓
Shipment
```

### Runtime Lifecycle

```text
Customer
   ↓
Activation
   ↓
IoT Connection
   ↓
Device Runtime
   ↓
Operational Data
```

The Device Digital Identity connects these two phases.

---

# 8. IoT Architecture

IoT is the runtime connection between the physical Device and the Industrial Cloud.

```text
┌─────────────────────┐
│   Physical Device   │
│                     │
│ Firmware / Runtime  │
└──────────┬──────────┘
           │
           │ MQTT
           ▼
┌────────────────────────────┐
│       MQTT Adapter         │
│                            │
│ Self-hosted MQTT           │
│ AWS IoT Core               │
└─────────────┬──────────────┘
              │
              ▼
┌────────────────────────────┐
│        IoT Service         │
│                            │
│ Authentication             │
│ Authorization              │
│ Device Status              │
│ Heartbeat                  │
│ Command                    │
│ Command ACK                │
│ Runtime Events             │
└─────────────┬──────────────┘
              │
       ┌──────┴───────┐
       ▼              ▼
   Device          Security
   Service         Service
```

The IoT layer does not need to know which MQTT implementation is being used.

The broker implementation is hidden behind the adapter layer.

```text
                  IoT Service
                       │
                       ▼
              DeviceMessageAdapter
                       │
              ┌────────┴────────┐
              │                 │
              ▼                 ▼
        MQTT Adapter      AWS Adapter
              │                 │
              ▼                 ▼
       MQTT Broker        AWS IoT Core
```

This allows the business layer to remain independent from the underlying messaging infrastructure.

---

# 9. MQTT Business Protocol

OpenIndustrial uses a simple device-specific MQTT topic model.

Each device has two primary business topics.

```text
devices/{resource_id}/status
devices/{resource_id}/command
```

For example:

```text
devices/10086/status
devices/10086/command
```

---

## 9.1 Device → Cloud

The Device publishes to:

```text
devices/{resource_id}/status
```

The status channel can carry different types of business messages.

```text
status
 ├── heartbeat
 ├── online
 ├── offline
 ├── command_ack
 └── event
```

Example:

```json
{
  "type": "heartbeat",
  "timestamp": "2026-09-24T01:20:00Z"
}
```

Event example:

```json
{
  "type": "event",
  "event": "production_completed",
  "timestamp": "2026-09-24T01:20:00Z",
  "data": {
    "xxx": "xxx"
  }
}
```

Command acknowledgement example:

```json
{
  "type": "command_ack",
  "command_id": 123,
  "success": true,
  "timestamp": "2026-09-24T01:20:00Z"
}
```

---

## 9.2 Cloud → Device

The Cloud publishes commands to:

```text
devices/{resource_id}/command
```

Example:

```json
{
  "command_id": 123,
  "command": "reboot",
  "payload": "{}"
}
```

The Device can then return the execution result through the status channel:

```text
Device
   │
   │ command_ack
   ▼
devices/{resource_id}/status
```

---

# 10. MQTT Authorization

Device access is scoped to its own Resource Identity.

For example, a device with:

```text
resource_id = 10086
```

can be authorized to:

```text
Publish:
devices/10086/status

Subscribe:
devices/10086/command
```

It must not access another device's topics:

```text
devices/10087/status
devices/10087/command
```

Conceptually:

```text
                    Device 10086
                         │
             ┌───────────┴───────────┐
             │                       │
          Publish                 Subscribe
             │                       │
             ▼                       ▼
    devices/10086/status    devices/10086/command
```

The Resource Identity therefore becomes the bridge between:

**business identity → security identity → IoT identity**

---

# 11. IoT and Device Identity

The relationship between Device and IoT is intentionally separated.

```text
              Manufacturing
                   │
                   ▼
          ┌─────────────────┐
          │     Device      │
          │                 │
          │ Resource ID     │
          │ Device ID       │
          │ Serial Number   │
          │ Certificate     │
          └────────┬────────┘
                   │
                   │ Existing Identity
                   ▼
          ┌─────────────────┐
          │      IoT        │
          │                 │
          │ Authentication  │
          │ Authorization   │
          │ MQTT Connection  │
          └────────┬────────┘
                   │
                   ▼
                Runtime
```

Therefore:

```text
Device = Who is this physical thing?

IoT = How does this physical thing connect and communicate?
```

This separation keeps the business identity lifecycle independent from the runtime communication layer.

---

# 12. Security and Digital Identity

Security provides the foundation for device identity and authentication.

Conceptually:

```text
Resource
   │
   ▼
Resource Identity
   │
   ├──────────────┐
   │              │
   ▼              ▼
Credential     Certificate
   │              │
   └──────┬───────┘
          ▼
      IoT Access
```

The security domain is the source of truth for device credentials and certificates.

IoT consumes these capabilities rather than implementing another independent device identity system.

---

# 13. Manufacturing to IoT

The most important end-to-end flow in OpenIndustrial is:

```text
                    MES
                     │
                     │ Production
                     ▼
             Production Result
                     │
                     │
                     ▼
              ┌──────────────┐
              │    Device    │
              └──────┬───────┘
                     │
                     │ Identity
                     ▼
              ┌──────────────┐
              │     WMS      │
              └──────┬───────┘
                     │
                     │ Shipment
                     ▼
              ┌──────────────┐
              │   Customer   │
              └──────┬───────┘
                     │
                     │ Activation
                     ▼
              ┌──────────────┐
              │     IoT      │
              └──────┬───────┘
                     │
                     │ MQTT
                     ▼
              ┌──────────────┐
              │   Runtime    │
              └──────────────┘
```

This is the central lifecycle that OpenIndustrial is designed to support.

The important point is that the device identity remains the same throughout the lifecycle.

---

# 14. Runtime Data

IoT runtime communication is divided conceptually into business messages and high-volume telemetry.

```text
                         Device
                           │
                           │
                    ┌──────┴──────┐
                    │             │
                    ▼             ▼
              Business Data   Telemetry
                    │             │
                    ▼             ▼
                  MQTT        Data Pipeline
                    │             │
                    ▼             ▼
               IoT Service       TSDB
```

Business-level messages include:

- Online / Offline
- Heartbeat
- Commands
- Command acknowledgements
- Important device events

High-volume telemetry can be handled through a separate data path rather than forcing every telemetry message through the normal business service layer.

This keeps the business architecture independent from the scale of device telemetry.

---

# 15. Multi-Tenant Architecture

OpenIndustrial is designed as a multi-tenant Industrial Cloud.

Conceptually:

```text
                 Industrial Cloud
                        │
        ┌───────────────┼───────────────┐
        │               │               │
      Tenant A        Tenant B        Tenant C
        │               │               │
      ┌─┴─┐           ┌─┴─┐           ┌─┴─┐
      │CRM│           │CRM│           │CRM│
      │MES│           │MES│           │MES│
      │WMS│           │WMS│           │WMS│
      │IoT│           │IoT│           │IoT│
      └───┘           └───┘           └───┘
```

Business resources are associated with their corresponding tenant context.

Device identity, production, warehouse, customer and IoT access therefore operate within the tenant boundary.

---

# 16. Backend Architecture

The backend follows a layered service architecture.

```text
┌─────────────────────────────────────────┐
│              API / Handler              │
├─────────────────────────────────────────┤
│                Service                  │
├─────────────────────────────────────────┤
│              Repository                 │
├─────────────────────────────────────────┤
│             Persistence                │
├─────────────────────────────────────────┤
│ PostgreSQL / Redis / External Services │
└─────────────────────────────────────────┘
```

The business service layer coordinates domain operations.

Repositories isolate persistence concerns.

Shared infrastructure such as persistence, transactions, authentication and external integrations remains separated from the business domain.

---

# 17. Resource Kernel

At the foundation of the domain model is the **Resource** concept.

```text
                    Resource
                       │
       ┌───────────────┼───────────────┐
       │               │               │
     Product         Device          Other
       │               │
       ▼               ▼
 Static Data      Runtime Identity
```

Resource provides a common foundation for physical and logical industrial objects.

The Product and Device domains then build their own business meaning on top of this foundation.

The purpose of the Resource Kernel is to provide a consistent identity and relationship model across different industrial domains.

---

# 18. Repository Structure

The backend is organized by business service and shared infrastructure.

Conceptually:

```text
cloud/
├── bootstrap/
│
├── internal/
│   ├── handler/
│   │
│   ├── services/
│   │   ├── customer/
│   │   ├── salesorder/
│   │   ├── product/
│   │   ├── bom/
│   │   ├── routing/
│   │   ├── planning/
│   │   ├── workorder/
│   │   ├── execution/
│   │   ├── executionresult/
│   │   ├── manufacturing/
│   │   ├── device/
│   │   ├── wms/
│   │   ├── security/
│   │   ├── identity/
│   │   ├── iot/
│   │   └── event/
│   │
│   ├── persistence/
│   │
│   └── ...
│
└── ...
```

The exact implementation structure may evolve, while the domain boundaries remain the primary organizational principle.

---

# 19. Technology Architecture

OpenIndustrial is built around a modern cloud backend architecture.

The current backend stack includes technologies such as:

```text
Go
 │
 ├── HTTP / API
 ├── Service Layer
 ├── Repository Layer
 └── Domain Services
        │
        ├── PostgreSQL
        ├── Redis
        ├── MQTT
        └── Cloud Infrastructure
```

The IoT communication layer is designed to support multiple MQTT infrastructure implementations through adapters.

```text
                    IoT Business Layer
                            │
                            ▼
                  DeviceMessageAdapter
                     │            │
                     ▼            ▼
                   MQTT        AWS IoT
```

---

# 20. Version 1.0

The current backend architecture is treated as the **v1.0 baseline**.

The major business chain is established:

```text
CRM
 │
 ▼
Sales
 │
 ▼
MES
 │
 ▼
Device Identity
 │
 ▼
WMS
 │
 ▼
Shipment
 │
 ▼
Customer
 │
 ▼
Device Activation
 │
 ▼
IoT
 │
 ▼
Runtime
```

The v1.0 baseline focuses on business closure rather than continued architectural redesign.

The core principles are:

- Device is created through manufacturing.
- IoT does not create Device.
- Device identity is persistent across the lifecycle.
- Security owns credentials and certificates.
- IoT owns runtime connectivity.
- MQTT infrastructure is accessed through adapters.
- CRM, MES, WMS and IoT remain separate business domains.
- Existing architecture is preferred over unnecessary redesign.

---

# 21. Current Development Phase

The backend is currently treated as the v1.0 baseline.

The next major development phase is the **Frontend Management Console**.

The frontend will provide the operational interface for the Industrial Cloud.

Conceptually:

```text
                    Management Console
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
       CRM                MES                WMS
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                           ▼
                    Device Management
                           │
                           ▼
                          IoT
```

The frontend is intended to expose the business lifecycle already established by the backend rather than redefine the underlying domain architecture.

---

# 22. Future Evolution

The platform can evolve around the established business foundation.

```text
                         OpenIndustrial
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
     Business              Identity              Runtime
        │                     │                     │
   CRM / MES / WMS       Device Identity          IoT
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
                              ▼
                       Industrial Data
                              │
                              ▼
                   Analytics / Intelligence
```

Potential future capabilities include:

- Advanced device telemetry
- Time-series data platform
- Device monitoring
- Remote diagnostics
- OTA updates
- Industrial analytics
- Production analytics
- Device lifecycle analytics
- Digital-twin capabilities
- Additional cloud infrastructure adapters

These capabilities can be added without changing the fundamental lifecycle:

```text
Manufacture
    ↓
Device Identity
    ↓
Deliver
    ↓
Activate
    ↓
Connect
    ↓
Operate
    ↓
Analyze
```

---

# 23. Development Principles

OpenIndustrial follows several core principles.

### 1. Business Lifecycle First

The platform is designed around the real industrial lifecycle rather than isolated CRUD modules.

```text
Customer
   ↓
Order
   ↓
Manufacturing
   ↓
Device
   ↓
Warehouse
   ↓
Shipment
   ↓
Customer
   ↓
IoT
   ↓
Runtime
```

---

### 2. Device Identity Is the Digital Thread

The Device identity connects the manufacturing world with the runtime world.

```text
Manufacturing
     │
     ▼
 Device Identity
     │
     ├── WMS
     ├── Shipment
     ├── Customer
     ├── Activation
     └── IoT
```

---

### 3. Manufacturing Creates Devices

A Device represents a real physical unit.

Therefore:

```text
Production Result
       │
       ▼
     Device
```

rather than:

```text
IoT
 │
 └──► Device
```

---

### 4. IoT Connects Existing Devices

IoT is responsible for runtime connectivity and communication.

```text
Existing Device
      │
      ▼
Authentication
      │
      ▼
Authorization
      │
      ▼
MQTT Connection
      │
      ▼
Runtime
```

---

### 5. Security Is a Shared Capability

Credentials, certificates and identity-related security capabilities are centralized rather than duplicated in individual business modules.

---

### 6. Infrastructure Is Replaceable

MQTT infrastructure is abstracted through an adapter layer.

```text
                    IoT
                     │
                     ▼
              Adapter Interface
                 │       │
                 ▼       ▼
               MQTT     AWS
```

The business layer should not depend directly on a specific broker implementation.

---

### 7. Avoid Unnecessary Abstraction

The project prefers meaningful business concepts over excessive technical abstraction.

The architecture should remain understandable from the business lifecycle itself.

---

# 24. The OpenIndustrial Digital Thread

The ultimate goal of OpenIndustrial can be summarized by the following digital thread:

```text
┌──────────┐
│ Customer │
└────┬─────┘
     │
     ▼
┌──────────┐
│   Order  │
└────┬─────┘
     │
     ▼
┌──────────────┐
│ Manufacturing│
└──────┬───────┘
       │
       ▼
╔══════════════════════╗
║  Device Digital ID   ║
║                      ║
║ Device               ║
║ Serial Number        ║
║ Resource Identity    ║
║ Certificate          ║
╚══════════╤═══════════╝
           │
           ▼
      ┌──────────┐
      │ Warehouse│
      └────┬─────┘
           │
           ▼
      ┌──────────┐
      │ Shipment │
      └────┬─────┘
           │
           ▼
      ┌──────────┐
      │ Customer │
      └────┬─────┘
           │
           ▼
      ┌──────────┐
      │Activation│
      └────┬─────┘
           │
           ▼
      ┌──────────┐
      │   IoT    │
      └────┬─────┘
           │
           ▼
      ┌──────────┐
      │ Runtime  │
      └────┬─────┘
           │
           ▼
    Operational Data
```

This is the central idea of OpenIndustrial:

> **A product starts as a business definition, becomes a physical device through manufacturing, receives a persistent digital identity, reaches the customer through the supply chain, and finally enters the IoT runtime world.**

The platform connects these stages into one continuous industrial lifecycle.

---

# 25. Roadmap

### v1.0

```text
CRM
 │
MES
 │
Device Identity
 │
WMS
 │
Shipment
 │
Activation
 │
IoT
```

The backend business foundation is established.

---

### Next Phase

```text
Frontend Management Console
            │
            ▼
     Business Operations
            │
            ▼
      Device Management
            │
            ▼
       IoT Operations
```

---

### Future

```text
IoT
 │
 ├── Telemetry
 ├── Time Series
 ├── Monitoring
 ├── Diagnostics
 ├── OTA
 └── Analytics
        │
        ▼
Industrial Intelligence
```

---

# 26. Summary

OpenIndustrial is an attempt to build an Industrial Cloud around the **complete lifecycle of physical products and devices**.

Its core is not a single module.

Its core is the connection between business, manufacturing, physical identity and runtime.

```text
        CRM
         │
         ▼
       Order
```
