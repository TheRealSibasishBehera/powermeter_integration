<h1>PowerMeter Monitoring</h1>
<h3>Using Prometheus and Graphana</h3>

Using OS- tools can more or less give accurate information about the energy consumption by the operating system and the CPU. In any general device, this could be sufficient as in such devices as the most energy supply Is used by CPU computation.

But in cases of SBCs like Raspberry Pi, Nvidia Jetson Boards this is not sufficient as there are other factors as well which consume energy comparable to that of CPU computation. Eg.Energy consumption by networking components like Ethernet, GPIO, UART, etc.

To overcome this problem and to get the complete power consumption, an external power meter needs to be integrated into the SBC

![f83f1a_e2d765a97e914125b200aeefdd810740_mv2](https://user-images.githubusercontent.com/95071627/188454874-a0691223-b8ec-4ee5-b1ac-c28d69ba6306.jpg)

For this we  would be using Tasmota EU plug V2 by Athom . This is based on tasmota-HLW8032 , providing control using MQTT,Web UI , HTTP.

For installation of flotta-operator and flotta-edge device  
follow the flotta guide

[kind installation](https://project-flotta.io/documentation/v0_2_0/gsg/kind.html)

[flotta-dev-cli](https://project-flotta.io/flotta/2022/07/20/developer-cli.html)

<h3>Deploying PoweMeter Workload</h3>

The powertop monitoring application would be deployed as workloads.
Details on how to deploying workloads are in
[flotta workloads deployment](https://project-flotta.io/documentation/v0_2_0/gsg/running_workloads.html)

The yaml for the workload  :-

```yaml
apiVersion: management.project-flotta.io/v1alpha1
kind: EdgeWorkload
metadata:
   name: powermeter
spec:
   metrics:
      interval: 5
      path: "/metrics"
      port: 8881
   deviceSelector:
      matchLabels:
         app: foo
   type: pod
   pod:
      spec:
         hostNetwork: true
         containers:
            - name: powermeter
              image: docker.io/sibseh/powermeter:v4    
```

<h3>Monitoring Using Thanos</h3>

A thanos and graphana set up can be used for monitoring visually

More details can be found in :-

[flotta observability](https://project-flotta.io/documentation/latest/operations/observability.html)

[writing-metrics-to-control-plane](https://project-flotta.io/flotta/2022/04/11/writing-metrics-to-control-plane.html
)


For Thanos receiver

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: thanos-receiver
  labels:
    app: thanos-receiver
spec:
  replicas: 1
  selector:
    matchLabels:
      app: thanos-receiver
  template:
    metadata:
      labels:
        app: thanos-receiver
    spec:
      containers:
        - name: receive
          image: quay.io/thanos/thanos:v0.24.0
          command:
            - /bin/thanos
            - receive
            - --log.level
            - debug
            - --label
            - "receiver=\"0\""
            - --remote-write.address
            - 0.0.0.0:10908
        - name: query
          image: quay.io/thanos/thanos:v0.24.0
          command:
            - /bin/thanos
            - query
            - --log.level
            - debug
            - --http-address
            - 0.0.0.0:9090
            - --grpc-address
            - 0.0.0.0:11901
            - --endpoint
            - 127.0.0.1:10901
---
apiVersion: v1
kind: Service
metadata:
  name: thanos-receiver
spec:
  type: NodePort
  selector:
    app: thanos-receiver
  ports:
    - port: 80
      targetPort: 10908
      nodePort: 30030
      name: endpoint
    - port: 9090
      targetPort: 9090
      name: admin
      
---
apiVersion: v1
kind: Service
metadata:
  name: thanos-receiver
spec:
  type: NodePort
  selector:
    app: thanos-receiver
  ports:
    - port: 80
      targetPort: 10908
      nodePort: 30030
      name: endpoint
    - port: 9090
      targetPort: 9090
      name: admin

```


For Graphana
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      name: grafana
      labels:
        app: grafana
    spec:
      containers:
        - name: grafana
          image: grafana/grafana:latest
          ports:
            - name: grafana
              containerPort: 3000
          resources:
            limits:
              memory: "1Gi"
              cpu: "1000m"
            requests:
              memory: 500M
              cpu: "500m"
          volumeMounts:
            - mountPath: /var/lib/grafana
              name: grafana-storage
      volumes:
        - name: grafana-storage
```

