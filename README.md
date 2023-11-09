# bvr logs
*bvr collects logs at the dam to centralize availability.*
---
## Architecture
**bvr logs** is composed of two applications written in **Go**.
### bvr client
**bvr client** consumes logs to send to **bvr dam**. Logs are sent via **HTTP POST** to the dam location.
### bvr dam
**bvr dam** has an endpoint for **bvr client** to POST logs to. **bvr dam** consolidates logs with a configurable TTL
## Status
No functionality has been implemented yet for **bvr client** and **bvr dam**.
