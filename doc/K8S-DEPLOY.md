# Deploy Mastro to a K8s Cluster
Mastro is a stateless service that can be easily deployed to a K8s cluster. 

## Precondition
In the examples below, we assume a previously deployed mongo database available on the same namespace or any reachable host at `mongo-mongodb:27017`.

For instance, we used the one using a StatefulSet and deployed as a Helm chart provided by bitnami (see [here](https://bitnami.com/stack/mongodb/helm)).

## Docker Images

The catalogue and the feature store are services that can be easily compiled statically and moved across environments. 
For this reason, we will be using the `pilillo/mastro-catalogue:20210306-static` and `pilillo/mastro-featurestore:20210306-static` in the examples below.
The main difference is the flag `CGO_ENABLED=1` (as by default) set in the dynamically compiled version (the static version has it set to 0). Please have a look at [Dockerfile](../Dockerfile) and [Dockerfile.static](../Dockerfile.static).
On the contrary, the crawler may depend on system libraries (e.g. Kerberos auth libraries) and requires being compiled dynamically. To this end, we will be using the `pilillo/mastro-crawlers:20210304` insted of its `pilillo/mastro-crawlers:20210304-static` counterpart.

## Catalogue

### Config Map

The config for the catalogue can be defined as a K8s config map, as follows:

```
apiVersion: v1
data:
  catalogue-conf.yaml: |
    type: catalogue
    details:
      port: 8085
    backend:
      name: catalogue-mongo
      type: mongo
      settings:
        database: mastro
        collection: mastro-catalogue
        connection-string: "mongodb://mastro:mastro@mongo-mongodb:27017/mastro"
kind: ConfigMap
metadata:
  name: catalogue-conf
```

Mind that in the example above we specified directly the DB user and password (i.e., `mastro:mastro`).
A K8s secret or one injected by an external vault (e.g. hashicorp) can be used for this purpose.

### Deployment

A deployment can be created to spawn multiple replicas for the catalogue.

The configuration is mounted as volume and its path set using the MASTRO_CONFIG variable.

```
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: mastro-catalogue
  name: mastro-catalogue
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mastro-catalogue
  strategy: {}
  template:
    metadata:
      labels:
        app: mastro-catalogue
    spec:
      containers:
      - image: pilillo/mastro-catalogue:20210306-static
        imagePullPolicy: Always
        name: mastro-catalogue
        resources: {}
        ports:
        - containerPort: 8085
          protocol: TCP
        env:
        - name: MASTRO_CONFIG
          value: /conf/catalogue-conf.yaml
        volumeMounts:
        - mountPath: /conf
          name: catalogue-conf-volume
      securityContext: {}
      volumes:
      - name: catalogue-conf-volume
        configMap:
          defaultMode: 420
          name: catalogue-conf
```

### Service

A service is created with:

```
apiVersion: v1
kind: Service
metadata:
  labels:
    app: mastro-catalogue
  name: mastro-catalogue
spec:
  ports:
  - name: rest-8085
    port: 8085
    protocol: TCP
    targetPort: 8085
  selector:
    app: mastro-catalogue
  type: ClusterIP
```

Mind that the service only exposes the catalogue across the namespace.

You will have to create an ingress or a route (respectively on plain K8s and openshift) to make it reachable from the outside world.

## Feature Store

### Config Map

```
apiVersion: v1
data:
  fs-conf.yaml: |
    type: featurestore
    details:
      port: 8085
    backend:
      name: fs-mongo
      type: mongo
      settings:
        database: mastro
        collection: mastro-featurestore
        connection-string: "mongodb://mastro:mastro@mongo-mongodb:27017/mastro"
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: fs-conf
```

### Deployment

```
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: mastro-featurestore
  name: mastro-featurestore
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mastro-featurestore
  strategy: {}
  template:
    metadata:
      labels:
        app: mastro-featurestore
    spec:
      containers:
      - image: pilillo/mastro-featurestore:20210306-static
        imagePullPolicy: Always
        name: mastro-featurestore
        resources: {}
        ports:
        - containerPort: 8085
          protocol: TCP
        env:
        - name: MASTRO_CONFIG
          value: /conf/fs-conf.yaml
        volumeMounts:
        - mountPath: /conf
          name: fs-conf-volume
      securityContext: {}
      volumes:
      - name: fs-conf-volume
        configMap:
          defaultMode: 420
          name: fs-conf
```

### Service

```
apiVersion: v1
kind: Service
metadata:
  labels:
    app: mastro-featurestore
  name: mastro-featurestore
spec:
  ports:
  - name: rest-8085
    port: 8085
    protocol: TCP
    targetPort: 8085
  selector:
    app: mastro-featurestore
  type: ClusterIP
```

## Crawler

For the following example we use the image `pilillo/mastro-crawlers:20210304` including a dynamically built binary of the crawlers.

The crawling agent can be easily debugged locally by overwriting the default entrypoint:

```
docker run --entrypoint "/bin/sh" -it pilillo/mastro-crawlers:20210304
```

### ConfigMap

A config for the crawler can be mounted as config map, for instance:

```
apiVersion: v1
data:
  crawler-conf.yaml: |
    type: crawler
    backend:
      name: impala-enterprise-datalake
      type: impala
      crawler:
        root: ""
        schedule-period: sunday
        schedule-value: 1
        start-now: true
        catalogue-endpoint: "http://mastro-catalogue:8085/assets/"
      settings:
        host: "impala.domain.com"
        port: "21000"
        use-kerberos: true
kind: ConfigMap
metadata:
  name: crawler-conf
```

This sets the agent to run every sunday, as well as right now after its Pod is created.

### Kerberos Authentication

For the example Impala crawler, we need to both spawn a mastro-crawler and a Kerberos authentication process.
To this end, we use an init container, doing a `kinit` on behalf the user and renewing the ticket cache upon expiration.
We previously documented this process in [this blog post](http://p111110.blogspot.com/2021/03/kerberos-auth-on-k8sopenshift-using.html).
Specifically, we rely on another Github project, named [Geronzio](https://github.com/pilillo/geronzio) to automatically build a Docker container including (see Dockerfile [here](https://github.com/pilillo/geronzio/blob/main/Dockerfile)) krb5 and kstart, respectively the Kerberos client libraries and k5start client.
For Kerberos configuration we need: i) a krb5.conf file, and ii) a keytab or password to authenticate.

#### krb5.conf

A `krb5.conf` file defines the REALM and location of the KDC. 
See [here](https://web.mit.edu/kerberos/krb5-1.12/doc/admin/conf_files/krb5_conf.html) for a full documentation.

```
apiVersion: v1
data:
  krb5.conf: |
    [logging]
    ...

    [libdefaults]
    ...
    
    [realms]
    ...
    
    [appdefaults]
    ..
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: krb5-conf
```

#### User Keytab file

The user keytab can be mounted as a secret, directly on K8s, or mounted from an external Vault.
For instance:

```
apiVersion: v1
data:
  user.keytab: blablablablablablablablablabla
kind: Secret
metadata:
  name: user-keytab
```


### Deployments, Jobs and CronJobs

Depending on the crawled source, a crawler may be scheduled to run once or periodically. 
There are 3 possibilities to deploy a crawler on K8s: using **i) a deployment, ii) a Job or a iii) Cron Job.**
When using a deployment the `github.com/go-co-op/gocron` library is used to schedule the agent runs. 
The deployment implies a Pod being created to run either once or periodically.
On K8s, however, the Job and CronJob resources can be used for a similar purpose, respectively for one-time and periodical jobs.

#### Deployment

To deploy the crawler with an auth sidecar container, the following steps are taken:
- add as a sidecar a container having an available kerberos client;
- mount the krb5.conf map on both the application container and the sidecar (as a read only volume); as you may try and see, it is not a good idea to mount something at /etc as kubernetes normally injects host and dns info at this location and may result in the Pod being rejected by the admission controller or an error. A similar behavior may occurr at /tmp. So just change the default paths with something creative. Mind that we can use KRB5_CONFIG and KRB5CCNAME to respectively overwrite the default location of the krb5.conf and cache files. Specifically, the cache file can be set to be written to an ephimeral volume used as communication means between main and sidecar container.
- mount the keytab secret on the sidecar container, i.e., as a read only volume at the /keytabs location.

```
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: mastro-impala-crawler
  name: mastro-impala-crawler
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mastro-impala-crawler
  strategy: {}
  template:
    metadata:
      labels:
        app: mastro-impala-crawler
    spec:
      containers:
      # sidecar container
      - image: pilillo/geronzio:20210305
        imagePullPolicy: IfNotPresent
        name: geronzio
        env:
        - name: KRB5_CONFIG
          value: /etc-krb5/krb5.conf
        - name: KRB5CCNAME
          value: /tmp-krb5/krb5cc
        - name: KRBUSER
          value: SMARTUSER01
        - name: REALM
          value: DOMAIN.COM
        command: ["kinit", "-kt", "/keytabs/user.keytab", "$(KRBUSER)@$(REALM)"]
        restartPolicy: OnFailure
        lifecycle:
          type: Sidecar
        volumeMounts:
        - mountPath: /keytabs
          name: keytab-volume
          readOnly: true
        - mountPath: /etc-krb5
          name: krb5-conf-volume
          readOnly: true
        - mountPath: /tmp-krb5
          name: shared-cache
      # actual crawler
      - image: pilillo/mastro-crawlers:20210305
        imagePullPolicy: Always
        #IfNotPresent
        name: mastro-crawler
        resources: {}
        env:
        - name: KRB5_CONFIG
          value: /etc-krb5/krb5.conf
        - name: KRB5CCNAME
          value: /tmp-krb5/krb5cc
        - name: MASTRO_CONFIG
          value: /conf/crawler-conf.yaml
        volumeMounts:
        - mountPath: /conf
          name: crawler-conf-volume
        - mountPath: /etc-krb5
          name: krb5-conf-volume
          readOnly: true
        - mountPath: /tmp-krb5
          name: shared-cache
      securityContext: {}
      volumes:
      - name: crawler-conf-volume
        configMap:
          defaultMode: 420
          name: crawler-conf
      - name: krb5-conf-volume
        configMap:
          defaultMode: 420
          name: krb5-conf
      - name: shared-cache
        emptyDir: {}
      - name: keytab-volume
        secret:
          secretName: user-keytab
```

#### Job

A K8s Job is a batch process meant to run once. For instance:

```
apiVersion: batch/v1
kind: Job
metadata:
  creationTimestamp: null
  name: mastro-impala-crawler
spec:
  template:
    metadata:
      creationTimestamp: null
    spec:
      containers:
      - image: pilillo/mastro-crawlers:20210305
        name: mastro-crawler
        resources: {}
      restartPolicy: Never
status: {}
```
This is only a Job example. Please refer to the full description provided in the Deployment case for the complete Impala deployment.

#### CronJob

A CronJob can be created with the following syntax.

```
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  creationTimestamp: null
  name: mastro-impala-crawler
spec:
  jobTemplate:
    metadata:
      creationTimestamp: null
      name: mastro-impala-crawler
    spec:
      template:
        metadata:
          creationTimestamp: null
        spec:
          containers:
          - image: pilillo/mastro-crawlers:20210305
            name: mastro-crawler
            resources: {}
          restartPolicy: OnFailure
  schedule: 0 0 * * 0
status: {}
```

This is only a CronJob example. Please refer to the full description provided in the Deployment case for the complete Impala deployment.

In this section we describe the introduced schedule format, i.e., the format string defining the schedule interval of cron jobs.

Specifically, this consists of the following 5 fields:

| Field | Description      | Values                              |
|-------|------------------|-------------------------------------|
| 1     | Minute           | 0 to 59, or *                       |
| 2     | Hour             | 0 to 23, or *                       |
| 3     | Day of the Month | 1 to 31, or *                       |
| 4     | Month            | 1 to 12, or *                       |
| 5     | Day of the Week  | 0 to 7, with (0 == 7, sunday), or * |

Mind that the string must contain entries for each field, or an asterisk (i.e., `*` otherwise).

For instance:
- `0 0 * * 0` schedules the job every sunday at midnight (00:00)


Also, instead of specifying a specific time, a slash and a period size can be defined to periodically schedule at the field granularity:
- `*/5 * * * *` schedules every 5 minutes
- `0 */2 * * *` schedule every second hour, at o'clock
