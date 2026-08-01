# automatomic

## Overview

Automatomic is a cloud-native software delivery platform that automates the process of building, testing, securing, deploying, and monitoring applications.

The goal is to help developers ship reliable software faster by reducing deployment complexity and improving visibility into application health.

## Why Automatomic?

Modern applications require complex deployment workflows involving containers, cloud infrastructure, and multiple services. Automatomic provides a centralized platform to automate these processes and reduce deployment failures.

## What Makes It Different?

Unlike traditional CI/CD tools focused only on running pipelines, Automatomic focuses on **software reliability** by combining:

- Automated CI/CD workflows
- Kubernetes-based deployments
- Security validation
- Real-time monitoring
- Deployment health tracking
- Rollback support

## Tech Stack

**Frontend**
- Next.js + TypeScript  
  - Dashboard and real-time deployment visualization

**Backend**
- Go  
  - APIs, pipeline orchestration, job management

**Infrastructure**
- Docker  
  - Containerized services and build environments
- Kubernetes  
  - Deployment orchestration and scalable runners
- AWS  
  - Cloud infrastructure

**Data**
- PostgreSQL  
  - Application data and deployment history
- Redis  
  - Job queues and real-time events

**CI/CD**
- GitHub Actions  
  - Automated testing and deployment

## C++ Component

A lightweight C++ runner handles low-level workload execution, including:

- Process management
- Resource monitoring
- Log streaming

C++ is used for performance-critical operations requiring efficient system-level control.

## Arduino Integration

An Arduino/ESP32 device acts as a physical deployment monitor, displaying live build and production status through WiFi communication.

Example:
AUTOMATOMIC

BUILD: SUCCESS
DEPLOY: ONLINE

