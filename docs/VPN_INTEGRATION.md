# VPN Integration Documentation

## Overview
This documentation outlines provider-agnostic VPN integration supporting WireGuard and OpenVPN.

## Supported Config Import Formats
- **WireGuard**: Configurations in `.conf` format.
- **OpenVPN**: Configurations in `.ovpn` format, alongside certificates/keys.

## Local-First Data Flow
The VPN integration ensures that all data flows locally before reaching the VPN.

## Kill Switch
A built-in kill switch is implemented to prevent data leakage in case of VPN disconnection.

## VPN-Required Mode for Remote Admin
Remote admin requires VPN connection to ensure secure access to resources.

## DNS-over-VPN
DNS requests are routed through the VPN to enhance privacy and security.

## Split Tunneling/Policy Routing
Policy routing enables selection of which traffic goes through the VPN and which does not.

## Privacy Guidance
Users are advised to minimize sharing of MAC/IP lists with third parties as VPN is primarily a transport mechanism, not an identity shield.

## Interoperability Requirements
Support for open-source clients and standard configurations is critical to ensure interoperability among different systems and services.