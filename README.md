# 🎯 Ping-Pong Microservices with TLS 🔐

Welcome to the **Ping-Pong Microservices** project! This microservice architecture showcases how two simple applications, **Pinger** and **Ponger**, communicate securely via TLS, containerized with Docker, and orchestrated using Kubernetes. The aim is to demonstrate secure, scalable communication between services using lightweight Go applications and a k3d cluster.

## 🌟 Key Features
- **Microservices Architecture**: Two services:
  - **Pinger**: Sends "ping" requests.
  - **Ponger**: Receives the "ping" and responds with "pong."
  
- **TLS Encryption** 🔐: Ensures secure communication between Pinger and Ponger using self-signed certificates.

- **Kubernetes-Orchestrated**: Both services run on a lightweight k3d cluster, making it perfect for development environments.

- **Containerized with Docker 🐳**: Ensures easy deployment and scalability of both services.

- **Port-forwarding** 🚀: For quick local access and testing of services.

---

### 📜 Requirements:

1. **Deploy a K3D cluster** with at least 1 server and 2 agents.
2. **Implement secure communication** between two microservices (`pinger` and `ponger`) using TLS certificates.
3. **Create a self-signed certificate** using `openssl` to secure the connection.
4. Configure Kubernetes manifests, including **services and deployments**, to enable secure communication.
5. Perform **port forwarding** to access services from an external environment (EC2 in this case).
6. Ensure the deployment is fully functional, check logs for correct communication, and validate the TLS connection using a browser like Mozilla Firefox.

---

## 🛠️ Project Setup & Tools
This project leverages several key technologies:
- **Go**: The language of choice for both services.
- **Docker**: To containerize the apps.
- **k3d**: For managing the lightweight Kubernetes cluster.
- **TLS Certificates**: Used to secure communication between services.

---

## 📂 File Structure Overview

```bash
├── Makefile                # Automation commands
├── README.md               # This README file
├── app
│   ├── certs               # TLS certificates
│   │   ├── san.cnf         # SAN config for certificate
│   │   ├── server.crt      # Self-signed certificate
│   │   └── server.key      # Private key
│   ├── pinger              # Pinger service source code
│   │   ├── Dockerfile      # Dockerfile for Pinger
│   │   ├── config.yaml     # Configuration for Pinger
│   │   └── manifests       # Kubernetes manifests for Pinger
│   ├── ponger              # Ponger service source code
│   │   ├── Dockerfile      # Dockerfile for Ponger
│   │   ├── config.yaml     # Configuration for Ponger
│   │   └── manifests       # Kubernetes manifests for Ponger
```

---

## 🚀 Running the Cluster
We use **k3d** to spin up a local cluster, perfect for small-scale development. Here’s the command to create a cluster with two agent nodes:

```bash
k3d cluster create cluster --k3s-arg '--disable=servicelb@server:0' --k3s-arg '--disable=traefik@server:0' --agents 2
```

---

## 🐳 Dockerizing the Microservices

Both Pinger and Ponger are containerized using Docker. Each service is built using a multi-stage Dockerfile for efficient image creation.

The **Ponger Dockerfile** example briefly:
- **Build stage**: Compiles the Go binary.
- **Final stage**: Uses a lightweight Alpine image to run the app.

To build the image for the **Ponger** service:
```bash
docker build -t ponger:latest app -f app/ponger/Dockerfile
```

---

## 🔐 Securing Communication with TLS
We use **self-signed certificates** (located in the `certs` directory) to secure inter-service communication between Pinger and Ponger. While this works for testing, in production, it's better to integrate a certificate authority or use **Secrets Managers** like AWS Secrets Manager for improved security. 

 - Using `openssl`, I created a self-signed certificate with a proper SAN (Subject Alternative Name) to ensure that the TLS certificate is valid for the service.
   - **Steps**:
     1. Create the SAN configuration:
        ```bash
        cat <<EOF > san.cnf
        [req]
        distinguished_name = req_distinguished_name
        req_extensions = req_ext
        x509_extensions = v3_req
        prompt = no

        [req_distinguished_name]
        C = CA
        ST = Ontario
        L = Toronto
        O = MyOrg
        CN = ponger

        [req_ext]
        subjectAltName = @alt_names

        [v3_req]
        keyUsage = keyEncipherment, dataEncipherment
        extendedKeyUsage = serverAuth
        subjectAltName = @alt_names

        [alt_names]
        DNS.1 = ponger
        EOF
        ```

     2. Generate the certificate:
        ```bash
        openssl req -x509 -nodes -newkey rsa:2048 -keyout tls.key -out tls.crt -days 365 -config san.cnf
        ```
---

## 🛡️ Kubernetes Deployment

The services are deployed using Kubernetes **manifests** located in the `manifests` directory. These define:
- **Deployments**: Create the pods for Pinger and Ponger.
- **Services**: Expose the pods internally within the cluster.
- **Network Policies**: Ensure secure communication between services.

Use `kubectl` commands to manage and observe logs:
```bash
kubectl logs svc/pinger
kubectl logs svc/ponger
```

### Sample Logs:
- **Pinger**:
  ```
  2024/10/14 07:38:49 Sent ping
  2024/10/14 07:38:50 Got pong
  ```
- **Ponger**:
  ```
  2024/10/14 07:38:49 Received GET /ping
  2024/10/14 07:38:50 Received GET /ping
  ```

### Port Forwarding for Local Testing:
To expose the **Ponger** service on port 8080:
```bash
kubectl port-forward svc/ponger --address=0.0.0.0 8080:443
```
You can then access the service using:
```bash
https://<your-ec2-public-ip>:8080/ping
```

---

## 🔐 Best Practices for Security
When pushing your code to GitHub, **never commit sensitive files** like **certificates or private keys**. You can achieve this by adding these to `.gitignore`:

```bash
# Exclude certificates and keys
certs/*.crt
certs/*.key
```

⚠️ **Note**: Failing to manage secrets properly can lead to serious security vulnerabilities.

---

## 🚀 Recommendations for Production
This project demonstrates a simple microservice architecture but to take it to production, additional considerations are essential to ensure scalability, security, and reliability. Here are some key enhancements:

🔐 Use Trusted Certificates: Self-signed certificates are a great start for local or test environments, but in production, it's best to use certificates issued by a trusted Certificate Authority (CA). Automate certificate management with tools like Let’s Encrypt or AWS Certificate Manager to ensure your certificates are always up to date and secure.

🎯 Set up Ingress Controllers with TLS Termination: Ingress controllers like NGINX or Traefik offer better scalability and flexibility by handling TLS termination at the edge. This removes the burden from your services, allowing for easier certificate management and more advanced routing rules.

📊 Implement Monitoring and Logging: Leverage tools like Prometheus for metrics collection and Grafana for visualization to monitor the health and performance of your services. This setup allows you to proactively catch potential issues and respond to security incidents in real time, ensuring production systems stay reliable.

🐳 Leverage Docker Hub or Private Registries: Host your Docker images on Docker Hub or use a private container registry for enhanced security. Private registries help you control access to your images, ensure consistency across deployments, and integrate easily into your CI/CD pipelines.

🎯 Streamline Kubernetes Deployments with Helm Charts: Use Helm Charts to package and manage Kubernetes resources in a more streamlined and scalable way. Helm allows you to handle multi-environment deployments effortlessly and apply rolling updates to minimize downtime during production releases.

🔐 Use a Secrets Manager for Sensitive Data: Store sensitive information such as certificates, API keys, and credentials in Secrets Managers like AWS Secrets Manager or HashiCorp Vault. This keeps sensitive data out of your config files and version control while enabling automatic rotation of secrets for enhanced security.

---

## 📚 Resources
For additional learning:
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [Docker Documentation](https://docs.docker.com/)
- [TLS with OpenSSL](https://www.openssl.org/docs/)

---

Feel free to fork, clone, and experiment with this project! 😊
